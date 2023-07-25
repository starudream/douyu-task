package main

import (
	"context"
	"os"

	"github.com/starudream/go-lib/app"
	"github.com/starudream/go-lib/config"
	"github.com/starudream/go-lib/constant"
	"github.com/starudream/go-lib/cronx"
	"github.com/starudream/go-lib/flag"
	"github.com/starudream/go-lib/log"
)

var rootCmd = &flag.Command{
	Use:     constant.NAME,
	Version: constant.VERSION + " (" + constant.BIDTIME + ")",
	Run: func(cmd *flag.Command, args []string) {
		app.Add(cronjob)
		err := app.Go()
		if err != nil {
			log.Fatal().Msgf("app init fail: %v", err)
		}
	},
}

func init() {
	config.SetDefault("cron.renewal", "0 0 12 * * 0")

	rootCmd.AddCommand(runCmd)
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func cronjob(context.Context) error {
	if cron := config.GetString("cron.refresh"); cron != "" {
		if _, err := cronx.AddJob(cron, NewRefresh()); err != nil {
			return err
		}
	}

	if cron := config.GetString("cron.renewal"); cron != "" {
		if _, err := cronx.AddJob(cron, NewRenewal()); err != nil {
			return err
		}
	}

	go cronx.Run()

	return nil
}
