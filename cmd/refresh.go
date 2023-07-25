package main

import (
	"fmt"

	"github.com/starudream/go-lib/bot"
	"github.com/starudream/go-lib/config"
	"github.com/starudream/go-lib/cronx"
	"github.com/starudream/go-lib/log"

	"github.com/starudream/douyu-task/api"
	"github.com/starudream/douyu-task/ws"
)

// NewRefresh 刷新背包荧光棒
func NewRefresh() *Refresh {
	r := &Refresh{
		room: config.GetInt("douyu.room"),
	}
	return r
}

type Refresh struct {
	room int
}

var _ cronx.Job = (*Refresh)(nil)

func (r *Refresh) Name() string {
	return "refresh"
}

func (r *Refresh) Run() {
	err := r.do()
	if err != nil {
		log.Error().Msgf("refresh error: %v", err)
		_ = bot.Send("刷新礼物失败：" + err.Error())
	} else {
		_ = bot.Send("刷新礼物成功")
	}
}

func (r *Refresh) do() error {
	client, err := api.NewFromEnv()
	if err != nil {
		return fmt.Errorf("init api error: %w", err)
	}

	gifts1, err := client.ListGifts()
	if err != nil {
		return fmt.Errorf("list gifts error: %w", err)
	}
	log.Info().Msgf("gifts before:\n%s", gifts1.TableString())

	err = ws.Login(ws.LoginParams{
		Room:     r.room,
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
	log.Info().Msgf("gifts after:\n%s", gifts2.TableString())

	return nil
}
