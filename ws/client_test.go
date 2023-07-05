package ws

import (
	"testing"

	"github.com/starudream/go-lib/config"
)

func TestLogin(t *testing.T) {
	t.Log(Login(LoginParams{
		Room:     config.GetInt("douyu.room"),
		Stk:      config.GetString("douyu.stk"),
		Ltkid:    config.GetString("douyu.ltkid"),
		Username: config.GetString("douyu.username"),
	}))
}
