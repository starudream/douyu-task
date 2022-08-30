package api

import (
	"testing"

	"github.com/starudream/go-lib/config"
	"github.com/starudream/go-lib/testx"
)

var client *Client

func init() {
	client = New(config.GetString("douyu.did"), config.GetString("douyu.uid"), config.GetString("douyu.auth"), config.GetString("douyu.ltp0"))
	if client.ltp0 != "" {
		err := client.Refresh()
		if err != nil {
			panic(err)
		}
	}
}

func TestNewWithCookie(t *testing.T) {
	cookie := "dy_did=654321; acf_auth=abc; acf_uid=123456"
	c := NewWithCookie(cookie)
	testx.RequireNotNilf(t, c, "NewWithCookie")
	testx.RequireEqualf(t, "654321", c.did, "did")
	testx.RequireEqualf(t, "abc", c.auth, "auth")
	testx.RequireEqualf(t, "123456", c.uid, "uid")
}
