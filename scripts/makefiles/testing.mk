TASK_JSON_CPP := $(shell cat test/data/cpp/task.json | tr -d '\n')
HYPERFINE_RESULT_CPP := test/data/cpp/hyperfine.json
test-cpp:
	@echo '=== Benchmarking C++ wasm ==='
	TASK_JSON='$(TASK_JSON_CPP)' \
	hyperfine --warmup 50 --runs 100 'bun ./scripts/runner/cpp.js' \
	--show-output \
	--export-json '$(HYPERFINE_RESULT_CPP)'

.PHONY: test-cpp
