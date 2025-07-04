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

	d := newDispatcher()
	go streamDirectoryEvents(ctx, w, c, d.trigger, d.done, dir, subDirs)
	return d, nil
}

func handleEvent(e fsnotify.Event, w *fsnotify.Watcher, c EventChannel, subDirs map[string]struct{}) {
	if e.Has(fsnotify.Remove) {
		if _, ok := subDirs[e.Name]; ok {
			// we do not want any event on deletion of a directory
			delete(subDirs, e.Name)
			return
		}
		c <- &RemoveEvent{
			source: source(e.Name),
		}
		return
	} else if e.Has(fsnotify.Write | fsnotify.Create) {
		if stats, err := os.Stat(e.Name); err != nil {
			c <- &ErrorEvent{
				error:  errors.WithStack(err),
				source: source(e.Name),
			}
			return
		} else if stats.IsDir() {
			if err := w.Add(e.Name); err != nil {
				c <- &ErrorEvent{
					error:  errors.WithStack(err),
					source: source(e.Name),
				}
			}
			subDirs[e.Name] = struct{}{}
			return
		}

		//#nosec G304 -- false positive
		data, err := os.ReadFile(e.Name)
		if err != nil {
			c <- &ErrorEvent{
				error:  err,
				source: source(e.Name),
			}
		} else {
			c <- &ChangeEvent{
				data:   data,
				source: source(e.Name),
			}
		}
	}
}

func streamDirectoryEvents(ctx context.Context, w *fsnotify.Watcher, c EventChannel, sendNow <-chan struct{}, sendNowDone chan<- int, dir string, subDirs map[string]struct{}) {
	for {
		select {
		case <-ctx.Done():
			_ = w.Close()
			return
		case e := <-w.Events:
			handleEvent(e, w, c, subDirs)
		case <-sendNow:
			var eventsSent int

			if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.IsDir() {
					//#nosec G304 -- false positive
					data, err := os.ReadFile(path)
					if err != nil {
						c <- &ErrorEvent{
							error:  err,
							source: source(path),
						}
					} else {
						c <- &ChangeEvent{
							data:   data,
							source: source(path),
						}
					}
					eventsSent++
				}
				return nil
			}); err != nil {
				c <- &ErrorEvent{
					error:  err,
					source: source(dir),
				}
				eventsSent++
			}

			sendNowDone <- eventsSent
		}
	}
}
