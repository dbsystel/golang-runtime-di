package examples_test

import (
	. "github.com/onsi/ginkgo/v2"
	"testing"

	. "github.com/onsi/gomega"
)

func TestExamples(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "golang-runtime-di-examples")
}
