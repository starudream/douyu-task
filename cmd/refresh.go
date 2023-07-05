package main

import (
	"github.com/starudream/go-lib/bot"
	"github.com/starudream/go-lib/config"
	"github.com/starudream/go-lib/cronx"
	"github.com/starudream/go-lib/log"

	"github.com/starudream/douyu-task/ws"
)

// NewRefresh 刷新背包荧光棒
func NewRefresh() *Refresh {
	r := &Refresh{
		room:     config.GetInt("douyu.room"),
		stk:      config.GetString("douyu.stk"),
		ltkid:    config.GetString("douyu.ltkid"),
		username: config.GetString("douyu.username"),
	}
	if r.stk == "" || r.ltkid == "" || r.username == "" {
		log.Fatal().Msgf("douyu refresh missing config")
	}
	return r
}

type Refresh struct {
	room     int
	stk      string
	ltkid    string
	username string
}

var _ cronx.Job = (*Refresh)(nil)

func (r *Refresh) Name() string {
	return "refresh"
}

func (r *Refresh) Run() {
	err := ws.Login(ws.LoginParams{
		Room:     r.room,
		Stk:      r.stk,
		Ltkid:    r.ltkid,
		Username: r.username,
	})
	if err != nil {
		log.Error().Msgf("refresh: %v", err)
		_ = bot.Send("刷新背包失败：" + err.Error())
	} else {
		_ = bot.Send("刷新背包成功")
	}
}
