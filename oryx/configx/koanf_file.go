// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package configx

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/v2"

	"github.com/pkg/errors"

	"github.com/ory/x/watcherx"
)

// KoanfFile implements a KoanfFile provider.
type KoanfFile struct {
	subKey string
	path   string
	parser koanf.Parser
}

// NewKoanfFile returns a file provider.
func NewKoanfFile(path string) (*KoanfFile, error) {
	return NewKoanfFileSubKey(path, "")
}

func NewKoanfFileSubKey(path, subKey string) (*KoanfFile, error) {
	kf := &KoanfFile{
		path:   filepath.Clean(path),
		subKey: subKey,
	}

	switch e := filepath.Ext(path); e {
	case ".toml":
		kf.parser = toml.Parser()
	case ".json":
		kf.parser = json.Parser()
	case ".yaml", ".yml":
		kf.parser = yaml.Parser()
	default:
		return nil, errors.Errorf("unknown config file extension: %s", e)
	}

	return kf, nil
}

// ReadBytes is not supported by KoanfFile.
func (f *KoanfFile) ReadBytes() ([]byte, error) {
	return nil, errors.New("file provider does not support this method")
}

// Read reads the file and returns the parsed configuration.
func (f *KoanfFile) Read() (map[string]interface{}, error) {
	//#nosec G304 -- false positive
	fc, err := os.ReadFile(f.path)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	v, err := f.parser.Unmarshal(fc)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if f.subKey == "" {
		return v, nil
	}

	path := strings.Split(f.subKey, Delimiter)
	for i := range path {
		v = map[string]interface{}{
			path[len(path)-1-i]: v,
		}
	}

	return v, nil
}

// WatchChannel watches the file and triggers a callback when it changes. It is a
// blocking function that internally spawns a goroutine to watch for changes.
func (f *KoanfFile) WatchChannel(ctx context.Context, c watcherx.EventChannel) (watcherx.Watcher, error) {
	return watcherx.WatchFile(ctx, f.path, c)
}
