package di_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"

	. "github.com/onsi/gomega"
)

func TestDI(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "golang-runtime-di")
}

type ValueA string
type ValueB string

type InvalidComponent struct {
	meh bool `inject:""` // nolint: structcheck,unused
}

type InterfaceA interface {
	GetA() string
}

type ComponentA1 struct {
	B     ValueB `inject:"optional"`
	A     ValueA `inject:""`
	Other string
}

func (c *ComponentA1) GetA() string { return string(c.A) }

type ComponentA2 struct {
	A ValueA `inject:"qualifier=a,optional"`
}

func (c *ComponentA2) GetA() string { return string(c.A) }

type InterfaceB interface {
	GetB() string
}

type ComponentB1 struct {
	A InterfaceA `inject:""`
}

func (c *ComponentB1) GetB() string {
	return c.A.GetA()
}

type ComponentC struct {
	B InterfaceB `inject:"b2"`
}

type AllValueA struct {
	Values []ValueA `inject:"qualifier=*"`
}
