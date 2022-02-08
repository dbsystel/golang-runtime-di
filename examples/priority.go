package examples

import (
	"github.com/dbsystel/golang-runtime-di/pkg/di"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// PriorityDependency is the interface for a dependency to be injected
type PriorityDependency interface {
	Woop()
}

// PriorityComponent is the component to be injected
type PriorityComponent struct{ Name string }

func (c *PriorityComponent) Woop() {}

// PriorityConsumer is the consumer for SimpleDependency
type PriorityConsumer struct {
	// Dependency will be injected by priority
	Dependency PriorityDependency `inject:""`
}

var _ = Describe("Priority example", func() {
	It("should wire components after priority", func() {
		scope := &di.Scope{}
		scope.MustRegister(&PriorityComponent{"soccer"})
		// lower value = higher priority
		scope.MustRegister(&PriorityComponent{"squash"}).WithPriority(-1)
		instance := &PriorityConsumer{}
		scope.MustWire(instance)
		Expect(instance).To(Equal(&PriorityConsumer{Dependency: &PriorityComponent{Name: "squash"}}))
	})
})
