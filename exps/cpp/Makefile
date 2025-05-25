API_ENDPOINT=http://localhost:8080/api/design
DESIGN_FILE=./out/design.json
JS_FILE=./out/lib_cpp.out.js
WASM_FILE=./out/lib_cpp.out.wasm

all: bench

build:
	mkdir -p ./out
	em++ ./src/lib.cpp \
		-s WASM=1 \
		-s EXPORTED_FUNCTIONS="['_x2Integrate','_x2IntegrateMock','_fibonacciRecursive','_fibonacciRecursiveMock','_fibonacciIterative','_fibonacciIterativeMock','_multiply','_multiplyMock','_multiplyVector','_multiplyVectorMock','_factorize','_factorizeMock']" \
		-s MODULARIZE -o ./out/lib_cpp.out.js

bench: build
	bun cmd/bench/main.js ./out/lib_cpp.out.wasm ./configs/design.json

pkg: build
	@echo "Выполняется сборка пакета... (pkg)"
	mkdir -p out
	cp configs/design.json out/design.json
	zip -r pkg.zip out/*

create-design: pkg
	@echo "Создание нового эксперимента через API..."
	python3 scripts/design-create.py \
		--endpoint $(API_ENDPOINT) \
		--design $(DESIGN_FILE) \
		--js $(JS_FILE) \
		--wasm $(WASM_FILE)

experiment-start:
	python3 scripts/experiment-start.py \
		--endpoint http://localhost:8080/api \
		--name cpp_exps \
		--repeats 1000 \
		--warmup
clean:
	@echo "Очистка сгенерированных файлов... (clean)"
	rm -rf out pkg.zip

.PHONY: all build pkg clean create-design experiment-start
