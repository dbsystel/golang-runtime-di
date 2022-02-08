package di

import (
	"reflect"

	"github.com/pkg/errors"
)

// Injection represents an injection to a reflect.StructField
type Injection reflect.StructField

// Apply injects the field to the provided reflect.Value
func (i Injection) Apply(target reflect.Value, val reflect.Value) error {
	// Inject to struct ptrs only
	if (target.Kind() != reflect.Ptr && target.Kind() != reflect.Interface) || target.Elem().Kind() != reflect.Struct {
		return errNoStructPtr(target.Type())
	}
	fld := target.Elem().FieldByName(i.Name)
	if !fld.IsValid() {
		return errors.Errorf("field '%v' is not valid in target: %v", i.Name, target.Type())
	}
	if !isCoercible(fld.Type(), val.Type()) {
		return errNotCoercible(fld.Type(), val.Type())
	}
	fld.Set(val)
	return nil
}

func isCoercible(tgt reflect.Type, src reflect.Type) bool {
	if tgt.Kind() == reflect.Interface {
		return src.Implements(tgt)
	}
	return src.AssignableTo(tgt)
}
