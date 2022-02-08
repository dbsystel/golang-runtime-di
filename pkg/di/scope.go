package di

import (
	"reflect"

	"github.com/pkg/errors"
)

var _ InstanceResolver = &Scope{}

// Scope is a scope for Registrations which is used to register and wire dependencies
type Scope struct {
	// Parent is the optional parent scope
	Parent        *Scope
	registrations Registrations
}

func (s *Scope) ResolveInstance(tpe reflect.Type, tag TagValue) (reflect.Value, error) {
	nilValue := reflect.ValueOf(nil)
	identifier := tpe.String()
	if len(tag.Qualifier) > 0 {
		identifier += " with qualifier " + tag.Qualifier
	}
	switch tpe.Kind() {
	case reflect.Array, reflect.Slice:
		candidates, err := s.resolveInjections(tpe.Elem(), tag, identifier)
		if err != nil {
			return nilValue, err
		}
		result := reflect.MakeSlice(tpe, candidates.Len(), candidates.Len())
		for idx, candidate := range candidates {
			instance, err := s.wiredInstance(candidate)
			if err != nil {
				return nilValue, err
			}
			result.Index(idx).Set(reflect.ValueOf(instance))
		}
		return result, nil
	default:
		candidates, err := s.resolveInjections(tpe, tag, identifier)
		if err != nil {
			return nilValue, err
		}
		if len(candidates) == 0 {
			return nilValue, nil
		}
		highestPriority := candidates[0].Priority
		candidates = candidates.FilterPriority(highestPriority)
		if len(candidates) > 1 {
			return nilValue, errors.Errorf("multiple candidates with priority %v for %v:\n\t%v",
				highestPriority, identifier, candidates)
		}
		instance, err := s.wiredInstance(candidates[0])
		if err != nil {
			return nilValue, err
		}
		return reflect.ValueOf(instance), nil
	}
}

// Register uses NewRegistration to register a component or factory func
func (s *Scope) Register(valOrFunc interface{}) (*Registration, error) {
	return s.doRegister(valOrFunc)
}

func (s *Scope) doRegister(valOrFunc interface{}) (*Registration, error) {
	registration, err := NewRegistration(valOrFunc, 2) // nolint:gomnd
	if err != nil {
		return nil, err
	}
	s.registrations = append(s.registrations, registration)
	return registration, nil
}

// MustRegister works like, Register but panics on error
func (s *Scope) MustRegister(valOrFunc interface{}) *Registration {
	result, err := s.doRegister(valOrFunc)
	s.panicOnErr(err)
	return result
}

// Wire wires the targets and all dependencies
func (s *Scope) Wire(targets ...interface{}) error {
	for _, target := range targets {
		if err := s.wireSingle(target); err != nil {
			return err
		}
	}
	return nil
}

// MustWire works like, but panics in case of error
func (s *Scope) MustWire(targets ...interface{}) {
	s.panicOnErr(s.Wire(targets...))
}

func (s *Scope) panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

func (s *Scope) wireSingle(target interface{}) error {
	injectable, err := InjectableFrom(reflect.TypeOf(target))
	if err != nil {
		return err
	}
	return injectable.Apply(reflect.ValueOf(target), s)
}

func (s *Scope) wiredInstance(candidate *Registration) (interface{}, error) {
	instance, created, err := candidate.GetInstance()
	if err != nil {
		return nil, err
	}
	if !created || reflect.TypeOf(instance).Kind() != reflect.Ptr {
		return instance, nil
	}
	return instance, s.wireSingle(instance)
}

func (s *Scope) resolveInjections(tpe reflect.Type, tag TagValue, identifier string) (Registrations, error) {
	candidates := s.registrations.FilterCoercible(tpe)
	if !tag.IsAllQualifier() {
		candidates = candidates.FilterQualifier(tag.Qualifier)
	}
	if s.Parent != nil {
		parentTag := *(&tag) // nolint:staticcheck
		parentTag.Required = false
		fromParent, _ := s.Parent.resolveInjections(tpe, parentTag, identifier)
		candidates = append(candidates, fromParent...)
	}
	candidates = candidates.ByPriority()
	if tag.Required && len(candidates) == 0 {
		return nil, errors.Errorf("no candidate found for: %v", identifier)
	}
	return candidates, nil
}
