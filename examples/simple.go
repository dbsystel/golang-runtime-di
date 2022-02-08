package examples

import (
	"github.com/dbsystel/golang-runtime-di/pkg/di"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// SimpleDependency is the interface for a dependency to be injected
type SimpleDependency interface {
	Woop()
}

// SimpleComponent is the component to be injected
type SimpleComponent struct{ Name string }

func (c *SimpleComponent) Woop() {}

// SimpleConsumer is the consumer for SimpleDependency
type SimpleConsumer struct {
	// Dependency will be injected without qualifier
	Dependency SimpleDependency `inject:""`
}

var _ = Describe("Simple wiring example", func() {
	It("should wire components after priority", func() {
		scope := &di.Scope{}
		scope.MustRegister(&SimpleComponent{"squash"})
		instance := &SimpleConsumer{}
		scope.MustWire(instance)
		Expect(instance).To(Equal(&SimpleConsumer{Dependency: &SimpleComponent{Name: "squash"}}))
	})
})
