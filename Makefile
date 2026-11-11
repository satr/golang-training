.PHONY: help install test lint fmt fmt-check vet race coverage vuln layout check reset reset-all \
	01_core_01_values_and_types 01_core_02_control_flow \
	01_core_03_loops_and_range 01_core_04_functions_and_multiple_returns \
	01_core_05_slices_and_arrays 01_core_14_variables_constants_and_zero_values \
	03_se \
	03_se_01_config_parsing 03_se_02_dependency_injection \
	03_se_03_interfaces_and_test_doubles 03_se_04_http_handlers \
	03_se_05_http_clients_timeouts 03_se_06_middleware \
	03_se_07_structured_logging 03_se_08_metrics_boundaries \
	03_se_09_package_architecture 03_se_10_graceful_shutdown \
	03_se_11_error_handling_patterns 03_se_12_config_env_merge \
	03_se_13_json_serialization 03_se_14_generics_containers \
	03_se_15_testing_patterns \
	03_se_16_feature_flags

GO ?= go
GO_VERSION := 1.27.1
GOTOOLCHAIN ?= go$(GO_VERSION)+auto
TOOLS_DIR := $(CURDIR)/.bin
GOLANGCI_LINT := $(TOOLS_DIR)/golangci-lint
GOVULNCHECK := $(TOOLS_DIR)/govulncheck
GOLANGCI_LINT_VERSION := v2.13.2
GOVULNCHECK_VERSION := v1.8.0
export GOTOOLCHAIN

help:
	@printf '%s\n' \
	  'install                                      Install pinned development tools' \
	  'test                                         Run all unit tests' \
	  'lint                                         Run golangci-lint' \
	  'fmt                                          Format all Go files' \
	  'vet                                          Run go vet' \
	  'race                                         Run all tests with the race detector' \
	  'coverage                                     Write coverage.out' \
	  'vuln                                         Scan dependencies and reachable code' \
	  'layout                                       Verify curriculum layout' \
	  'check                                        Run layout, format, vet, lint, and tests' \
	  '01_core_01_values_and_types                  Test Values and basic types' \
	  '01_core_02_control_flow                      Test Control flow' \
	  '01_core_03_loops_and_range                   Test Loops and range' \
	  '01_core_04_functions_and_multiple_returns    Test Functions and multiple returns' \
	  '01_core_05_slices_and_arrays                 Test Slices and arrays' \
	   '01_core_14_variables_constants_and_zero_values Test Variables, constants, and zero values' \
	'03_se                                          Run all Software Engineering exercises' \
	'## Software Engineering:' \
	'03_se_01_config_parsing                        Test Config parsing' \
	'03_se_02_dependency_injection                  Test Dependency injection' \
	'03_se_03_interfaces_and_test_doubles            Test Interfaces and test doubles' \
	'03_se_04_http_handlers                         Test HTTP handlers' \
	'03_se_05_http_clients_timeouts                  Test HTTP client timeouts' \
	'03_se_06_middleware                            Test Middleware pattern with context' \
	'03_se_07_structured_logging                    Test Structured logging with context' \
	'03_se_08_metrics_boundaries                    Test Metrics boundaries and state machines' \
	'03_se_09_package_architecture                  Test Package architecture and visibility' \
	'03_se_10_graceful_shutdown                     Test Graceful shutdown with context and signals' \
	'03_se_11_error_handling_patterns               Test Idiomatic error handling patterns' \
	'03_se_12_config_env_merge                      Test Config loading with environment merge' \
	'03_se_13_json_serialization                    Test JSON serialization with custom methods' \
	'03_se_14_generics_containers                   Test Generics and type constraints' \
	'03_se_15_testing_patterns                      Test Testing patterns and table-driven tests' \
	   '03_se_16_feature_flags                       Test Feature flags and deterministic A/B testing' \
	'reset EXERCISE=<dir>                         Restore one solution from TODO.go' \
	  'reset-all                                    Restore every solution after confirmation'

install:
	@mkdir -p "$(TOOLS_DIR)"
	GOBIN="$(TOOLS_DIR)" "$(GO)" install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)
	GOBIN="$(TOOLS_DIR)" "$(GO)" install golang.org/x/vuln/cmd/govulncheck@$(GOVULNCHECK_VERSION)

test:
	"$(GO)" test ./...

lint:
	@test -x "$(GOLANGCI_LINT)" || { echo 'golangci-lint is missing; run make install' >&2; exit 1; }
	"$(GOLANGCI_LINT)" run

fmt:
	gofmt -w .

fmt-check:
	@test -z "$$(gofmt -l .)" || { echo 'Run make fmt; these files need formatting:' >&2; gofmt -l . >&2; exit 1; }

vet:
	"$(GO)" vet ./...

race:
	"$(GO)" test -race ./...

coverage:
	"$(GO)" test -coverprofile=coverage.out ./...

vuln:
	@test -x "$(GOVULNCHECK)" || { echo 'govulncheck is missing; run make install' >&2; exit 1; }
	"$(GOVULNCHECK)" ./...

layout:
	"$(GO)" run ./scripts/verify_layout.go

check: layout fmt-check vet lint test

01_core_01_values_and_types:
	"$(GO)" test ./exercises/01_core/01_values_and_types

01_core_02_control_flow:
	"$(GO)" test ./exercises/01_core/02_control_flow

01_core_03_loops_and_range:
	"$(GO)" test ./exercises/01_core/03_loops_and_range

01_core_04_functions_and_multiple_returns:
	"$(GO)" test ./exercises/01_core/04_functions_and_multiple_returns

01_core_05_slices_and_arrays:
	"$(GO)" test ./exercises/01_core/05_slices_and_arrays

03_se: 03_se_01_config_parsing 03_se_02_dependency_injection \
	03_se_03_interfaces_and_test_doubles 03_se_04_http_handlers \
	03_se_05_http_clients_timeouts 03_se_06_middleware \
	03_se_07_structured_logging 03_se_08_metrics_boundaries \
	03_se_09_package_architecture 03_se_10_graceful_shutdown \
	03_se_11_error_handling_patterns 03_se_12_config_env_merge \
	03_se_13_json_serialization 03_se_14_generics_containers \
	03_se_15_testing_patterns \
	03_se_16_feature_flags

03_se_01_config_parsing:
	"$(GO)" test ./exercises/03_software_engineering/01_config-parsing/...

03_se_02_dependency_injection:
	"$(GO)" test ./exercises/03_software_engineering/02_dependency-injection/...

03_se_03_interfaces_and_test_doubles:
	"$(GO)" test ./exercises/03_software_engineering/03_interfaces-and-test-doubles/...

03_se_04_http_handlers:
	"$(GO)" test ./exercises/03_software_engineering/04_http-handlers/...

03_se_05_http_clients_timeouts:
	"$(GO)" test ./exercises/03_software_engineering/05_http-clients-timeouts/...

03_se_06_middleware:
	"$(GO)" test ./exercises/03_software_engineering/06_middleware/...

03_se_07_structured_logging:
	"$(GO)" test ./exercises/03_software_engineering/07_structured-logging/...

03_se_08_metrics_boundaries:
	"$(GO)" test ./exercises/03_software_engineering/08_metrics-boundaries/...

03_se_09_package_architecture:
	"$(GO)" test ./exercises/03_software_engineering/09_package-architecture/...

03_se_10_graceful_shutdown:
	"$(GO)" test ./exercises/03_software_engineering/10_graceful-shutdown/...

03_se_11_error_handling_patterns:
	"$(GO)" test ./exercises/03_software_engineering/11_error-handling-patterns/...

03_se_12_config_env_merge:
	"$(GO)" test ./exercises/03_software_engineering/12_config-env-merge/...

03_se_13_json_serialization:
	"$(GO)" test ./exercises/03_software_engineering/13_json-serialization/...

03_se_14_generics_containers:
	"$(GO)" test ./exercises/03_software_engineering/14_generics-containers/...

03_se_15_testing_patterns:
	"$(GO)" test ./exercises/03_software_engineering/15_testing-patterns/...

03_se_16_feature_flags:
	"$(GO)" test ./exercises/03_software_engineering/16_feature-flags/...

reset:
	$(if $(EXERCISE),,$(error Set EXERCISE=<exercise-directory>))
	"$(GO)" run ./scripts/reset.go -exercise "$(EXERCISE)"

reset-all:
	"$(GO)" run ./scripts/reset.go -all
