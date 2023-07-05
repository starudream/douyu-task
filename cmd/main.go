package main

import (
	"context"

	"github.com/starudream/go-lib/app"
	"github.com/starudream/go-lib/config"
	"github.com/starudream/go-lib/cronx"
	"github.com/starudream/go-lib/log"
)

func init() {
	config.SetDefault("cron.send_gift", "0 0 12 * * 0")
}

func main() {
	app.Add(startup)
	app.Add(cronjob)
	err := app.Go()
	if err != nil {
		log.Fatal().Msgf("app init fail: %v", err)
	}
}

func startup(context.Context) error {
	if config.GetBool("startup") {
		NewRenewal().Run()
	}
	return nil
}

func cronjob(context.Context) error {
	if cron := config.GetString("cron.refresh"); cron != "" {
		if _, err := cronx.AddJob(cron, NewRefresh()); err != nil {
			return err
		}
	}

	if cron := config.GetString("cron.send_gift"); cron != "" {
		if _, err := cronx.AddJob(cron, NewRenewal()); err != nil {
			return err
		}
	}

	go cronx.Run()

	return nil
}
