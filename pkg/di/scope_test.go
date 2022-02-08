package di_test

import (
	"errors"
	"reflect"

	"github.com/dbsystel/golang-runtime-di/pkg/di"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("DefaultScope", func() {
	var sut *di.Scope
	BeforeEach(func() {
		sut = &di.Scope{}
	})
	Context("Register()", func() {
		It("should register a component", func() {
			Expect(sut.Register(&ComponentA2{})).NotTo(BeNil())
			Expect(sut.ResolveInstance(reflect.TypeOf(&ComponentA2{}), di.TagValue{Required: false})).NotTo(
				WithTransform(func(value reflect.Value) interface{} { return value.Interface() }, BeNil()),
			)
		})
	})
	Context("MustRegister()", func() {
		It("should panic on error from registration", func() {
			Expect(func() { sut.MustRegister(uintptr(0)) }).To(Panic())
		})
	})
	Context("MustWire()", func() {
		It("should wire an instance", func() {
			sut.MustRegister(ValueA("a"))
			sut.MustRegister(func() InterfaceA { return &ComponentA1{Other: "x"} })
			instance := &ComponentB1{}
			sut.MustWire(instance)
			Expect(instance).To(Equal(&ComponentB1{A: &ComponentA1{A: "a", Other: "x"}}))
		})
		It("inherit from parent", func() {
			sut.Parent = &di.Scope{}
			sut.Parent.MustRegister(ValueA("a"))
			instance := &ComponentA1{}
			sut.MustWire(instance)
			Expect(instance).To(Equal(&ComponentA1{A: "a"}))
		})
		It("should wire an instance with qualifier and priority", func() {
			sut.MustRegister(ValueA("b")).WithQualifier("a")
			sut.MustRegister(ValueA("a")).WithQualifier("a").WithPriority(-1)
			sut.MustRegister(func() InterfaceA { return &ComponentA1{Other: "x"} }).WithPriority(1)
			sut.MustRegister(func() InterfaceA { return &ComponentA2{} })
			instance := &ComponentB1{}
			sut.MustWire(instance)
			Expect(instance).To(Equal(&ComponentB1{A: &ComponentA2{A: "a"}}))
		})
		It("should wire all known", func() {
			sut.MustRegister(ValueA("b")).WithQualifier("a")
			sut.MustRegister(ValueA("a"))
			instance := &AllValueA{}
			sut.MustWire(instance)
			Expect(instance).To(Equal(&AllValueA{
				Values: []ValueA{"a", "b"},
			}))
		})
		It("should not error on optional missing", func() {
			instance := &ComponentA2{}
			sut.MustWire(instance)
			Expect(instance).To(Equal(&ComponentA2{}))
		})
	})
	Context("Wire()", func() {
		It("should return error from factory", func() {
			sut.MustRegister(func() (ValueA, error) {
				return "a", errors.New("meh")
			})
			Expect(sut.Wire(&ComponentA1{})).To(MatchError(ContainSubstring("meh")))
		})
		It("should error if invalid type", func() {
			Expect(sut.Wire(ValueA(""))).To(HaveOccurred())
		})
		It("should error if dependency not found", func() {
			Expect(sut.Wire(&ComponentA1{})).To(HaveOccurred())
		})
		It("should error if on required all known", func() {
			Expect(sut.Wire(&AllValueA{})).To(HaveOccurred())
		})
		It("should error if all known factory fails", func() {
			sut.MustRegister(func() (ValueA, error) {
				return "a", errors.New("meh")
			})
			Expect(sut.Wire(&AllValueA{})).To(MatchError(ContainSubstring("meh")))
		})
		It("should error on multiple candidates", func() {
			sut.MustRegister(ValueA("a"))
			sut.MustRegister(ValueA("b"))
			Expect(sut.Wire(&ComponentA1{})).To(HaveOccurred())
		})
	})

})
