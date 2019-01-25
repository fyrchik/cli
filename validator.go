package cli

import (
	"reflect"

	"github.com/pkg/errors"
)

var (
	_ Validator = OneOf()
)

// OneOf returns Validator which requires value to be in provided set.
func OneOf(values ...interface{}) Validator {
	return func(val interface{}) error {
		for _, v := range values {
			if reflect.DeepEqual(val, v) {
				return nil
			}
		}
		return errors.Errorf("must be one of: %v", values)
	}
}
