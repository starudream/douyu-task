package api

import (
	"fmt"
	"net/http"

	"github.com/starudream/go-lib/config"
	"github.com/starudream/go-lib/httpx"

	"github.com/starudream/douyu-task/consts"
)

const URL = "https://www.douyu.com"

func init() {
	config.SetDefault("douyu.ua", consts.UserAgent)

	httpx.SetUserAgent(config.GetString("douyu.ua"))
}

type Client struct {
	Did  string // cookie: dy_did
	Ltp0 string // cookie: ltp0

	Uid      string // cookie: acf_uid
	Auth     string // cookie: acf_auth
	Stk      string // cookie: acf_stk
	Ltkid    string // cookie: acf_ltkid
	Username string // cookie: acf_username
}

func New(did, ltp0 string) (*Client, error) {
	c := &Client{
		Did:  did,
		Ltp0: ltp0,
	}
	if c.Did == "" || c.Ltp0 == "" {
		return nil, fmt.Errorf("missing config")
	}
	return c, c.Refresh()
}

func NewFromEnv() (*Client, error) {
	return New(config.GetString("douyu.did"), config.GetString("douyu.ltp0"))
}

func (c *Client) genAuthCookies() []*http.Cookie {
	return []*http.Cookie{
		{Name: "dy_did", Value: c.Did},
		{Name: "acf_uid", Value: c.Uid},
		{Name: "acf_auth", Value: c.Auth},
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
