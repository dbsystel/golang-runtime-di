package examples

import (
	"github.com/dbsystel/golang-runtime-di/pkg/di"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// QualifierDependency is the interface for a dependency to be injected
type QualifierDependency interface {
	Woop()
}

// QualifierComponent is the component to be injected
type QualifierComponent struct{ Name string }

func (c *QualifierComponent) Woop() {}

// QualifierConsumer is the consumer for SimpleDependency
type QualifierConsumer struct {
	// Dependency will be injected without qualifier
	Dependency QualifierDependency `inject:"qualifier=squash"`
}

var _ = Describe("Qualifier example", func() {
	It("should wire components after priority", func() {
		scope := &di.Scope{}
		scope.MustRegister(&QualifierComponent{"soccer"})
		scope.MustRegister(&QualifierComponent{"squash"}).WithQualifier("squash")
		instance := &QualifierConsumer{}
		scope.MustWire(instance)
		Expect(instance).To(Equal(&QualifierConsumer{Dependency: &QualifierComponent{Name: "squash"}}))
	})
})
