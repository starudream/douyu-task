package ws

import (
	"testing"

	"github.com/starudream/go-lib/core/v2/utils/testutil"

	"github.com/starudream/douyu-task/api"
	"github.com/starudream/douyu-task/config"
)

func TestLogin(t *testing.T) {
	client := api.MustNew(config.C().Douyu)
	testutil.Nil(t, Login(LoginParams{
		Room:     config.C().Douyu.Room,
		Stk:      client.Stk,
		Ltkid:    client.Ltkid,
		Username: client.Username,
	}))
}
