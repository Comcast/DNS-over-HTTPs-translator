# Copyright 2020 Comcast Cable Communications Management, LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
# SPDX-License-Identifier: Apache-2.0

VERSION=$(shell cat VERSION)
BUILD_TIME=$(shell date +%FT%T%z)
REVISION=$(shell git rev-parse --short HEAD)
BUILD_DIR=.build
TIMESTAMP:=$(shell date +%s)
BINARY=doh-translator

# Use linker flags to provide version/build settings to the target
LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

.PHONY: clean
clean:
	@rm -rf ${BUILD_DIR}

.PHONY: test
test:	build
	@go test -race -v ./...

.PHONY: build
build: clean
	@echo Building for Darwin
	@env GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY}-darwin-amd64 cmd/doh-translator/main.go
	@echo Building for Linux
	@env GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY}-linux-amd64 cmd/doh-translator/main.go

.PHONY: dev
dev:
	go run -race ${LDFLAGS} cmd/doh-translator/main.go
