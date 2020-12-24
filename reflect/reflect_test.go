package reflect_test

import (
    "testing"
    "reflect"
    "encoding/json"
)

func TestKind(t *testing.T) {
    memberId := 1
    //memberId := "1"

    var memberInterface interface{}
    memberInterface = memberId
    t.Log(reflect.TypeOf(memberInterface).Kind() == reflect.Int)

    switch reflect.TypeOf(memberInterface).Kind() {
    case reflect.Int:

    case reflect.String:

    }
}

type UserData struct {
    Identity int
}

func TestAssertType(t *testing.T) {
    var d = UserData{
        Identity: 0,
    }

    var di interface{}
    di = d

    if c, ok := di.(UserData); ok {
        t.Log(c, ok)
    }
}

type M struct {
    o int `json:"o"`
}

type SS struct {
    p interface{} `json:"o"`
}

func TestInt(t *testing.T) {
    //var v interface{} = 1000
    //
    //t.Log(reflect.TypeOf(v).Kind())

    var s = M{
        o: 1000,
    }

    b, err := json.Marshal(s)
    if err != nil {
        t.Fatal(err)
    }

    var payload SS
    if err := json.Unmarshal(b, &payload); err != nil {
        t.Fatal(err)
    }

    t.Log(reflect.TypeOf(payload).Kind())
}