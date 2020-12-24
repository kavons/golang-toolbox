package uuid

import (
	"fmt"
	"testing"

	"github.com/satori/go.uuid"
)

func TestCreateUUID(t *testing.T) {
	fmt.Println(uuid.NewV4().String())
}
