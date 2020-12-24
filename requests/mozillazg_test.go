package request_test

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/mozillazg/request"
)

type TokenInfo struct {
	UserId               int    `json:"user_id"`
	AccessToken          string `json:"access_token"`
	ExpireSeconds        string `json:"expire_seconds"`
	RefreshToken         string `json:"refresh_token"`
	RefreshExpireSeconds string `json:"refresh_expire_seconds"`
}

func TestGet(t *testing.T) {
	req := &Request{
		Url: "http://127.0.0.1:9001/app/token/get",
		Params: map[string]string{
			"user_id": "1",
		},
		Response: Response{
			Data: &TokenInfo{},
		},
	}
	err := req.Get()
	if err != nil {
		log.Fatalf("shit %s\n", err.Error())
	}

	var data *TokenInfo
	data = req.Data.(*TokenInfo)
	fmt.Printf("%#v\n", data)
}

func TestPost(t *testing.T) {
	req := &Request{
		Url: "http://127.0.0.1:9001/app/token/create",
		Params: map[string]string{
			"user_id": "1",
		},
		Response: Response{
			Data: &TokenInfo{},
		},
	}
	err := req.Post()
	if err != nil {
		log.Fatalf("shit %s\n", err.Error())
	}

	var data *TokenInfo
	data = req.Data.(*TokenInfo)
	fmt.Printf("%#v\n", data)
}

type Request struct {
	Url    string
	Params map[string]string
	Response
}

func (rq *Request) Get() error {
	var resp *request.Response
	var err error

	c := new(http.Client)
	req := request.NewRequest(c)

	if rq.Params != nil {
		params := ""
		for k, v := range rq.Params {
			params += k + "=" + v
		}
		rq.Url += "?" + params
	}
	resp, err = req.Get(rq.Url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	content, err := resp.Content()
	if err != nil {
		return err
	}

	err = json.Unmarshal(content, &(rq.Response))
	if err != nil {
		return err
	}

	return nil
}

func (rq *Request) Post() error {
	var resp *request.Response
	var err error

	c := new(http.Client)
	req := request.NewRequest(c)
	req.Data = rq.Params
	resp, err = req.Post(rq.Url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	content, err := resp.Content()
	if err != nil {
		return err
	}

	err = json.Unmarshal(content, &(rq.Response))
	if err != nil {
		return err
	}

	return nil
}
