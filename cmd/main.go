package main

import (
	"context"

	"github.com/starudream/go-lib/app"
	"github.com/starudream/go-lib/config"
	"github.com/starudream/go-lib/cronx"
	"github.com/starudream/go-lib/log"
)

func init() {
	config.SetDefault("cron.refresh", "0 0 1 * * *")
	config.SetDefault("cron.send_gift", "0 0 12 * * 0")
}

func main() {
	app.Add(startup)
	app.Add(cron)
	err := app.Go()
	if err != nil {
		log.Fatal().Msgf("app init fail: %v", err)
	}
}

func startup(context.Context) error {
	if config.GetBool("startup") {
		NewRenewal(false).Run()
	}
	return nil
}

func cron(context.Context) error {
	if _, err := cronx.AddJob(config.GetString("cron.refresh"), NewRenewal(true)); err != nil {
		return err
	}

	if _, err := cronx.AddJob(config.GetString("cron.send_gift"), NewRenewal(false)); err != nil {
		return err
	}

	go cronx.Run()

	return nil
}
