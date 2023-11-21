package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/starudream/go-lib/cobra/v2"
	"github.com/starudream/go-lib/core/v2/slog"
	"github.com/starudream/go-lib/cron/v2"
	"github.com/starudream/go-lib/ntfy/v2"

	"github.com/starudream/douyu-task/config"
	"github.com/starudream/douyu-task/job"
)

var cronCmd = cobra.NewCommand(func(c *cobra.Command) {
	c.Use = "cron"
	c.Short = "Run as cron job"
	c.RunE = func(cmd *cobra.Command, args []string) error {
		return cronRun()
	}
})

func init() {
	rootCmd.AddCommand(cronCmd)
}

func cronRun() error {
	err1 := cron.AddJob(config.C().Cron.Refresh, "douyu-cron-refresh", cronRefresh)
	if err1 != nil {
		return fmt.Errorf("add cron job error: %w", err1)
	}
	err2 := cron.AddJob(config.C().Cron.Renewal, "douyu-cron-renewal", cronRenewal)
	if err2 != nil {
		return fmt.Errorf("add cron job error: %w", err2)
	}
	cron.Run()
	return nil
}

func cronRefresh() {
	msg := "斗鱼每日刷新礼物"
	err := job.Refresh()
	if err != nil {
		msg += fmt.Sprintf("失败：%v", err)
		slog.Error(msg)
	} else {
		msg += "成功"
		slog.Info(msg)
	}
	err = ntfy.Notify(context.Background(), msg)
	if err != nil && !errors.Is(err, ntfy.ErrNoConfig) {
		slog.Error("cron notify error: %v", err)
	}
	return
}

func cronRenewal() {
	msg := "斗鱼送礼物续粉丝牌"
	err := job.Renewal()
	if err != nil {
		if errors.Is(err, job.ErrTodayNotExpired) {
			return
		}
		msg += fmt.Sprintf("失败：%v", err)
		slog.Error(msg)
	} else {
		msg += "成功"
		slog.Info(msg)
	}
	err = ntfy.Notify(context.Background(), msg)
	if err != nil && !errors.Is(err, ntfy.ErrNoConfig) {
		slog.Error("cron notify error: %v", err)
	}
	return
}
