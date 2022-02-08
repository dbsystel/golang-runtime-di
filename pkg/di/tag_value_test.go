package di_test

import (
	"github.com/dbsystel/golang-runtime-di/pkg/di"
	. "github.com/onsi/ginkgo/v2"

	. "github.com/onsi/gomega"
)

var _ = Describe("TagValueFrom()", func() {
	It("should parse empty", func() {
		Expect(di.TagValueFrom("")).To(Equal(di.TagValue{Required: true, Qualifier: ""}))
	})
	It("should parse optional", func() {
		Expect(di.TagValueFrom("optional")).To(Equal(di.TagValue{Required: false, Qualifier: ""}))
	})
	It("should parse all options, skipping invalid", func() {
		Expect(di.TagValueFrom("qualifier=meh,optional,meh")).To(Equal(di.TagValue{Required: false, Qualifier: "meh"}))
	})
	It("should parse all qualifiers selector", func() {
		Expect(di.TagValueFrom("qualifier=*").IsAllQualifier()).To(BeTrue())
	})
})
