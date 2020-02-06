# Copyright 2019 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


# git options
GIT_COMMIT ?= $(shell git rev-parse HEAD)
GIT_TAG    ?= $(shell git tag --points-at HEAD)
DIST_TYPE  ?= snapshot
BRANCH     ?= $(shell git rev-parse --abbrev-ref HEAD)

PROJECT_NAME := gospal
PKG_ORG      := github.com/contentsquare
PKG 		 := $(PKG_ORG)/$(PROJECT_NAME)
PKG_LIST 	 := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES 	 := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)

GO			 := go
GOFMT		 := gofmt
GOOS		 ?= linux
GOARCH		 ?= amd64

tidy:
	$(info =====  $@  =====)
	GO111MODULE=on go mod tidy

deps:
	$(info =====  $@  =====)
	GO111MODULE=on go mod vendor

format:
	$(info =====  $@  =====)
	$(GOFMT) -w -s $(GO_FILES)

test:
	$(info =====  $@  =====)
	$(GO) test -v -race -cover -coverprofile=coverage.out  $(PKG_LIST)

fmt:
	$(info =====  gofmt =====)
	$(GOFMT) -d -e -s $(GO_FILES)

lint:
	$(info =====  $@  =====)
	$(GO) vet $(PKG_LIST)
	$(GO) list ./... | grep -Ev /vendor/ | xargs -L1 golint -set_exit_status

version:
	$(info =====  $@  =====)
ifneq ($(GIT_TAG),)
	$(eval VERSION := $(GIT_TAG))
else
	$(eval VERSION := $(subst /,-,$(BRANCH)))
	$(eval VERSION_FILE := $(GIT_COMMIT)-SNAPSHOT)
endif
	@test -n "$(VERSION)"
	$(info Building $(VERSION)/$(VERSION_FILE) on sha1 $(GIT_COMMIT))

.PHONY: tidy \
		deps \
		format \
		test \
		fmt \
		lint \
		version

