package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func (c *Client) Refresh() error {
	resp, err := c.R().
		SetCookie(&http.Cookie{Name: "dy_did", Value: c.Did}).
		SetCookie(&http.Cookie{Name: "LTP0", Value: c.Ltp0}).
		SetHeader("referer", URL).
		SetQueryParam("client_id", "1").
		SetQueryParam("t", strconv.Itoa(int(time.Now().UnixMilli()))).
		SetQueryParam("_", strconv.Itoa(int(time.Now().UnixMilli()))).
		SetQueryParam("callback", "axiosJsonpCallback").
		Get("https://passport.douyu.com/lapi/passport/iframe/safeAuth")
	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("http status: %s", resp.Status())
	}

	for _, cookie := range resp.Cookies() {
		if cookie.Name == "acf_uid" {
			c.Uid = cookie.Value
		}
		if cookie.Name == "acf_auth" {
			c.Auth = cookie.Value
		}
		if cookie.Name == "acf_stk" {
			c.Stk = cookie.Value
		}
		if cookie.Name == "acf_ltkid" {
			c.Ltkid = cookie.Value
		}
		if cookie.Name == "acf_username" {
			c.Username = cookie.Value
		}
	}

	if c.Uid != "" && c.Auth != "" && c.Stk != "" && c.Ltkid != "" && c.Username != "" {
		return nil
	}

	return fmt.Errorf("cookies not found")
}
