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

package controller

import (
	"context"

	"github.com/Comcast/DNS-over-HTTPs-translator"
	"go.uber.org/zap"
)

// ctrl is the type that implements the Controller interface
type ctrl struct {
	logger *zap.Logger
	proxy  translator.Proxy
}

/* Below lies the implementation to create and set a new instance of Controller */

// New returns an instance of Controller
func New(options ...Option) (result translator.Controller, err error) {
	result = &ctrl{}

	// apply options
	for _, opt := range options {
		opt(result.(*ctrl))
	}

	return
}

func (c *ctrl) Proxy() translator.Proxy {
	return c.proxy
}

func (c *ctrl) Start(ctx context.Context) (err error) {
	err = c.proxy.Start(ctx)
	if err != nil {
		c.logger.Error("Proxy execution error", zap.Error(err))
	}

	for {
		select {
		case <-ctx.Done():
			c.logger.Info("Controller Start - Context cancelled")
			return
		}
	}
}