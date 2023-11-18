package job

import (
	"fmt"
	"time"

	"github.com/starudream/go-lib/core/v2/slog"
	"github.com/starudream/go-lib/core/v2/utils/timeutil"

	"github.com/starudream/douyu-task/api"
	"github.com/starudream/douyu-task/config"
	"github.com/starudream/douyu-task/ws"
)

const (
	RefreshLoop        = 5
	RefreshMinDuration = 3 * time.Second
	RefreshMaxDuration = 8 * time.Second
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

	loopCount := 0

login:

	err = ws.Login(ws.LoginParams{
		Room:     config.C().Douyu.Room,
		Stk:      client.Stk,
		Ltkid:    client.Ltkid,
		Username: client.Username,
	})
	if err != nil {
		if loopCount >= RefreshLoop {
			return fmt.Errorf("login error: %w", err)
		}
		slog.Error("login error when attempt %d: %v", loopCount, err)
		loopCount++
		duration := timeutil.JitterDuration(RefreshMinDuration, RefreshMaxDuration, loopCount)
		slog.Info("login will retry after %s", duration.String())
		time.Sleep(duration)
		goto login
	}

	gifts2, err := client.ListGifts()
	if err != nil {
		return fmt.Errorf("list gifts error: %w", err)
	}
	slog.Info("gifts after:\n%s", gifts2.TableString())

	return nil
}
