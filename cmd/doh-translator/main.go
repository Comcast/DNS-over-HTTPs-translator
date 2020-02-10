/**
* Copyright 2020 Comcast Cable Communications Management, LLC
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
* http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*
* SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/urfave/cli"
	"github.com/Comcast/DNS-over-HTTPs-translator"
	"github.com/Comcast/DNS-over-HTTPs-translator/pkg/config"
	"github.com/Comcast/DNS-over-HTTPs-translator/pkg/proxy"
	"github.com/Comcast/DNS-over-HTTPs-translator/pkg/controller"
	"go.uber.org/zap"
)

var (
	logger *zap.Logger
	cfg    translator.Config
	ctrl   translator.Controller
	pxy    translator.Proxy
)

func main() {
	// Initialize the cli package.
	app := cli.NewApp()

	app.Name = "DNS-over-HTTPs-translator"
	app.Version = "1.0.0"
	app.Compiled = time.Now()

	app.Before = func(c *cli.Context) (err error) {

		// Initialize Config
		if cfg, err = config.New(); err != nil {
			return err
		}

		// Initialize the logger
		if logger, err = cfg.Logger(); err != nil {
			return err
		}

		return
	}

	app.After = func(c *cli.Context) (err error) {
		// Flush logger buffers
		return logger.Sync()
	}

	app.Commands = []cli.Command{
		{
			Name:   "start",
			Usage:  "start the DOH translator",
			Action: TranslatorStart,
		},
	}

	// Run the cli app
	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
	}
}

// TranslatorStart starts the translator
func TranslatorStart(c *cli.Context) (err error) {
	logger.Info("Starting translator...")

	// Initialize Proxy
	if pxy, err = proxy.New(
		proxy.WithLogger(logger),
		proxy.WithResolver(cfg.Resolver()),
		proxy.WithListen(cfg.Listen()),
	); err != nil {
		logger.Error("Unable to create the Proxy instance", zap.Error(err))
		return err
	}

	// Initialize Controller
	if ctrl, err = controller.New(
		controller.WithLogger(logger),
		controller.WithProxy(pxy),
	); err != nil {
		logger.Error("Unable to create the Controller instance", zap.Error(err))
		return err
	}

	// Prepare main parent context with cancel function
	mainCtx, cf := context.WithCancel(context.Background())

	// Run controller
	go ctrl.Start(mainCtx)

	// Listen for signals to stop the translator
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// Wait for cancellation signal
	<-signals

	logger.Info("Received signal to stop translator, initiating graceful shutdown.")

	// Cancel main parent context
	cf()

	logger.Info("All contexts cancelled. Hit Ctrl + C to close program.")
	logger.Info("translator: Goodbye!")

	return
}
