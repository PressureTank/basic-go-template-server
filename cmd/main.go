// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package main

import (
	"context"

	"github.com/spacemonkeygo/monkit/v3"
	"github.com/spf13/cobra"
	"github.com/zeebo/errs"
	"go.uber.org/zap"

	"storj.io/private/cfgstruct"
	"storj.io/private/process"
	"storj.io/qa-storj/server"
	"storj.io/qa-storj/static"
)

var cfg struct {
	server.Config

	// other config options could go here
}

func main() {
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run the web server",
		RunE:  cmdRun,
	}
	rootCmd := &cobra.Command{
		Use:   "webserver",
		Short: "web server",
	}
	rootCmd.AddCommand(runCmd)
	process.SetHardcodedApplicationName("webserver")
	process.Bind(runCmd, &cfg, cfgstruct.DefaultsFlag(runCmd))
	process.Exec(rootCmd)
}

func cmdRun(cmd *cobra.Command, args []string) (err error) {
	ctx := context.Background()
	log := zap.L()
	err = process.InitMetricsWithHostname(ctx, log, monkit.Default)
	if err != nil {
		return errs.New("Error initializing metrics: %+v", err)
	}

	s, err := server.New(log, cfg.Config, static.FS)
	if err != nil {
		return errs.New("Error creating web server: %+v", err)
	}

	runError := s.Serve()
	closeError := s.Close()

	return errs.Combine(runError, closeError)
}
