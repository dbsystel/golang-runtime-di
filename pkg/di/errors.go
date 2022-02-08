package di

import (
	"reflect"

	"github.com/pkg/errors"
)

func errNoStructPtr(tpe reflect.Type) error {
	return errors.Errorf("expected a struct pointer, but got: %v", tpe)
}

func errNotCoercible(tgt reflect.Type, src reflect.Type) error {
	return errors.Errorf("cannot coerce '%v' to: %v", src, tgt)
}

func errFieldNotExported(tpe reflect.Type, fld reflect.StructField) error {
	return errors.Errorf("field not exported in type '%v': %v", tpe, fld.Name)
}
