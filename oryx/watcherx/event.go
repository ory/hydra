// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package watcherx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

type (
	Event interface {
		// MarshalJSON is required to work multiple times
		json.Marshaler

		Reader() io.Reader
		Source() string
		String() string
		setSource(string)
	}
	source     string
	ErrorEvent struct {
		error
		source
	}
	ChangeEvent struct {
		data []byte
		source
	}
	RemoveEvent struct {
		source
	}
	serialEventType string
	serialEvent     struct {
		Type   serialEventType `json:"type"`
		Data   []byte          `json:"data"`
		Source source          `json:"source"`
	}
)

func NewErrorEvent(err error, source_ string) *ErrorEvent {
	return &ErrorEvent{
		error:  err,
		source: source(source_),
	}
}

const (
	serialTypeChange serialEventType = "change"
	serialTypeRemove serialEventType = "remove"
	serialTypeError  serialEventType = "error"
)

func (e *ErrorEvent) Reader() io.Reader {
	return bytes.NewBufferString(e.Error())
}

func (e *ErrorEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(serialEvent{
		Type:   serialTypeError,
		Data:   []byte(e.Error()),
		Source: e.source,
	})
}

func (e *ErrorEvent) String() string {
	return fmt.Sprintf("error: %+v; source: %s", e.error, e.source)
}

func (e source) Source() string {
	return string(e)
}

func (e *source) setSource(nsrc string) {
	*e = source(nsrc)
}

func (e *ChangeEvent) Reader() io.Reader {
	return bytes.NewBuffer(e.data)
}

func (e *ChangeEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(serialEvent{
		Type:   serialTypeChange,
		Data:   e.data,
		Source: e.source,
	})
}

func (e *ChangeEvent) String() string {
	return fmt.Sprintf("data: %s; source: %s", e.data, e.source)
}

func (e *RemoveEvent) Reader() io.Reader {
	return nil
}

func (e *RemoveEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(serialEvent{
		Type:   serialTypeRemove,
		Source: e.source,
	})
}

func (e *RemoveEvent) String() string {
	return fmt.Sprintf("removed source: %s", e.source)
}
