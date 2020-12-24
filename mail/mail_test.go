package mail_test

import (
	"fmt"
	"regexp"
	"testing"
)

func ValidateEmail(email string) bool {
	pattern := `^[0-9a-zA-Z][_.0-9a-zA-Z]{0,31}@([0-9a-zA-Z][0-9a-zA-Z-]{0,30}[0-9a-zA-Z]\.){1,4}[a-zA-Z]{2,4}$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

func TestEmailFormat(t *testing.T) {
	m1 := "mjsc1023@163.com"
	fmt.Printf("%s: %v\n", m1, ValidateEmail(m1))

	m2 := "echo.songchi@163.com"
	fmt.Printf("%s: %v\n", m2, ValidateEmail(m2))

	m3 := "a@163.com"
	fmt.Printf("%s: %v\n", m3, ValidateEmail(m3))

	m4 := ""
	fmt.Printf("%s: %v\n", m4, ValidateEmail(m4))
}
