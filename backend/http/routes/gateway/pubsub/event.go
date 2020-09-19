package pubsub

import (
	"encoding/json"
	"reflect"

	"github.com/pkg/errors"
)

type Event struct {
	Event string          `json:"e"`
	Data  json.RawMessage `json:"d"`
}

var zeroEvent = Event{}

func NewEvent(v interface{}) (Event, error) {
	j, err := json.Marshal(v)
	if err != nil {
		return zeroEvent, errors.Wrap(err, "failed to encode event data")
	}

	return Event{
		Event: typeName(v),
		Data:  j,
	}, nil
}

func typeName(v interface{}) string {
	return reflect.TypeOf(v).String()
}
