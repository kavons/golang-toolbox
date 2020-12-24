package page_query_test

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"reflect"
)

type QueryResult struct {
	Total int64       `json:"total"`
	Rows  interface{} `json:"rows"`
}

func (r *QueryResult) PageQuery(
	qs orm.QuerySeter,
	page, pageSize int64,
	list interface{}) (*QueryResult, error) {

	if page <= 1 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize
	total, err := qs.Limit(pageSize, offset).All(list)
	if err != nil && err != orm.ErrNoRows {
		return nil, err
	}

	return &QueryResult{
		Total: total,
		Rows:  list,
	}, nil
}

func (r *QueryResult) QueryFilter(
	c interface{},
	m string) error {

	class := reflect.ValueOf(c)
	method := class.MethodByName(m)
	if method.IsValid() == false {
		err := fmt.Errorf("method not found param name: %s", m)
		return err
	}

	r.Rows = method.Call([]reflect.Value{reflect.ValueOf(r.Rows)})[0].Interface()
	return nil
}

func (r *QueryResult) PageQueryWithFilter(
	qs orm.QuerySeter,
	page, pageSize int,
	list interface{},
	c interface{},
	m string) (*QueryResult, error) {

	if page <= 1 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize
	total, err := qs.Limit(pageSize, offset).All(list)
	if err != nil && err != orm.ErrNoRows {
		return nil, err
	}

	result := &QueryResult{
		Total: total,
		Rows:  list,
	}

	if c != nil && m != "" {
		err := result.QueryFilter(c, m)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}
