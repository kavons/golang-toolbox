package _struct

import (
	"fmt"
	"testing"
)

type Name struct {
	name string
}

type Dog struct {
	name  Name
	color string
	age   int8
	kind  string
}

func TestStructValueCopy(t *testing.T) {
	d1 := Dog{Name{"豆豆"}, "黑色", 2, "二哈"}
	fmt.Printf("d1: %T , %v , %p \n", d1, d1, &d1)
	d2 := d1 //值拷贝
	fmt.Printf("d2: %T , %v , %p \n", d2, d2, &d2)

	d2.name.name = "毛毛"
	fmt.Println("d2修改后：", d2)
	fmt.Println("d1：", d1)
	fmt.Println("------------------")
}

func TestStructPointerCopy(t *testing.T) {
	d1 := Dog{Name{"豆豆"}, "黑色", 2, "二哈"}
	d3 := &d1
	fmt.Printf("d3: %T , %v , %p \n", d3, d3, d3)
	d3.name.name = "球球"
	d3.color = "白色"
	d3.kind = "萨摩耶"
	fmt.Println("d3修改后：", d3)
	fmt.Println("d1：", d1)
	fmt.Println("------------------")
}
