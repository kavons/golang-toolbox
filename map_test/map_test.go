package map_test

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
	"utils/.vendor-new/github.com/pkg/errors"
)

func TestMap(t *testing.T) {
	subscriptions := make(map[string]map[string]bool)

	channel := make(map[string]bool)
	channel["event"] = true
	subscriptions["channel"] = channel

	if channel, ok := subscriptions["channel"]; ok {
		fmt.Println(len(channel))
		channel["event"] = false
	}

	fmt.Printf("%#v", subscriptions)
}

type MapStruct struct {
	Progress map[string]bool
}

func TestMap2(t *testing.T) {
	m := &MapStruct{}
	s := `{"sms": true}`
	if _, ok := m.Progress["sms"]; !ok {
		fmt.Println("bad")
	}
	err := json.Unmarshal([]byte(s), &m.Progress)
	fmt.Println(err, m.Progress)
	m.Progress = make(map[string]bool)
	m.Progress["sms"] = true

	fmt.Printf("%#v", m)
}

func TestMap3(t *testing.T) {
	var nilMap map[string]bool

	for k := range nilMap {
		fmt.Println(k)
	}
}

func Same() (count int) {
	fmt.Println(&count)
	count, err := 10, errors.New("test")
	fmt.Println(&count, err)
	return
}

func TestVar(t *testing.T) {
	fmt.Println(Same())

	a := "ma#jun"
	fmt.Println("a:", strings.Index(a, "#"), a[1:], a)
	for b, c := range a {
		fmt.Println(b, c)
	}

	d := map[string]interface{}{}
	if _, ok := d["p"]; ok {
		fmt.Println("true")
	} else {
		d["p"] = "yes"
	}
	fmt.Println(len(d))
}

func Test4Test(t *testing.T) {
	fileInfo, err := os.Stat("./conf/locales/en")
	os.IsNotExist(err)
	fmt.Println(fileInfo, err)
}

type Member struct {
	Id int
}

func TestMapList(t *testing.T) {
	var members []*Member
	members = append(members, &Member{
		Id: 1,
	}, &Member{
		Id: 2,
	})

	memberMap := make(map[int][]*Member)
	for _, m := range members {
		memberMap[m.Id] = append(memberMap[m.Id], m)
	}

	//fmt.Printf("%#v", members)
	fmt.Printf("%#v", memberMap)
}

func TestStrList(t *testing.T) {
	l := []string{
		"a",
		"b",
		"c",
	}

	for i, s := range l {
		l[i] = fmt.Sprintf("\"%s\"", s)
	}
	fmt.Println(strings.Join(l, ","))

	a := strings.Split("", ",")
	fmt.Println(a, len(a), a[0])

	var b []string
	fmt.Println(b, len(b), b[0])
}

func TestShadeVar(t *testing.T) {
	var o = 1
	var r = 2
	{
		o, r := 3, 4
		fmt.Println(o, r)
	}
	fmt.Println(o, r)
}

func TestSwitch(t *testing.T) {
	a := 1
	switch a {
	case 1:
		fmt.Println(1)
	default:
		fmt.Println("default")
	}
}

type User struct {
	name string
	attr Attr
}

type Attr struct {
	age     int
	address string
}

func TestNestedUser(t *testing.T) {
	u := User{
		name: "majun",
		attr: Attr{
			age:     10,
			address: "china",
		},
	}

	m := u
	fmt.Printf("%p, %p\n", &u.attr, &m.attr)

	fmt.Println(u.name)
	fmt.Printf("%#v\n", u.attr)

	m.name = "songchi"
	m.attr.age = 20
	fmt.Println(u.name)
	fmt.Printf("%#v\n", u.attr)
}

func TestMakeMap(t *testing.T) {
	m := make(map[string]int, 10)
	m["first"] = 1
	t.Log(len(m))
}
