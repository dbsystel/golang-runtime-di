package examples

import (
	"github.com/dbsystel/golang-runtime-di/pkg/di"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// ParentDependency is the interface for a dependency to be injected
type ParentDependency interface {
	Woop()
}

// ParentComponent is the component to be injected
type ParentComponent struct{ Name string }

func (c *ParentComponent) Woop() {}

// ParentConsumer is the consumer for ParentDependency
type ParentConsumer struct {
	// Dependency will be injected without qualifier
	Dependency ParentDependency `inject:""`
}

var _ = Describe("Parent wiring example", func() {
	It("should wire components after priority", func() {
		parent := &di.Scope{}
		scope := &di.Scope{Parent: parent}
		parent.MustRegister(&ParentComponent{"squash"})
		instance := &ParentConsumer{}
		scope.MustWire(instance)
		Expect(instance).To(Equal(&ParentConsumer{Dependency: &ParentComponent{Name: "squash"}}))
	})
})
