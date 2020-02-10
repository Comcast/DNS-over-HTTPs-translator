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
	"github.com/Comcast/DNS-over-HTTPs-translator"
	"go.uber.org/zap"
)

// Option defines the type for functionals options of the Controller
type Option func(*ctrl)

// WithLogger adds an instance of the logger to the Controller instance
func WithLogger(logger *zap.Logger) Option {
	return func(c *ctrl) {
		c.logger = logger
	}
}

// WithProxy adds an instance of Proxy to the Controller instance
func WithProxy(proxy translator.Proxy) Option {
	return func(c *ctrl) {
		c.proxy = proxy
	}
}