package beego_orm

import (
	"testing"

	"won-utils"
    "won-models"
    "time"
    "strings"
)

func TestBirthDate(t *testing.T) {
	idd, _ := models.GetIdDocumentByMemberId(12)
    t.Log(idd.BirthDate)
	birthDate := idd.BirthDate.Format("2006-01-02")
	t.Log(birthDate)

    b, err := time.ParseInLocation(utils.DateFormat, birthDate, time.UTC)
    if err != nil {
        t.Fatal(err)
    }
    t.Log(b)

    b2 := b.Format("2006-01-02")
    t.Log(b2)
}

func TestTrim(t *testing.T) {
    s := " 123 "
    d := strings.Trim(s, " ")
    t.Log(s, len(s))
    t.Log(d, len(d))
}