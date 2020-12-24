package rand_test

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestRand(t *testing.T) {
	code := rand.Intn(999)
	fmt.Println(code)
}
