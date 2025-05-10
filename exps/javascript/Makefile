API_ENDPOINT=http://localhost:8080/api/design
DESIGN_FILE=./out/design.json
JS_FILE=./out/lib_js.out.js

pkg:
	@echo "Выполняется сборка пакета... (pkg)"
	mkdir -p out
	cp configs/design.json out/design.json
	cp src/lib.js out/lib_js.out.js

create-design: pkg
	@echo "Создание нового эксперимента через API..."
	python3 scripts/design-create.py \
		--endpoint $(API_ENDPOINT) \
		--design $(DESIGN_FILE) \
		--js $(JS_FILE)

experiment-start:
	python3 scripts/experiment-start.py \
		--endpoint http://localhost:8080/api \
		--name js_exps \
		--repeats 1000 \
		--warmup
clean:
	@echo "Очистка сгенерированных файлов... (clean)"
	rm -rf out pkg.zip

.PHONY: all build pkg clean create-design experiment-start
