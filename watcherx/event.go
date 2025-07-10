package watcherx

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
)

type (
	Event interface {
		Source() string
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

const (
	serialTypeChange serialEventType = "change"
	serialTypeRemove serialEventType = "remove"
	serialTypeError  serialEventType = "error"
)

var errUnknownEvent = errors.New("unknown event type")

func (e *ErrorEvent) String() string {
	return fmt.Sprintf("error: %+v; source: %s", e.error, e.source)
}

func (e source) Source() string {
	return string(e)
}

func (e *source) setSource(nsrc string) {
	*e = source(nsrc)
}

func unmarshalEvent(data []byte) (Event, error) {
	var serialEvent serialEvent
	if err := json.Unmarshal(data, &serialEvent); err != nil {
		return nil, errors.WithStack(err)
	}
	switch serialEvent.Type {
	case serialTypeRemove:
		return &RemoveEvent{
			source: serialEvent.Source,
		}, nil
	case serialTypeChange:
		return &ChangeEvent{
			data:   serialEvent.Data,
			source: serialEvent.Source,
		}, nil
	case serialTypeError:
		return &ErrorEvent{
			error:  errors.New(string(serialEvent.Data)),
			source: serialEvent.Source,
		}, nil
	}
	return nil, errUnknownEvent
}
