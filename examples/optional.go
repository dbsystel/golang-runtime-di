package examples

import (
	"github.com/dbsystel/golang-runtime-di/pkg/di"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// OptionalDependency is the interface for a dependency to be injected
type OptionalDependency interface {
	Woop()
}

// OptionalConsumer is the consumer for SimpleDependency
type OptionalConsumer struct {
	// Dependency cannot be resolved
	Dependency OptionalDependency `inject:"optional"`
}

var _ = Describe("Optional example", func() {
	It("should wire components after Optional", func() {
		scope := &di.Scope{}
		instance := &OptionalConsumer{}
		scope.MustWire(instance)
		Expect(instance).To(Equal(&OptionalConsumer{Dependency: nil}))
	})
})
