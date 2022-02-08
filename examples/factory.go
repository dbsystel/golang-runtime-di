package examples

import (
	"github.com/dbsystel/golang-runtime-di/pkg/di"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// FactoryDependency is the interface for a dependency to be injected
type FactoryDependency interface {
	Woop()
}

// FactoryComponent is the component to be injected
type FactoryComponent struct{ Name string }

func (c *FactoryComponent) Woop() {}

// FactoryConsumer is the consumer for SimpleDependency
type FactoryConsumer struct {
	// Dependency will be injected by Factory
	Dependency FactoryDependency `inject:""`
}

var _ = Describe("Factory example", func() {
	It("should wire components after Factory", func() {
		scope := &di.Scope{}
		// KINDLY NOTE:
		// - the return type may also be *FactoryComponent
		// - the error return parameter is optional for the factory function
		scope.MustRegister(func() (FactoryDependency, error) {
			return &FactoryComponent{"squash"}, nil
		})
		instance := &FactoryConsumer{}
		scope.MustWire(instance)
		Expect(instance).To(Equal(&FactoryConsumer{Dependency: &FactoryComponent{Name: "squash"}}))
	})
})
