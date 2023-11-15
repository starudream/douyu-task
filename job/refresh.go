package job

import (
	"fmt"

	"github.com/starudream/go-lib/core/v2/slog"

	"github.com/starudream/douyu-task/api"
	"github.com/starudream/douyu-task/config"
	"github.com/starudream/douyu-task/ws"
)

func Refresh() error {
	client, err := api.NewC(config.C().Douyu)
	if err != nil {
		return fmt.Errorf("new client error: %w", err)
	}

	gifts1, err := client.ListGifts()
	if err != nil {
		return fmt.Errorf("list gifts error: %w", err)
	}
	slog.Info("gifts before:\n%s", gifts1.TableString())

	err = ws.Login(ws.LoginParams{
		Room:     config.C().Douyu.Room,
		Stk:      client.Stk,
		Ltkid:    client.Ltkid,
		Username: client.Username,
	})
	if err != nil {
		return fmt.Errorf("login error: %w", err)
	}

	gifts2, err := client.ListGifts()
	if err != nil {
		return fmt.Errorf("list gifts error: %w", err)
	}
	slog.Info("gifts after:\n%s", gifts2.TableString())

	return nil
}
