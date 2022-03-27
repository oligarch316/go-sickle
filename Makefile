#!/usr/bin/make -f

# Disable all default suffixes
.SUFFIXES:

cmd_dhall := $(shell command -v dhall || echo "dhall")
cmd_go := $(shell command -v go || echo "go")
cmd_golint := $(shell command -v golint || echo "golint")


# ----- Default
.PHONY: default

default: binary


# ----- Source
.PHONY: source.fmt source.lint

source.fmt: ## Format source
	$(info Formatting source)
	@$(cmd_go) fmt ./...

source.lint: ## Lint source
	$(info Linting source)
	@$(cmd_golint) ./...


# ----- Config
.PHONY: config.assets config.fmt

config_sources := $(shell find config -type f -name '*.dhall')

config.fmt: ## Format config
	$(info Formatting config)
	@$(cmd_dhall) format $(config_sources)


# ----- Assets
.PHONY: assets

asset_dir := pkg/embedded/assets

asset_target_config := $(asset_dir)/config.bytes
asset_targets := $(asset_target_config)

$(asset_target_config): $(config_sources)
	$(info Building '$@')
	@mkdir -p $(@D)
	@$(cmd_dhall) --plain --ascii \
		<<< 'let x = ./config/schemas.dhall in x.Config.default' \
		| $(cmd_dhall) encode > $@

assets: $(asset_targets)  ## Build embedded assets

assets.clean: ## Clean embedded assets
	$(info Cleaning embedded assets)
	@rm pkg/embedded/assets/* 2> /dev/null || true


# ----- Binary
.PHONY: binary binary.clean

binary_tail := $(shell go env GOOS)_$(shell go env GOARCH)
binary_flags :=
binary_ld_flags :=

binary_sources := cmd/main.go $(shell find pkg -type f -name '*.go')
binary_target := bin/sickle_$(binary_tail)

$(binary_target): $(binary_sources) $(asset_targets) go.mod go.sum
	$(info Building '$@')
	@mkdir -p $(@D)
	@$(cmd_go) build $(build_flags) $(build_ld_flags) -o $@ $<

binary: $(binary_target) ## Compile binary

binary.clean: ## Clean binary artifacts
	$(info Cleaning build artifacts)
	@rm bin/* 2> /dev/null || true


# ----- Test
# TODO


# ----- Tooling
.PHONY: fmt tidy

fmt: source.fmt config.fmt ## Format all

tidy: ## Tidy modules
	$(info Tidying)
	@$(cmd_go) mod tidy

# ----- Clean
.PHONY: clean

clean: assets.clean binary.clean ## Clean all artifacts


# ----- HELP
.PHONY: help

print-%: ; @echo "$($*)"

help: ## Show help information
	@awk -F ':|##' '/^[^\t].+?:.*?##/ {printf "\033[36m%-30s\033[0m %s\n", $$1, $$NF }' $(MAKEFILE_LIST);
