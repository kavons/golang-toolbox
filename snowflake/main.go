package main

import (
	"fmt"
	"log"

	"github.com/bwmarrin/snowflake"
)

func main() {
	node, err := snowflake.NewNode(1)
	if err != nil {
		log.Fatal(err)
	}
	var x, y snowflake.ID
	for i := 0; i < 1000000; i++ {
		y = node.Generate()
		if x == y {
			fmt.Errorf("x(%d) & y(%d) are the same", x, y)
		}
		x = y
	}
}