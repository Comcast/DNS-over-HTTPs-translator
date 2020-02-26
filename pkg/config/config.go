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

package config

import (
	"fmt"
	"strings"

	translator "github.com/Comcast/DNS-over-HTTPs-translator"
	"github.com/spf13/viper"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type cfg struct {
	viperCfg *viper.Viper
}

// New returns an instance of Config
func New(options ...Option) (result translator.Config, err error) {
	vip := viper.New()

	// viper configuration
	vip.SetConfigName("config-doh-translator")
	vip.AddConfigPath("/etc/doh-translator/")

	if err := vip.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failure reading DOH translator config file: %s", err)
	}

	result = &cfg{
		viperCfg: vip,
	}

	// Apply options
	for _, opt := range options {
		opt(result.(*cfg))
	}

	return
}

// Resolver returns the DNS resolver that is used by the translator
func (c *cfg) Resolver() string {
	c.viperCfg.SetDefault("resolver", "75.75.75.75:53")
	return c.viperCfg.GetString("resolver")
}

// Listen returns the address/port that the translator/proxy will listen on
func (c *cfg) Listen() string {
	c.viperCfg.SetDefault("listen", ":80")
	return c.viperCfg.GetString("listen")
}

// Logger returns an instantiated logger for the translator
func (c *cfg) Logger() (logger *zap.Logger, err error) {
	logConfig := c.viperCfg.Sub("logger")
	if logConfig == nil {
		logConfig = viper.New()
	}

	logConfig.SetDefault("loglevel", "debug")
	logConfig.SetDefault("logfile", "")
	logConfig.SetDefault("stdout", false)

	// Initialize the logger
	al := zap.NewAtomicLevelAt(zapcore.DebugLevel)
	if err = al.UnmarshalText([]byte(logConfig.GetString("level"))); err != nil {
		return nil, err
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = al
	if logger, err = cfg.Build(); err != nil {
		return
	}

	cfg.OutputPaths = []string{}

	var logFile string
	if logFile = strings.TrimSpace(logConfig.GetString("logfile")); logFile != "" {
		cfg.OutputPaths = append(cfg.OutputPaths, logFile)
	}

	// Add stdout logging to output path.
	if logConfig.GetBool("stdout") {
		cfg.OutputPaths = append(cfg.OutputPaths, "stdout")
	}

	if logger, err = cfg.Build(); err != nil {
		return nil, err
	}

	return
}
