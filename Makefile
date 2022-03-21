#!/usr/bin/make -f

# Disable all default suffixes
.SUFFIXES:

cmd_go := $(shell command -v go || echo "go")


# ----- Aliases
.PHONY: default

default: build


# ----- Build
build_sources := $(shell find . -type f -name '*.go')
build_tail := _$(shell go env GOOS)_$(shell go env GOARCH)
build_flags :=
build_ld_flags :=

bin/sickle: $(build_sources) go.mod go.sum
	$(info Building '$@')
	@mkdir -p $(@D)
	@$(cmd_go) build $(build_flags) $(build_ld_flags) -o $@ cmd/main.go

build: bin/sickle ## Compile binary

build.clean: ## Clean build artifacts
	$(info Cleaning build artifacts)
	@rm -r bin/sickle 2> /dev/null || true


# ----- Clean
.PHONY: clean

clean: build.clean


# ----- HELP
.PHONY: help

print-%: ; @echo "$($*)"

help: ## Show help information
	@awk -F ':|##' '/^[^\t].+?:.*?##/ {printf "\033[36m%-30s\033[0m %s\n", $$1, $$NF }' $(MAKEFILE_LIST);
