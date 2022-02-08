package di

import "reflect"

// InstanceResolver is the interface for resolving instances for a reflect.Type
type InstanceResolver interface {
	// ResolveInstance resolves the reflect.Value for the provided reflect.Type and qualifier
	ResolveInstance(tpe reflect.Type, tag TagValue) (reflect.Value, error)
}
