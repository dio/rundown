# Copyright 2022 Dhi Aurrahman
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

.PHONY: license

# Include versions of tools we build or fetch on-demand.
include Tools.mk

name := rundown

# Root dir returns absolute path of current directory. It has a trailing "/".
root_dir := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

# Currently we resolve it using which. But more sophisticated approach is to use infer GOROOT.
go     := $(shell which go)
goarch := $(shell $(go) env GOARCH)
goexe  := $(shell $(go) env GOEXE)
goos   := $(shell $(go) env GOOS)

# Local cache directory.
CACHE_DIR ?= $(root_dir).cache

# Go tools directory holds the binaries of Go-based tools.
go_tools_dir := $(CACHE_DIR)/tools/go
# Prepackaged tools may have more than precompiled binaries, e.g. for protoc, it also has an include
# directory which contains well-known proto files: https://github.com/protocolbuffers/protobuf/tree/master/src/google/protobuf.
prepackaged_tools_dir := $(CACHE_DIR)/tools/prepackaged

# By default, a protoc-gen-<name> program is expected to be on your PATH so that it can be
# discovered and executed by buf. This makes sure the Go-based and prepackaged tools dirs are
# registered in the PATH for buf to pick up. As an alternative, we can specify "path"
# https://docs.buf.build/configuration/v1/buf-gen-yaml#path for each plugin entry in buf.gen.yaml,
# however that means we need to override buf.gen.yaml at runtime. Note: since remote plugin
# execution https://docs.buf.build/bsr/remote-generation/remote-plugin-execution is available, one
# should check that out first before downloading local protoc plugins.
export PATH := $(go_tools_dir):$(prepackaged_tools_dir)/bin:$(PATH)

# Pre-packaged targets.
clang-format := $(prepackaged_tools_dir)/bin/clang-format

# Go-based tools.
addlicense          := $(go_tools_dir)/addlicense
buf                 := $(go_tools_dir)/buf
goimports           := $(go_tools_dir)/goimports
golangci-lint       := $(go_tools_dir)/golangci-lint
protoc-gen-go       := $(go_tools_dir)/protoc-gen-go
protoc-gen-validate := $(go_tools_dir)/protoc-gen-validate

# Assorted tools required for processing proto files.
proto_tools := \
	$(buf) \
	$(protoc-gen-go) \
	$(protoc-gen-validate)

# We cache the deps fetched by buf locally (in-situ) by setting BUF_CACHE_DIR
# https://docs.buf.build/bsr/overview#module-cache, so it can be referenced by other tools.
export BUF_CACHE_DIR := $(root_dir).cache/buf
BUF_V1_MODULE_DATA   := $(BUF_CACHE_DIR)/v1/module/data/buf.build


# This is adopted from https://github.com/tetratelabs/func-e/blob/3df66c9593e827d67b330b7355d577f91cdcb722/Makefile#L60-L76.
# ANSI escape codes. f_ means foreground, b_ background.
# See https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_(Select_Graphic_Rendition)_parameters.
f_black            := $(shell printf "\33[30m")
b_black            := $(shell printf "\33[40m")
f_white            := $(shell printf "\33[97m")
f_gray             := $(shell printf "\33[37m")
f_dark_gray        := $(shell printf "\33[90m")
f_bright_cyan      := $(shell printf "\33[96m")
b_bright_cyan      := $(shell printf "\33[106m")
ansi_reset         := $(shell printf "\33[0m")
ansi_$(name)       := $(b_black)$(f_black)$(b_bright_cyan)$(name)$(ansi_reset)
ansi_format_dark   := $(f_gray)$(f_bright_cyan)%-10s$(ansi_reset) $(f_dark_gray)%s$(ansi_reset)\n
ansi_format_bright := $(f_white)$(f_bright_cyan)%-10s$(ansi_reset) $(f_black)$(b_bright_cyan)%s$(ansi_reset)\n

# This formats help statements in ANSI colors. To hide a target from help, don't comment it with a trailing '##'.
help: ## Describe how to use each target
	@printf "$(ansi_$(name))$(f_white)\n"
	@awk 'BEGIN {FS = ":.*?## "} /^[0-9a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "$(ansi_format_dark)", $$1, $$2}' $(MAKEFILE_LIST)

# By default, unless GOMAXPROCS is set via an environment variable or explicity in the code, the
# tests are run with GOMAXPROCS=1. This is problematic if the tests require more than one CPU, for
# example when running t.Parallel() in tests.
export GOMAXPROCS ?=4
test: ## Run all unit tests
	@$(go) test ./internal/...

# $(call buf-generate,proto)

# "externals/authservice" has a special rule for now, since we do not want to touch its *.proto.
gen: $(BUF_V1_MODULE_DATA) ## To generate generated files from *.proto

define buf-generate
	$(buf) generate $1 --template generators/$1.gen.yaml
endef

update: ## Update authservice to latest commit
	@git submodule update --remote --merge

check: ## Make sure we follow the rules
	@rm -fr generated
	@$(MAKE) gen format lint license
	@$(go) mod tidy
	@if [ ! -z "`git status -s`" ]; then \
		echo "The following differences will fail CI until committed:"; \
		git diff --exit-code; \
	fi

license_ignore :=
license_files  := api examples internal proto buf.*.yaml Makefile *.mk
license: $(addlicense) ## To add license
	@$(addlicense) $(license_ignore) -c "Dhi Aurrahman"  $(license_files) 1>/dev/null 2>&1

all_nongen_go_sources := $(wildcard api/*.go example/*.go internal/*.go internal/*/*.go internal/*/*/*.go)
format: go.mod $(all_nongen_go_sources) $(goimports)
	@$(go) mod tidy
	@$(go)fmt -s -w $(all_nongen_go_sources)
# Workaround inconsistent goimports grouping with awk until golang/go#20818 or incu6us/goimports-reviser#50
	@for f in $(all_nongen_go_sources); do \
			awk '/^import \($$/,/^\)$$/{if($$0=="")next}{print}' $$f > /tmp/fmt; \
	    mv /tmp/fmt $$f; \
	done
	@$(goimports) -local $$(sed -ne 's/^module //gp' go.mod) -w $(all_nongen_go_sources)

# Override lint cache directory. https://golangci-lint.run/usage/configuration/#cache.
export GOLANGCI_LINT_CACHE=$(CACHE_DIR)/golangci-lint
lint: .golangci.yml $(all_nongen_go_sources) $(golangci-lint) ## Lint all Go sources
	@$(golangci-lint) run --timeout 5m --config $< ./...

authservice_dir := $(root_dir)externals/authservice
# BUF_V1_MODULE_DATA can only be generated by buf generate or build.
# Note that since we use newer buf binary, the buf.lock contains "version: v1" entry which is not
# backward compatible with older version of buf.
$(BUF_V1_MODULE_DATA): $(authservice_dir)/buf.yaml $(authservice_dir)/buf.lock $(proto_tools)
	@$(buf) lint
	@$(buf) build

$(authservice_dir)/buf.yaml:
	@git submodule update --init

# Catch all rules for Go-based tools.
$(go_tools_dir)/%:
	@GOBIN=$(go_tools_dir) go install $($(notdir $@)@v)

# Install clang-format from https://github.com/angular/clang-format. We don't support win32 yet as
# this script will fail.
clang-format-download-archive-url = https://$(subst @,/archive/refs/tags/,$($(notdir $1)@v)).tar.gz
clang-format-dir                  = $(subst github.com/angular/clang-format@v,clang-format-,$($(notdir $1)@v))
$(clang-format):
	@printf "$(ansi_format_dark)" tools "installing $($(notdir $@)@v)..."
	@mkdir -p $(dir $@)
	@curl -sSL $(call clang-format-download-archive-url,$@) | tar xzf - -C $(prepackaged_tools_dir)/bin \
		--strip 3 $(call clang-format-dir,$@)/bin/$(goos)_x64
	@printf "$(ansi_format_bright)" tools "ok"
