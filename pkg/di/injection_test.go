package di_test

import (
	"reflect"

	"github.com/dbsystel/golang-runtime-di/pkg/di"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Injection", func() {
	value := ValueA("A")
	var tgtA ComponentA1
	var tgtB ComponentB1
	var sut di.Injection
	BeforeEach(func() {
		tgtA = ComponentA1{}
		tgtB = ComponentB1{}
		sut = di.Injection(reflect.TypeOf(tgtA).Field(1))
	})
	Context("Apply()", func() {
		reflectValValue := reflect.ValueOf(value)
		It("should inject into struct ptr", func() {
			Expect(sut.Apply(reflect.ValueOf(&tgtA), reflectValValue)).NotTo(HaveOccurred())
			Expect(tgtA.A).To(Equal(value))
		})
		It("should inject interface", func() {
			Expect(sut.Apply(reflect.ValueOf(&tgtB), reflect.ValueOf(&tgtA))).NotTo(HaveOccurred())
			Expect(tgtB.A).To(Equal(&tgtA))
		})
		It("should error if not struct ptr", func() {
			Expect(sut.Apply(reflect.ValueOf(tgtA), reflectValValue)).To(
				MatchError(ContainSubstring("expected a struct pointer")),
			)
		})
		It("should error if field does not exist on target", func() {
			Expect(sut.Apply(reflect.ValueOf(&ComponentC{}), reflectValValue)).To(
				MatchError(ContainSubstring("field 'A' is not valid in target")),
			)
		})
		It("should error, if not coercible", func() {
			Expect(sut.Apply(reflect.ValueOf(&tgtB), reflectValValue)).To(
				MatchError(ContainSubstring("cannot coerce ")),
			)
		})
	})
})
