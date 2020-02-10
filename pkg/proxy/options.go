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

package proxy

import (
	"go.uber.org/zap"
)

// Option defines the type for functional options of the Proxy instance
type Option func(*proxy)

// WithLogger adds an instance of the logger to the Proxy instance
func WithLogger(logger *zap.Logger) Option {
	return func(p *proxy) {
		p.logger = logger
	}
}

// WithResolver adds an instance of Config to the Proxy instance
func WithResolver(resolver string) Option {
	return func(p *proxy) {
		p.resolver = resolver
	}
}
