TASK_JSON_CPP := $(shell cat test/data/cpp/task.json | tr -d '\n')
HYPERFINE_RESULT_CPP := test/data/cpp/hyperfine.json
test-cpp:
	@echo '=== Benchmarking C++ wasm ==='
	TASK_JSON='$(TASK_JSON_CPP)' \
	hyperfine --warmup 50 --runs 100 'bun ./scripts/runner/cpp.js' \
	--show-output \
	--export-json '$(HYPERFINE_RESULT_CPP)'

TASK_JSON_GO := $(shell cat test/data/go/task.json | tr -d '\n')
HYPERFINE_RESULT_GO := test/data/go/hyperfine.json
test-go:
	@echo '=== Benchmarking Go wasm ==='
	TASK_JSON='$(TASK_JSON_GO)' \
	hyperfine --warmup 50 --runs 100 'bun ./scripts/runner/go.js' \
	--show-output \
	--export-json '$(HYPERFINE_RESULT_GO)'

.PHONY: test-cpp test-go
