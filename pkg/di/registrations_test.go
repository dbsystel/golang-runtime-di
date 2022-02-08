package di_test

import (
	"reflect"
	"strings"

	"github.com/dbsystel/golang-runtime-di/pkg/di"
	. "github.com/onsi/ginkgo/v2"

	. "github.com/onsi/gomega"
)

var _ = Describe("Registrations", func() {
	var sut di.Registrations
	BeforeEach(func() {
		sut = di.Registrations{
			{Type: reflect.TypeOf(&ComponentA1{}), Qualifier: ""},
			{Type: reflect.TypeOf(&ComponentA1{}), Qualifier: "a"},
			{Type: reflect.TypeOf((*InterfaceA)(nil)).Elem(), Qualifier: "", Priority: -1},
		}
	})
	Context("ByPriority()", func() {
		It("should sort correctly", func() {
			Expect(sut.ByPriority()).To(Equal(di.Registrations{sut[2], sut[0], sut[1]}))
		})
	})
	Context("FilterQualifier()", func() {
		It("should filter correctly", func() {
			Expect(sut.FilterQualifier("a")).To(ConsistOf(sut[1]))
		})
	})
	Context("FilterPriority()", func() {
		It("should filter correctly", func() {
			Expect(sut.FilterPriority(-1)).To(ConsistOf(sut[2]))
		})
	})
	Context("FilterCoercible()", func() {
		It("should filter correctly by ptr", func() {
			Expect(sut.FilterCoercible(sut[0].Type)).To(Equal(di.Registrations{sut[0], sut[1]}))
		})
		It("should filter correctly by ptr", func() {
			Expect(sut.FilterCoercible(sut[2].Type)).To(Equal(sut))
		})
	})
	Context("String()", func() {
		It("should build correctly", func() {
			Expect(strings.Split(sut.String(), "\n\t")).To(HaveLen(sut.Len()))
		})
	})
})
