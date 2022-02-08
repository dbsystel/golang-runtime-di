package di_test

import (
	"errors"
	"reflect"

	"github.com/dbsystel/golang-runtime-di/pkg/di"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ di.InstanceResolver = testResolver{}

type testResolver struct {
	err     error
	value   reflect.Value
	fbValue reflect.Value
}

func (r testResolver) ResolveInstance(tpe reflect.Type, _ di.TagValue) (reflect.Value, error) {
	if r.value.IsValid() && r.value.Type().AssignableTo(tpe) {
		return r.value, r.err
	}
	return r.fbValue, r.err
}

var _ = Describe("InjectableFrom()", func() {
	It("should error on nil", func() {
		_, err := di.InjectableFrom(nil)
		Expect(err).To(HaveOccurred())
	})
	It("should return valid injectable for struct", func() {
		res, err := di.InjectableFrom(reflect.TypeOf(ComponentA1{}))
		Expect(res, err).To(BeAssignableToTypeOf(&di.Injectable{}))
		Expect(res.Injections).To(HaveLen(2))
	})
	It("should return valid injectable for struct ptr", func() {
		res, err := di.InjectableFrom(reflect.TypeOf(&ComponentA1{}))
		Expect(res, err).To(BeAssignableToTypeOf(&di.Injectable{}))
		Expect(res.Type)
		Expect(res.Injections).To(HaveLen(2))
	})
	It("should return valid injectable for interface struct ptr", func() {
		res, err := di.InjectableFrom(reflect.TypeOf((InterfaceA)(&ComponentA1{})))
		Expect(res, err).To(BeAssignableToTypeOf(&di.Injectable{}))
		Expect(res.Injections).To(HaveLen(2))
	})
	It("should return error if not struct ptr", func() {
		_, err := di.InjectableFrom(reflect.TypeOf((InterfaceA)(nil)))
		Expect(err).To(MatchError(ContainSubstring("expected a struct pointer")))
	})
	It("should return error if field not exported", func() {
		_, err := di.InjectableFrom(reflect.TypeOf(InvalidComponent{}))
		Expect(err).To(MatchError(ContainSubstring("field not exported")))
	})
})

var _ = Describe("Injectable", func() {
	var resolver testResolver
	var sut *di.Injectable
	var tgt ComponentA1
	valueA := ValueA("meh")
	BeforeEach(func() {
		var err error
		tgt = ComponentA1{Other: "b"}
		sut, err = di.InjectableFrom(reflect.TypeOf(&tgt))
		Expect(err).NotTo(HaveOccurred())
		resolver = testResolver{}
	})
	Context("Apply()", func() {
		It("should inject values", func() {
			resolver.value = reflect.ValueOf(valueA)
			Expect(sut.Apply(reflect.ValueOf(&tgt), resolver)).NotTo(HaveOccurred())
			Expect(tgt).To(Equal(ComponentA1{A: valueA, Other: "b"}))
		})
		It("should return error if not ptr", func() {
			Expect(sut.Apply(reflect.ValueOf(tgt), resolver)).To(MatchError(
				ContainSubstring("expected a struct pointer"),
			))
		})
		It("should return error if invalid type", func() {
			Expect(sut.Apply(reflect.ValueOf(&ComponentA2{}), resolver)).To(MatchError(
				ContainSubstring("cannot coerce"),
			))
		})
		It("should return error from resolver", func() {
			errMsg := "mehmeh"
			resolver.err = errors.New(errMsg)
			Expect(sut.Apply(reflect.ValueOf(&tgt), resolver)).To(MatchError(
				ContainSubstring(errMsg),
			))
		})
		It("should return error from injection", func() {
			resolver.fbValue = reflect.ValueOf("mehmehmeh")
			Expect(sut.Apply(reflect.ValueOf(&tgt), resolver)).To(MatchError(
				ContainSubstring("could not inject field"),
			))
		})
	})
})
