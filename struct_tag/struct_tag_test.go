package struct_tag_test

import (
	"fmt"
	"reflect"
	"testing"
)

func TestReflectTag(t *testing.T) {
	type User struct {
		UserId   int    `json:"user_id"`
		UserName string `json:"user_name"`
	}

	u := &User{UserId: 1, UserName: "Kobe"}
	tt := reflect.TypeOf(u)
	field := tt.Elem().Field(0)
	fmt.Println(field.Name)
	fmt.Println(field.Tag.Get("json"))

	ss := tt.Elem()
	for i := 0; i < ss.NumField(); i++ {
		fmt.Println(ss.Field(i).Tag)
	}
}
