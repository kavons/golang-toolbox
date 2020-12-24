package switch_test

import (
	"fmt"
	"testing"
)

func TestSwitch(t *testing.T) {
	action := "submit"
	switch action {
	case "submit":
		fmt.Println("submit")
	case "cancel":
		fmt.Println("cancel")
	default:
		fmt.Println("default")
	}
}
