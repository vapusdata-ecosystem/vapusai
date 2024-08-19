package datasvc

import (
	"reflect"
)

func validateType(t reflect.Type, expected reflect.Kind) (reflect.Type, error) {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != expected {
		return nil, ErrInvalidDestinationType
	}
	return t, nil
}
