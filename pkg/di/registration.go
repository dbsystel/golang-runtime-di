package di

import (
	"fmt"
	"reflect"
	"runtime"

	"github.com/pkg/errors"
)

// Registration is a registration for a single component, created by a single instance
type Registration struct {
	// FactoryFn is the factory function to be used to create a new instance
	FactoryFn func() (interface{}, error)
	// Type is the target type of the registration
	Type reflect.Type
	// Qualifier is an optional qualifier for the component
	Qualifier string
	// Priority denotes the resolution priority (lower = higher)
	Priority int
	// Source is the string which provides the information of the component's origin (file:line)
	Source string
	// instance is being used to cache the wired instance, once the component is created
	instance interface{}
}

func NewRegistration(val interface{}, skipCaller int) (*Registration, error) {
	if val != nil {
		tpe := reflect.TypeOf(val)
		switch tpe.Kind() {
		case reflect.Invalid, reflect.Uintptr, reflect.UnsafePointer:
			break
		case reflect.Func:
			return newFactoryRegistration(reflect.ValueOf(val), skipCaller+1)
		default:
			return newRegistration(func() (interface{}, error) { return val, nil }, tpe, skipCaller+1), nil
		}
	}
	return nil, errors.Errorf("invalid component type: %v", reflect.TypeOf(val))
}

func newRegistration(fn func() (interface{}, error), tpe reflect.Type, skipCaller int) *Registration {
	_, file, line, _ := runtime.Caller(skipCaller + 1)
	return &Registration{
		FactoryFn: fn,
		Type:      tpe,
		Source:    fmt.Sprintf("%v:%v", file, line),
	}
}

func newFactoryRegistration(val reflect.Value, skipCaller int) (*Registration, error) {
	tpe := val.Type()
	if paramCount := tpe.NumIn(); paramCount > 0 {
		return nil, errors.Errorf("function should not have parameters, but got: %v", paramCount)
	}
	returnCount := tpe.NumOut()
	if returnCount > 0 && returnCount < 3 {
		resultTpe := tpe.Out(0)
		switch resultTpe.Kind() {
		case reflect.Invalid, reflect.Uintptr, reflect.UnsafePointer, reflect.Func:
			return nil, errors.Errorf("invalid factory result type: %v", resultTpe)
		}
		switch returnCount {
		case 1:
			return newRegistration(func() (interface{}, error) {
				return val.Call(nil)[0].Interface(), nil
			}, resultTpe, skipCaller+1), nil
		case 2: // nolint:gomnd
			if errParam := tpe.Out(1); !errParam.Implements(reflect.TypeOf((*error)(nil)).Elem()) {
				return nil, errors.Errorf("second function return value should be error, but is: %v", errParam)
			}
			return newRegistration(func() (interface{}, error) {
				results := val.Call(nil)
				result := results[0].Interface()
				if err := results[1].Interface(); err != nil {
					return result, err.(error)
				}
				return result, nil
			}, resultTpe, skipCaller+1), nil
		}
	}
	return nil, errors.Errorf("function should provide 1 or 2 return values, but has: %v", returnCount)
}

// GetInstance returns the instance of the registration
func (r *Registration) GetInstance() (result interface{}, first bool, err error) {
	// If we have not created an instance for this registration, create it and wire it
	if r.instance == nil {
		first = true
		r.instance, err = r.FactoryFn()
		if err != nil {
			return nil, first, errors.Wrapf(err, "could not create instance: %v", r)
		}
	}
	return r.instance, first, nil
}

// WithQualifier sets the Registration#Qualifier for the registered component returning the same ptr as in the receiver
func (r *Registration) WithQualifier(qualifier string) *Registration {
	r.Qualifier = qualifier
	return r
}

// WithPriority sets the Registration#Priority for the registered component returning the same ptr as in the receiver
func (r *Registration) WithPriority(priority int) *Registration {
	r.Priority = priority
	return r
}

// String returns a descriptor for the Registration
func (r *Registration) String() string {
	name := r.Type.String()
	if len(r.Qualifier) > 0 {
		name = fmt.Sprintf("%v(%v)", name, r.Qualifier)
	}
	return fmt.Sprintf("component %v with priority %v registered at: %v", name, r.Priority, r.Source)
}
