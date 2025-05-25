TASK_JSON_CPP := $(shell cat test/data/cpp/task.json | tr -d '\n')
HYPERFINE_RESULT_CPP := test/data/cpp/hyperfine.json
test-cpp:
	@echo '=== Benchmarking C++ wasm ==='
	TASK_JSON='$(TASK_JSON_CPP)' \
	hyperfine --warmup 15 --runs 1000 'bun ./scripts/runner/cpp.js' \
	--show-output \
	--export-json '$(HYPERFINE_RESULT_CPP)'

TASK_JSON_GO := $(shell cat test/data/go/task.json | tr -d '\n')
HYPERFINE_RESULT_GO := test/data/go/hyperfine.json
test-go:
	@echo '=== Benchmarking Go wasm ==='
	TASK_JSON='$(TASK_JSON_GO)' \
	hyperfine --warmup 15 --runs 1000 'bun ./scripts/runner/go.js' \
	--show-output \
	--export-json '$(HYPERFINE_RESULT_GO)'


TASK_JSON_RUST := $(shell cat test/data/rs/task.json | tr -d '\n')
HYPERFINE_RESULT_RUST := test/data/rs/hyperfine.json
test-rust:
	@echo '=== Benchmarking Rust wasm ==='
	TASK_JSON='$(TASK_JSON_RUST)' \
	hyperfine --warmup 15 --runs 1000 'bun ./scripts/runner/rs.js' \
	--show-output \
	--export-json '$(HYPERFINE_RESULT_RUST)'


TASK_JSON_JS := $(shell cat test/data/js/task.json | tr -d '\n')
HYPERFINE_RESULT_JS := test/data/js/hyperfine.json
test-js:
	@echo '=== Benchmarking JS ==='
	TASK_JSON='$(TASK_JSON_JS)' \
	hyperfine --warmup 15 --runs 1000 'bun ./scripts/runner/js.js' \
	--show-output \
	--export-json '$(HYPERFINE_RESULT_JS)'

.PHONY: test-cpp test-go test-rust test-js
