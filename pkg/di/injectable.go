package di

import (
	"reflect"

	"github.com/pkg/errors"
)

const (
	TagKey = "inject"
)

type Injectable struct {
	// Type is the target struct type to be injected to
	Type reflect.Type
	// Injections denote the fields to be injected to
	Injections []Injection
}

func (i Injectable) Apply(target reflect.Value, resolver InstanceResolver) error {
	if (target.Kind() != reflect.Ptr && target.Kind() != reflect.Interface) || target.Elem().Kind() != reflect.Struct {
		return errNoStructPtr(target.Type())
	}
	// Make sure we end up with the correct type to be injected to
	if !target.Type().Elem().AssignableTo(i.Type) {
		return errNotCoercible(i.Type, target.Type())
	}
	// Apply all injections
	for _, injection := range i.Injections {
		resolved, err := resolver.ResolveInstance(injection.Type, TagValueFrom(injection.Tag.Get(TagKey)))
		if err != nil {
			return errors.Wrapf(err, "could not resolve component for field: %v", injection.Name)
		}
		if resolved == reflect.ValueOf(nil) {
			continue
		}
		if err = injection.Apply(target, resolved); err != nil {
			return errors.Wrapf(err, "could not inject field: %v", injection.Name)
		}
	}
	return nil
}

// InjectableFrom creates an Injectable from a reflect.Type
func InjectableFrom(tpe reflect.Type) (*Injectable, error) {
	// Unwrap pointers and interfaces
	for tpe != nil && (tpe.Kind() == reflect.Ptr || tpe.Kind() == reflect.Interface) {
		tpe = tpe.Elem()
	}
	// Make sure we end up with a struct
	if tpe == nil || tpe.Kind() != reflect.Struct {
		return nil, errNoStructPtr(tpe)
	}
	result := Injectable{Type: tpe}
	// Scan each field
	for i := 0; i < tpe.NumField(); i++ {
		structFld := tpe.Field(i)
		// In case we have a tag for the field, get the Qualifier and create a new field injection
		if _, hasTag := structFld.Tag.Lookup(TagKey); hasTag {
			if !structFld.IsExported() {
				return nil, errFieldNotExported(tpe, structFld)
			}
			result.Injections = append(result.Injections, Injection(structFld))
		}
	}
	return &result, nil
}
