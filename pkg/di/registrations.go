package di

import (
	"reflect"
	"sort"
	"strings"
)

type Registrations []*Registration

func (r Registrations) String() string {
	var sb strings.Builder
	for idx, reg := range r {
		if idx > 0 {
			sb.WriteString("\n\t")
		}
		sb.WriteString(reg.String())
	}
	return sb.String()
}

func (r Registrations) Len() int {
	return len(r)
}

func (r Registrations) Less(i, j int) bool {
	if len(r[i].Qualifier) == 0 && len(r[j].Qualifier) > 0 {
		return true
	}
	return r[i].Priority < r[j].Priority
}

func (r Registrations) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r Registrations) FilterCoercible(target reflect.Type) Registrations {
	if target.Kind() == reflect.Interface {
		return r.filter(func(reg *Registration) bool { return reg.Type.Implements(target) })
	}
	return r.filter(func(reg *Registration) bool { return reg.Type.AssignableTo(target) })
}

func (r Registrations) FilterQualifier(qualifier string) Registrations {
	return r.filter(func(reg *Registration) bool {
		return reg.Qualifier == qualifier
	})
}

func (r Registrations) FilterPriority(priority int) Registrations {
	return r.filter(func(reg *Registration) bool {
		return reg.Priority == priority
	})
}

func (r Registrations) ByPriority() Registrations {
	result := make(Registrations, len(r))
	copy(result, r)
	sort.Sort(result)
	return result
}

func (r Registrations) filter(f func(reg *Registration) bool) (result Registrations) {
	for _, reg := range r {
		if f(reg) {
			result = append(result, reg)
		}
	}
	return result
}
