API_ENDPOINT=http://localhost:8080/api/design
EXPERIMENT_ENDPOINT=http://localhost:8080/api
DESIGN_FILE=./configs/design.json
JS_FILE=./out/rs_exps.js
WASM_FILE=./out/rs_exps_bg.wasm

all: pkg

build:
	wasm-pack build --target nodejs

pkg: build
	mkdir -p out
	cp pkg/rs_exps.js out/rs_exps.js
	cp pkg/rs_exps_bg.wasm out/rs_exps_bg.wasm

create-design: pkg
	@echo "Создание нового эксперимента через API..."
	python3 scripts/design-create.py \
		--endpoint $(API_ENDPOINT) \
		--design $(DESIGN_FILE) \
		--js $(JS_FILE) \
		--wasm $(WASM_FILE)

experiment-start:
	python3 scripts/experiment-start.py \
		--endpoint $(EXPERIMENT_ENDPOINT) \
		--name rs_exps \
		--repeats 1000 \
		--warmup

clean:
	@echo "Очистка сгенерированных файлов... (clean)"
	rm -rf out

.PHONY: all build pkg clean create-design experiment-start
