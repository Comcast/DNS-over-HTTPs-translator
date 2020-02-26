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

package translator

import (
	"context"
	"go.uber.org/zap"
)

type Config interface {
	Resolver() string
	Listen() string
	Logger() (logger *zap.Logger, err error)
}

type Controller interface {
	Start(ctx context.Context) (err error)
	Proxy() Proxy
}

type Proxy interface {
	Start(ctx context.Context) (err error)
	Stop() (err error)
}