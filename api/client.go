package api

import (
	"fmt"
	"net/http"

	"github.com/starudream/go-lib/core/v2/utils/osutil"
	"github.com/starudream/go-lib/resty/v2"

	"github.com/starudream/douyu-task/config"
)

const URL = "https://www.douyu.com"

type Client struct {
	c *resty.Client

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

func NewC(douyu config.Douyu) (*Client, error) {
	return New(douyu.Did, douyu.Ltp0)
}

func MustNew(douyu config.Douyu) *Client {
	c, err := NewC(douyu)
	osutil.PanicErr(err)
	return c
}

func (c *Client) R() *resty.Request {
	if c.c == nil {
		c.c = resty.New().SetHeader(resty.HeaderUserAgent, resty.UAWindowsChrome)
	}
	return c.c.R()
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
