API_ENDPOINT=http://localhost:8080/api/design
DESIGN_FILE=./configs/design.json
WASM_FILE=./out/integrate.wasm

.PHONY: build create-design experiment-start

export GOARCH=wasm
export GOOS=js

build:
	@echo "Выполняется сборка WebAssembly файла... (build)"
	mkdir -p out
	go build -o $(PWD)/out/integrate.wasm $(PWD)/lib/lib.go

create-design: build
	@echo "Создание нового Go-эксперимента через API..."
	python3 scripts/design-create.py \
		--endpoint $(API_ENDPOINT) \
		--design $(DESIGN_FILE) \
		--wasm $(WASM_FILE)

experiment-start:
	@echo "Запуск Go-эксперимента через API..."
	python3 scripts/experiment-start.py \
		--endpoint http://localhost:8080/api \
		--name go_exps \
		--repeats 1000 \
		--warmup
