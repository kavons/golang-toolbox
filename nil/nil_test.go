package nil_test

import (
	"reflect"
	"testing"
)

type PersonInt interface {
	Age()
}

type PersonImpl struct {
	Name string
}

func (p *PersonImpl) Age() {
	return
}

func TestNilCompare(t *testing.T) {
	var m map[string]bool
	t.Log("m == nil:", m == nil)
	t.Log("reflect.ValueOf(m).IsNil():", reflect.ValueOf(m).IsNil())

	var s []string
	t.Log("s == nil:", s == nil)
	t.Log("reflect.ValueOf(s).IsNil():", reflect.ValueOf(s).IsNil())

	var pImpl *PersonImpl
	t.Log("pImpl == nil:", pImpl == nil)
	t.Log("reflect.ValueOf(pImpl).IsNil():", reflect.ValueOf(pImpl).IsNil())

	var pInt PersonInt
	t.Log("pInt == nil:", pInt == nil) // 这里相等
	pInt = pImpl
	t.Log("pInt == nil:", pInt == nil) // 这里不相等
	t.Log("reflect.ValueOf(pInt).IsNil():", reflect.ValueOf(pInt).IsNil())
}
