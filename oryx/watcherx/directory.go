// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package watcherx

import (
	"context"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
)

func WatchDirectory(ctx context.Context, dir string, c EventChannel) (Watcher, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	subDirs := make(map[string]struct{})
	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.WithStack(err)
		}
		if info.IsDir() {
			if err := w.Add(path); err != nil {
				return errors.WithStack(err)
			}
			subDirs[path] = struct{}{}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	dw := &directoryWatcher{
		dispatcher: newDispatcher(),
		c:          c,
		dir:        dir,
		subDirs:    subDirs,
		w:          w,
	}
	go dw.streamDirectoryEvents(ctx)
	return dw, nil
}

type directoryWatcher struct {
	*dispatcher
	c       EventChannel
	dir     string
	subDirs map[string]struct{}
	w       *fsnotify.Watcher
}

func (w *directoryWatcher) handleEvent(ctx context.Context, e fsnotify.Event) {
	if e.Has(fsnotify.Remove) {
		if _, ok := w.subDirs[e.Name]; ok {
			// we do not want any event on deletion of a directory
			delete(w.subDirs, e.Name)
			return
		}
		w.maybeSend(ctx, &RemoveEvent{
			source: source(e.Name),
		})
		return
	} else if e.Has(fsnotify.Write | fsnotify.Create) {
		if stats, err := os.Stat(e.Name); err != nil {
			w.maybeSend(ctx, &ErrorEvent{
				error:  errors.WithStack(err),
				source: source(e.Name),
			})
			return
		} else if stats.IsDir() {
			if err := w.w.Add(e.Name); err != nil {
				w.maybeSend(ctx, &ErrorEvent{
					error:  errors.WithStack(err),
					source: source(e.Name),
				})
			}
			w.subDirs[e.Name] = struct{}{}
			return
		}

		//#nosec G304 -- false positive
		data, err := os.ReadFile(e.Name)
		if err != nil {
			w.maybeSend(ctx, &ErrorEvent{
				error:  err,
				source: source(e.Name),
			})
		} else {
			w.maybeSend(ctx, &ChangeEvent{
				data:   data,
				source: source(e.Name),
			})
		}
	}
}

func (w *directoryWatcher) maybeSend(ctx context.Context, e Event) bool {
	select {
	case <-ctx.Done():
		return false
	case w.c <- e:
		return true
	}
}

func (w *directoryWatcher) streamDirectoryEvents(ctx context.Context) {
	defer func() {
		close(w.done)
		close(w.c)
		_ = w.w.Close()
	}()
	for {
		select {
		case <-ctx.Done():
			return
		case e := <-w.w.Events:
			w.handleEvent(ctx, e)
		case <-w.trigger:
			var eventsSent int

			if err := filepath.Walk(w.dir, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.IsDir() {
					//#nosec G304 -- false positive
					data, err := os.ReadFile(path)
					if err != nil {
						if !w.maybeSend(ctx, &ErrorEvent{
							error:  err,
							source: source(path),
						}) {
							return errors.WithStack(context.Canceled)
						}
					} else {
						if !w.maybeSend(ctx, &ChangeEvent{
							data:   data,
							source: source(path),
						}) {
							return errors.WithStack(context.Canceled)
						}
					}
					eventsSent++
				}
				return nil
			}); err != nil {
				if !w.maybeSend(ctx, &ErrorEvent{
					error:  err,
					source: source(w.dir),
				}) {
					return
				}
				eventsSent++
			}

			w.done <- eventsSent
		}
	}
}
