package main

import (
	"github.com/starudream/go-lib/cobra/v2"
)

var rootCmd = cobra.NewRootCommand(func(c *cobra.Command) {
	c.Use = "douyu-task"
	c.Run = cronCmd.Run

	cobra.AddConfigFlag(c)
})
