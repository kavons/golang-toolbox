package validator_test

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/shopspring/decimal"
)

func WonValidator(v reflect.Value) error {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)

        switch f.Type() {
        case reflect.TypeOf(decimal.Decimal{}):
            tag := string(t.Field(i).Tag.Get("decimal"))
            if tag == "" {
                continue
            }

            vals := strings.SplitN(tag, "=", 2)
            if len(vals) != 2 {
                continue
            }

            value := f.Interface().(decimal.Decimal)
            benchmark, err := decimal.NewFromString(vals[1])
            if err != nil {
                continue
            }

            switch vals[0] {
            case "equal":
                if !value.Equal(benchmark) {
                    return errors.New(fmt.Sprintf("%s-%s: %s must be %s", t.Name(), t.Field(i).Name, value.String(), benchmark.String()))
                }
            case "lt":
                if value.GreaterThanOrEqual(benchmark) {
                    return errors.New(fmt.Sprintf("%s-%s: %s must be less than %s", t.Name(), t.Field(i).Name, value.String(), benchmark.String()))
                }
            case "lte":
                if value.GreaterThan(benchmark) {
                    return errors.New(fmt.Sprintf("%s-%s: %s must be %s or less", t.Name(), t.Field(i).Name, value.String(), benchmark.String()))
                }
            case "gt":
                if value.LessThanOrEqual(benchmark) {
                    return errors.New(fmt.Sprintf("%s-%s: %s must be greater than %s", t.Name(), t.Field(i).Name, value.String(), benchmark.String()))
                }
            case "gte":
                if value.LessThan(benchmark) {
                    return errors.New(fmt.Sprintf("%s-%s: %s must be %s or greater", t.Name(), t.Field(i).Name, value.String(), benchmark.String()))
                }
            }
        }
	}

	return nil
}

func TestEqual(t *testing.T) {
    type TestStruct struct {
        Value decimal.Decimal `decimal:"equal=0"`
    }

    v := TestStruct{
        Value: decimal.NewFromFloat(0.1),
    }

    err := WonValidator(reflect.ValueOf(&v))
    t.Log(err.Error())
}

func TestLt(t *testing.T) {
    type TestStruct struct {
        Value decimal.Decimal `decimal:"lt=0"`
    }

    v := TestStruct{
        Value: decimal.NewFromFloat(0.0),
    }

    err := WonValidator(reflect.ValueOf(&v))
    t.Log(err.Error())
}

func TestLte(t *testing.T) {
    type TestStruct struct {
        Value decimal.Decimal `decimal:"lte=0"`
    }

    v := TestStruct{
        Value: decimal.NewFromFloat(0.5),
    }

    err := WonValidator(reflect.ValueOf(&v))
    t.Log(err.Error())
}

func TestGt(t *testing.T) {
    type TestStruct struct {
        Value decimal.Decimal `decimal:"gt=0"`
    }

    v := TestStruct{
        Value: decimal.NewFromFloat(0.0),
    }

    err := WonValidator(reflect.ValueOf(&v))
    t.Log(err.Error())
}

func TestGte(t *testing.T) {
    type Gte struct {
        Value decimal.Decimal `decimal:"gte=0"`
    }

    v := Gte{
        Value: decimal.NewFromFloat(-0.5),
    }

    err := WonValidator(reflect.ValueOf(&v))
    t.Log(err.Error())
}