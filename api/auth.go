package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/starudream/go-lib/httpx"
)

func (c *Client) Refresh() error {
	resp, err := httpx.R().
		SetCookie(&http.Cookie{Name: "dy_did", Value: c.did}).
		SetCookie(&http.Cookie{Name: "LTP0", Value: c.ltp0}).
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
		if cookie.Name == "acf_auth" {
			c.auth = cookie.Value
			return nil
		}
	}

	return fmt.Errorf("cookies not found")
}
