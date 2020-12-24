package request_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"context"

	"github.com/gemalto/requester"
)

type Response struct {
	Error     string      `json:"error"`
	ErrorDesc string      `json:"error_description,omitempty"`
	Data      interface{} `json:"data,omitempty"`
}

type RespAppTokenAuth struct {
	UserId int `json:"user_id"`
}

type ReqTest struct {
	Sdk     string `json:"sdk"`
	Version string `json:"version"`
}

func TestRequester(t *testing.T) {
	req := &ReqTest{
		Sdk:     "ios",
		Version: "0.0.1",
	}
	body, _ := json.Marshal(req)

	var r Response
	r.Data = &RespAppTokenAuth{}
	resp, body, err := requester.ReceiveContext(context.Background(), &r,
		requester.AddHeader("Authorization", "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzYWx0IjoiSkg0T05STkhHS3JUUWpCbmlLUGxvWDVhYlBlYndzQUJDWjEyIiwidXNlcl9pZCI6IjEifQ.dnj27yeLXgn3b5bab9n3lgL6GMp_qOaxAzsSXQnWTPPiGnuSXCmFJ-Am4Z2jrWj1DEBGNrEYy7p0iDTTV2uxVWzafoxEKysLGPmI3S0YvgGES0awSgrGDhy-XCKf6EBE_T3kfvvln-wThCGvSnNKwhHfRY41Lc0Nifg_syOdGcV8gQ547EyVjz04Ya5YhL54AD9Zrdti7AAWQnm1Z91arhJgGyFGpXx0BL3-DPfy2gtvxTgtmPwkFGdPqr_K5OgCRykgfet-hKJPrjuf0mZgJRZQuQy2J3g3bqFYUO7WkkQ7jN6MvmkQyShWDEWtEHEf-l5pJWsfAxfulkGp0gBjmQ"),
		//requester.AddHeader("Accept", "application/json"),
		//requester.AddHeader("Content-Type", "application/x-www-form-urlencoded"),
		//requester.Form(),
		requester.Post("http://127.0.0.1:5000/app/customer/get_user_data"),
		requester.Body(body),
		requester.ExpectCode(200),
	)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n %s\n %#v\n", resp, string(body), r.Data)
}
