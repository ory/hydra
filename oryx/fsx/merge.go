// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package fsx

import (
	"io"
	"io/fs"
	"sort"
	"time"

	"github.com/pkg/errors"
)

type (
	mergedFS   []fs.FS
	mergedFile struct {
		files                 []fs.File
		unprocessedDirEntries dirEntries
	}
	mergedFileInfo []fs.FileInfo
	dirEntries     []fs.DirEntry
)

var (
	_ fs.StatFS      = (mergedFS)(nil)
	_ fs.ReadDirFS   = (mergedFS)(nil)
	_ fs.ReadDirFile = (*mergedFile)(nil)
	_ fs.FileInfo    = (mergedFileInfo)(nil)
	_ sort.Interface = (dirEntries)(nil)
)

// Merge multiple filesystems. Later file systems are shadowed by previous ones.
func Merge(fss ...fs.FS) fs.FS {
	return mergedFS(fss)
}

func (m mergedFS) Open(name string) (fs.File, error) {
	var file mergedFile
	for _, fsys := range m {
		f, err := fsys.Open(name)
		if errors.Is(err, fs.ErrNotExist) {
			continue
		}
		if err != nil {
			return nil, errors.WithStack(err)
		}

		file.files = append(file.files, f)
	}
	if len(file.files) == 0 {
		return nil, errors.WithStack(fs.ErrNotExist)
	}

	return &file, nil
}

func (m mergedFS) Stat(name string) (fs.FileInfo, error) {
	for i, fsys := range m {
		info, err := fs.Stat(fsys, name)
		if errors.Is(err, fs.ErrNotExist) {
			continue
		}

		switch {
		case err != nil:
			return nil, errors.WithStack(err)
		case info.IsDir():
			dirs := mergedFileInfo{info}
			for j := i + 1; j < len(m); j++ {
				info, err := fs.Stat(m[j], name)
				if errors.Is(err, fs.ErrNotExist) {
					continue
				}
				if err != nil {
					return nil, err
				}
				dirs = append(dirs, info)
			}
			return dirs, nil
		default:
			return info, nil
		}
	}
	return nil, errors.WithStack(fs.ErrNotExist)
}

func (m mergedFS) ReadDir(name string) ([]fs.DirEntry, error) {
	var entries dirEntries

	for _, fsys := range m {
		e, err := fs.ReadDir(fsys, name)
		if errors.Is(err, fs.ErrNotExist) {
			continue
		}
		if err != nil {
			return nil, err
		}
		entries = append(entries, e...)
	}
	if len(entries) == 0 {
		return nil, errors.WithStack(fs.ErrNotExist)
	}

	entries.clean()
	return entries, nil
}

func (m mergedFileInfo) Name() string {
	return m[0].Name()
}

func (m mergedFileInfo) Size() int64 {
	return m[0].Size()
}

func (m mergedFileInfo) Mode() fs.FileMode {
	return m[0].Mode()
}

func (m mergedFileInfo) ModTime() time.Time {
	return m[0].ModTime()
}

func (m mergedFileInfo) IsDir() bool {
	return m[0].IsDir()
}

func (m mergedFileInfo) Sys() interface{} {
	return m
}

func (d dirEntries) Len() int {
	return len(d)
}

func (d dirEntries) Less(i, j int) bool {
	return d[i].Name() < d[j].Name()
}

func (d dirEntries) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

func (d *dirEntries) clean() {
	sort.Sort(d)

	for i := 1; i < len(*d); i++ {
		if (*d)[i-1].Name() == (*d)[i].Name() {
			if len(*d)-i == 1 {
				// remove the last entry; we're done
				*d = (*d)[:i]
				return
			}
			// remove the duplicate entry at index i
			*d = append((*d)[:i], (*d)[i+1:]...)

			// need to check the same index again
			i--
		}
	}
}

func (m *mergedFile) Stat() (fs.FileInfo, error) {
	return m.files[0].Stat()
}

func (m *mergedFile) Read(bytes []byte) (int, error) {
	return m.files[0].Read(bytes)
}

func (m *mergedFile) Close() error {
	var firstErr error
	for _, f := range m.files {
		if err := f.Close(); err != nil {
			if firstErr == nil {
				firstErr = errors.WithStack(err)
			}
		}
	}
	return firstErr
}

func (m *mergedFile) ReadDir(n int) ([]fs.DirEntry, error) {
	if m.unprocessedDirEntries != nil {
		if n <= 0 {
			entries := m.unprocessedDirEntries
			m.unprocessedDirEntries = nil
			return entries, nil
		}
		if n >= len(m.unprocessedDirEntries) {
			entries := m.unprocessedDirEntries
			m.unprocessedDirEntries = nil
			return entries, io.EOF
		}

		var entries dirEntries
		entries, m.unprocessedDirEntries = m.unprocessedDirEntries[:n], m.unprocessedDirEntries[n:]
		return entries, nil
	}

	var entries dirEntries
	for _, f := range m.files {
		if f, ok := f.(fs.ReadDirFile); ok {
			e, err := f.ReadDir(-1)
			if err != nil && !errors.Is(err, fs.ErrNotExist) {
				return nil, err
			}
			entries = append(entries, e...)
		}
	}
	if entries == nil {
		if n > 0 {
			return nil, io.EOF
		}
		return nil, nil
	}

	entries.clean()
	if n <= 0 {
		return entries, nil
	}
	if n >= len(entries) {
		return entries, io.EOF
	}

	entries, m.unprocessedDirEntries = entries[:n], entries[n:]
	return entries, nil
}
