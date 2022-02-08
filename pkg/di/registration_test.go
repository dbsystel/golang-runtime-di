package di_test

import (
	"errors"
	"reflect"

	"github.com/dbsystel/golang-runtime-di/pkg/di"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("NewRegistration()", func() {
	It("should error on nil", func() {
		_, err := di.NewRegistration(nil, 0)
		Expect(err).To(MatchError(ContainSubstring("invalid component type")))
	})
	It("should error on invalid", func() {
		_, err := di.NewRegistration(uintptr(100), 0)
		Expect(err).To(MatchError(ContainSubstring("invalid component type")))
	})
	Context("instance", func() {
		It("should create a new registration from a struct ptr", func() {
			comp := &ComponentA1{Other: "B"}
			registration, err := di.NewRegistration(comp, 0)
			Expect(registration, err).To(BeAssignableToTypeOf(&di.Registration{}))
			Expect(registration.Type).To(Equal(reflect.TypeOf(&ComponentA1{})))
			Expect(registration.Qualifier).To(BeEmpty())
			Expect(registration.Priority).To(BeNumerically("==", 0))
			Expect(registration.FactoryFn).NotTo(BeNil())
			Expect(registration.Source).To(ContainSubstring("registration_test.go:"))
			instance, _, err := registration.GetInstance()
			Expect(instance, err).To(Equal(comp))
		})
	})
	Context("factory function", func() {
		It("should create a new registration from a struct", func() {
			registration, err := di.NewRegistration(func() *ComponentA1 {
				return &ComponentA1{}
			}, 0)
			Expect(registration, err).To(BeAssignableToTypeOf(&di.Registration{}))
			Expect(registration.Type).To(Equal(reflect.TypeOf(&ComponentA1{})))
			instance, _, err := registration.GetInstance()
			Expect(instance, err).To(BeAssignableToTypeOf(&ComponentA1{}))
		})
		It("should create a new registration from interface", func() {
			registration, err := di.NewRegistration(func() InterfaceA {
				return &ComponentA1{}
			}, 0)
			Expect(registration, err).To(BeAssignableToTypeOf(&di.Registration{}))
			Expect(registration.Type).To(Equal(reflect.TypeOf((*InterfaceA)(nil)).Elem()))
			instance, _, err := registration.GetInstance()
			Expect(instance, err).To(BeAssignableToTypeOf(&ComponentA1{}))
		})
		It("should create a new registration from interface and error", func() {
			registration, err := di.NewRegistration(func() (InterfaceA, error) {
				return &ComponentA1{}, nil
			}, 0)
			Expect(registration, err).To(BeAssignableToTypeOf(&di.Registration{}))
			Expect(registration.Type).To(Equal(reflect.TypeOf((*InterfaceA)(nil)).Elem()))
			instance, _, err := registration.GetInstance()
			Expect(instance, err).To(BeAssignableToTypeOf(&ComponentA1{}))
		})
		It("should return error from factory", func() {
			registration, err := di.NewRegistration(func() (InterfaceA, error) {
				return nil, errors.New("meh")
			}, 0)
			Expect(registration, err).To(BeAssignableToTypeOf(&di.Registration{}))
			Expect(registration.Type).To(Equal(reflect.TypeOf((*InterfaceA)(nil)).Elem()))
			_, _, err = registration.GetInstance()
			Expect(err).To(MatchError(ContainSubstring("meh")))
		})
		It("should error on invalid function", func() {
			_, err := di.NewRegistration(func() {}, 0)
			Expect(err).To(MatchError(ContainSubstring("function should provide 1 or 2 return values")))
		})
		It("should error on invalid in params", func() {
			_, err := di.NewRegistration(func(a int) InterfaceA { return nil }, 0)
			Expect(err).To(MatchError(ContainSubstring("function should not have parameters")))
		})
		It("should error on invalid 1st return type", func() {
			_, err := di.NewRegistration(func() func() { return nil }, 0)
			Expect(err).To(MatchError(ContainSubstring("invalid factory result type")))
		})
		It("should error on invalid 2nd return type", func() {
			_, err := di.NewRegistration(func() (InterfaceA, int) { return nil, 0 }, 0)
			Expect(err).To(MatchError(ContainSubstring("second function return value should be error")))
		})
	})
})

var _ = Describe("Registration", func() {
	var sut *di.Registration
	BeforeEach(func() {
		sut = &di.Registration{
			FactoryFn: func() (interface{}, error) {
				return &ComponentA1{}, nil
			},
			Type:   reflect.TypeOf(&ComponentA1{}),
			Source: "test.go:123",
		}
	})
	Context("GetInstance()", func() {
		It("should return the instance from the factory and cache it", func() {
			instance, first, err := sut.GetInstance()
			Expect(instance, err).To(BeAssignableToTypeOf(&ComponentA1{}))
			Expect(first).To(BeTrue())
			sut.FactoryFn = nil
			instance, first, err = sut.GetInstance()
			Expect(instance, err).To(BeAssignableToTypeOf(&ComponentA1{}))
			Expect(first).To(BeFalse())
		})
		It("should return error from factory", func() {
			errMsg := "meh"
			sut.FactoryFn = func() (interface{}, error) {
				return nil, errors.New(errMsg)
			}
			instance, _, err := sut.GetInstance()
			Expect(instance).To(BeNil())
			Expect(err).To(MatchError(And(
				ContainSubstring("could not create instance"),
				ContainSubstring(errMsg),
			)))
		})
	})
	Context("WithQualifier()", func() {
		It("should set the qualifier", func() {
			qualifier := "meh"
			Expect(sut.WithQualifier(qualifier)).To(Equal(sut))
			Expect(sut.Qualifier).To(Equal(qualifier))
		})
	})
	Context("WithPriority()", func() {
		It("should set the priority", func() {
			priority := 666
			Expect(sut.WithPriority(priority)).To(Equal(sut))
			Expect(sut.Priority).To(Equal(priority))
		})
	})
	Context("String()", func() {
		It("should include type and priority", func() {
			Expect(sut.String()).To(Equal(
				"component *di_test.ComponentA1 with priority 0 registered at: test.go:123",
			))
		})
		It("should include qualifier", func() {
			Expect(sut.WithQualifier("meh").String()).To(Equal(
				"component *di_test.ComponentA1(meh) with priority 0 registered at: test.go:123",
			))
		})
	})
})
