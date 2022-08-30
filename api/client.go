package api

import (
	"net/http"
	"strings"
)

const URL = "https://www.douyu.com"

type Client struct {
	did  string // cookie: dy_did
	uid  string // cookie: acf_uid
	auth string // cookie: acf_auth
	ltp0 string // cookie: ltp0
}

func New(did, uid, auth, ltp0 string) *Client {
	return &Client{
		did:  did,
		uid:  uid,
		auth: auth,
		ltp0: ltp0,
	}
}

func NewWithCookie(cookie string) *Client {
	client := &Client{}
	ss := strings.Split(cookie, ";")
	for _, s := range ss {
		kvs := strings.Split(s, "=")
		if len(kvs) != 2 {
			continue
		}
		k, v := strings.TrimSpace(kvs[0]), strings.TrimSpace(kvs[1])
		switch k {
		case "dy_did":
			client.did = v
		case "acf_auth":
			client.auth = v
		case "acf_uid":
			client.uid = v
		}
	}
	if client.did == "" || client.uid == "" || client.auth == "" {
		return nil
	}
	return client
}

func (c *Client) genAuthCookies() []*http.Cookie {
	return []*http.Cookie{
		{Name: "dy_did", Value: c.did},
		{Name: "acf_uid", Value: c.uid},
		{Name: "acf_auth", Value: c.auth},
	}
}

type CommonResp struct {
	Error int    `json:"error"`
	Msg   string `json:"msg"`
}

func (resp CommonResp) GetError() int {
	return resp.Error
}

func (resp CommonResp) GetMsg() string {
	return resp.Msg
}
