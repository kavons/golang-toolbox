package sort_test

import (
	"fmt"
	"sort"
	"testing"
)

type Person struct {
	Name string
	Age  int
}

var persons = []Person{
	{"Majun", 31},
	{"Wangtong", 25},
	{"Jiangyuepeng", 3},
	{"Madanbao", 60},
}

// simple sorter
type SimpleSorter []Person

func (s SimpleSorter) Len() int {
	return len(s)
}

func (s SimpleSorter) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s SimpleSorter) Less(i, j int) bool {
	//return s[i].Age < s[j].Age
	return s[i].Name < s[j].Name
}

func TestSimpleSort(t *testing.T) {
	s := SimpleSorter(persons)
	fmt.Println(persons)
	sort.Sort(s)
	fmt.Println(persons)
}

// advanced sorter
type By func(p1, p2 *Person) bool

type AdvancedSorter struct {
	Persons []Person
	SortBy  By
}

func (s AdvancedSorter) Len() int {
	return len(s.Persons)
}

func (s AdvancedSorter) Swap(i, j int) {
	s.Persons[i], s.Persons[j] = s.Persons[j], s.Persons[i]
}

func (s AdvancedSorter) Less(i, j int) bool {
	return s.SortBy(&s.Persons[i], &s.Persons[j])
}

func TestAdvancedSortByAge(t *testing.T) {
	SortByAge := func(p1, p2 *Person) bool {
		return p1.Age < p2.Age
	}

	s := AdvancedSorter{
		persons,
		SortByAge,
	}

	fmt.Println(persons)
	sort.Sort(s)
	fmt.Println(persons)
}

func TestAdvancedSortByName(t *testing.T) {
	SortByName := func(p1, p2 *Person) bool {
		return p1.Name < p2.Name
	}

	s := AdvancedSorter{
		persons,
		SortByName,
	}

	fmt.Println(persons)
	sort.Sort(s)
	fmt.Println(persons)
}
