package main

import (
	"flag"
	"fmt"
	"github.com/dongri/phonenumber"
)

func main() {
	country := flag.String("country", "", "e.g. US")
	phone := flag.String("phone", "", "e.g. 1112223333")
	flag.Parse()

	if *country == "" || *phone == "" {
		fmt.Println("Usage Example: ./check_phone -country US -phone 1112223333")
		return
	}

	result := phonenumber.Parse(*phone, *country)
	if result == "" {
		fmt.Println("unknown phone number")
	} else {
		fmt.Println("valid phone number")
	}
}
