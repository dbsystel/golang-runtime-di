package examples

import (
	"github.com/dbsystel/golang-runtime-di/pkg/di"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// AllDependency is the interface for a dependency to be injected
type AllDependency interface {
	Woop()
}

// AllComponent is the component to be injected
type AllComponent struct{ Name string }

func (c *AllComponent) Woop() {}

// AllConsumer is the consumer for SimpleDependency
type AllConsumer struct {
	// Dependencies will be injected with all known AllDependency (qualifier = *)
	Dependencies []AllDependency `inject:"qualifier=*"`
}

var _ = Describe("All example", func() {
	It("should wire components after All", func() {
		scope := &di.Scope{}
		scope.MustRegister(&AllComponent{"soccer"}).WithQualifier("soccer")
		scope.MustRegister(&AllComponent{"squash"})
		instance := &AllConsumer{}
		scope.MustWire(instance)
		Expect(instance).To(Equal(&AllConsumer{Dependencies: []AllDependency{
			&AllComponent{"squash"},
			&AllComponent{"soccer"},
		}}))
	})
})
