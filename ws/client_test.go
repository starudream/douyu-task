package ws

import (
	"testing"

	"github.com/starudream/go-lib/config"

	"github.com/starudream/douyu-task/api"
)

func TestLogin(t *testing.T) {
	client, err := api.NewFromEnv()
	if err != nil {
		t.Fatal(err)
	}
	err = Login(LoginParams{
		Room:     config.GetInt("douyu.room"),
		Stk:      client.Stk,
		Ltkid:    client.Ltkid,
		Username: client.Username,
	})
	if err != nil {
		t.Fatal(err)
	}
}
