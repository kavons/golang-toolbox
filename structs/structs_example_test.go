package structs

import (
	"fmt"
	"testing"
	"time"

	"github.com/fatih/structs"
)

func TestNew(t *testing.T) {
	type Server struct {
		Name    string
		ID      int32
		Enabled bool
	}

	server := &Server{
		Name:    "Arslan",
		ID:      123456,
		Enabled: true,
	}

	s := structs.New(server)

	fmt.Printf("Name        : %v\n", s.Name())
	fmt.Printf("Values      : %v\n", s.Values())
	fmt.Printf("Value of ID : %v\n", s.Field("ID").Value())
	// Output:
	// Name        : Server
	// Values      : [Arslan 123456 true]
	// Value of ID : 123456

}

func TestExampleMap(t *testing.T) {
	type Server struct {
		Name    string
		ID      int32
		Enabled bool
	}

	s := &Server{
		Name:    "Arslan",
		ID:      123456,
		Enabled: true,
	}

	m := structs.Map(s)

	fmt.Printf("%#v\n", m["Name"])
	fmt.Printf("%#v\n", m["ID"])
	fmt.Printf("%#v\n", m["Enabled"])
	// Output:
	// "Arslan"
	// 123456
	// true

}

func TestMap_tags(t *testing.T) {
	// Custom tags can change the map keys instead of using the fields name
	type Server struct {
		Name    string `structs:"server_name"`
		ID      int32  `structs:"server_id"`
		Enabled bool   `structs:"enabled"`
	}

	s := &Server{
		Name: "Zeynep",
		ID:   789012,
	}

	m := structs.Map(s)

	// access them by the custom tags defined above
	fmt.Printf("%#v\n", m["server_name"])
	fmt.Printf("%#v\n", m["server_id"])
	fmt.Printf("%#v\n", m["enabled"])
	// Output:
	// "Zeynep"
	// 789012
	// false

}

func TestMap_omitNested(t *testing.T) {
	type Date struct {
		Year  int `structs:"year"`
		Month int `structs:"month"`
		Day   int `structs:"day"`
	}

	// By default field with struct types are processed too. We can stop
	// processing them via "omitnested" tag option.
	type Server struct {
		Name string `structs:"server_name"`
		ID   int32  `structs:"server_id"`
		D    Date   `structs:"date,omitnested"` // do not convert to map[string]interface{}
		//D Date `structs:"date"` // convert to map[string]interface{}
	}

	s := &Server{
		Name: "Zeynep",
		ID:   789012,
		D:    Date{2019, 1, 31},
	}

	m := structs.Map(s)

	// access them by the custom tags defined above
	fmt.Printf("%v\n", m["server_name"])
	fmt.Printf("%v\n", m["server_id"])
	fmt.Printf("%#v\n", m["date"])
	// Output:
	// Zeynep
	// 789012
	// structs.Date{Year:2019, Month:1, Day:31}
	// map[string]interface {}{"day":31, "year":2019, "month":1}
}

func TestMap_omitEmpty(t *testing.T) {
	// By default field with struct types of zero values are processed too. We
	// can stop processing them via "omitempty" tag option.
	type Server struct {
		Name     string `structs:",omitempty"`
		ID       int32  `structs:"server_id,omitempty"`
		Location string
	}

	// Only add location
	s := &Server{
		Location: "Tokyo",
	}

	m := structs.Map(s)

	// map contains only the Location field
	fmt.Printf("%v\n", m)
	// Output:
	// map[Location:Tokyo]
}

func TestExampleValues(t *testing.T) {
	type Server struct {
		Name    string
		ID      int32
		Enabled bool
	}

	s := &Server{
		Name:    "Fatih",
		ID:      135790,
		Enabled: false,
	}

	m := structs.Values(s)

	fmt.Printf("Values: %v\n", m)
	// Output:
	// Values: [Fatih 135790 false]
}

func TestValues_omitEmpty(t *testing.T) {
	// By default field with struct types of zero values are processed too. We
	// can stop processing them via "omitempty" tag option.
	type Server struct {
		Name     string `structs:",omitempty"`
		ID       int32  `structs:"server_id,omitempty"`
		Location string
	}

	// Only add location
	s := &Server{
		Location: "Ankara",
	}

	m := structs.Values(s)

	// values contains only the Location field
	fmt.Printf("Values: %v\n", m)
	// Output:
	// Values: [Ankara]
}

func TestValues_tags(t *testing.T) {
	type Location struct {
		City    string
		Country string
	}

	type Server struct {
		Name     string
		ID       int32
		Enabled  bool
		Location Location `structs:"-"` // values from location are not included anymore
	}

	s := &Server{
		Name:     "Fatih",
		ID:       135790,
		Enabled:  false,
		Location: Location{City: "Ankara", Country: "Turkey"},
	}

	// Let get all values from the struct s. Note that we don't include values
	// from the Location field
	m := structs.Values(s)

	fmt.Printf("Values: %+v\n", m)
	// Output:
	// Values: [Fatih 135790 false]
}

func TestExampleFields(t *testing.T) {
	type Access struct {
		Name         string
		LastAccessed time.Time
		Number       int
	}

	s := &Access{
		Name:         "Fatih",
		LastAccessed: time.Now(),
		Number:       1234567,
	}

	fields := structs.Fields(s)

	for i, field := range fields {
		fmt.Printf("[%d] %+v\n", i, field.Name())
	}

	// Output:
	// [0] Name
	// [1] LastAccessed
	// [2] Number
}

func TestFields_nested(t *testing.T) {
	type Person struct {
		Name   string
		Number int
	}

	type Access struct {
		Person        Person
		HasPermission bool
		LastAccessed  time.Time
	}

	s := &Access{
		Person:        Person{Name: "fatih", Number: 1234567},
		LastAccessed:  time.Now(),
		HasPermission: true,
	}

	// Let's get all fields from the struct s.
	fields := structs.Fields(s)

	for _, field := range fields {
		if field.Name() == "Person" {
			fmt.Printf("Access.Person.Name: %+v\n", field.Field("Name").Value())
		}
	}

	// Output:
	// Access.Person.Name: fatih
}

func TestField(t *testing.T) {
	type Person struct {
		Name   string
		Number int
	}

	type Access struct {
		Person        Person
		HasPermission bool
		LastAccessed  time.Time
	}

	access := &Access{
		Person:        Person{Name: "fatih", Number: 1234567},
		LastAccessed:  time.Now(),
		HasPermission: true,
	}

	// Create a new Struct type
	s := structs.New(access)

	// Get the Field type for "Person" field
	p := s.Field("Person")

	// Get the underlying "Name field" and print the value of it
	name := p.Field("Name")

	fmt.Printf("Value of Person.Access.Name: %+v\n", name.Value())

	// Output:
	// Value of Person.Access.Name: fatih

}

func TestExampleIsZero(t *testing.T) {
	type Server struct {
		Name    string
		ID      int32
		Enabled bool
	}

	// Nothing is initialized
	a := &Server{}
	isZeroA := structs.IsZero(a)

	// Name and Enabled is initialized, but not ID
	b := &Server{
		Name:    "Golang",
		Enabled: true,
	}
	isZeroB := structs.IsZero(b)

	fmt.Printf("%#v\n", isZeroA)
	fmt.Printf("%#v\n", isZeroB)
	// Output:
	// true
	// false
}

func TestExampleHasZero(t *testing.T) {
	// Let's define an Access struct. Note that the "Enabled" field is not
	// going to be checked because we added the "structs" tag to the field.
	type Access struct {
		Name         string
		LastAccessed time.Time
		Number       int
		Enabled      bool `structs:"-"`
	}

	// Name and Number is not initialized.
	a := &Access{
		LastAccessed: time.Now(),
	}
	hasZeroA := structs.HasZero(a)

	// Name and Number is initialized.
	b := &Access{
		Name:         "Fatih",
		LastAccessed: time.Now(),
		Number:       12345,
	}
	hasZeroB := structs.HasZero(b)

	fmt.Printf("%#v\n", hasZeroA)
	fmt.Printf("%#v\n", hasZeroB)
	// Output:
	// true
	// false
}
