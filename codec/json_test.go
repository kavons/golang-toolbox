package codec

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"testing"
)

func TestDecoder(tt *testing.T) {
	const jsonStream = `
    {"Name": "Ed", "Text": "Knock knock."}
    {"Name": "Sam", "Text": "Who's there?"}
    {"Name": "Ed", "Text": "Go fmt."}
    {"Name": "Sam", "Text": "Go fmt who?"}
    {"Name": "Ed", "Text": "Go fmt yourself!"}
`
	type Message struct {
		Name, Text string
	}
	dec := json.NewDecoder(strings.NewReader(jsonStream))

	// read open bracket
	//t, err := dec.Token()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//// t的类型是json.Delim
	//fmt.Printf("%v\n", t)

	// while the array contains values
	for dec.More() {
		var m Message
		// decode an array value (Message)
		err := dec.Decode(&m)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%v: %v\n", m.Name, m.Text)
	}

	// read closing bracket
	//t, err = dec.Token()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//// t的类型是json.Delim
	//fmt.Printf("%v\n", t)
}
