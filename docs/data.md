Для с++ я уже реализовал
// scripts/runner/cpp.js
const fs = require("fs");
const path = require("path");

function errorExit(msg, code) {
    process.stderr.write(msg + "\n");
    process.exit(code);
}

function main() {
    let taskJson = process.env.TASK_JSON;
    if (!taskJson)
        errorExit("[runner/cpp] Не определена переменная окружения TASK_JSON", 2);

    let task;
    try {
        task = typeof taskJson === "object" ? taskJson : JSON.parse(taskJson);
    } catch (e) {
        errorExit("[runner/cpp] Ошибка разбора TASK_JSON: " + e.message, 2);
    }

    const { function: functionName, args = [], wasmPath } = task;
    if (!wasmPath)
        errorExit("[runner/cpp] Не указан путь wasmPath в TASK_JSON", 3);
    if (!functionName)
        errorExit("[runner/cpp] Не указана функция function в TASK_JSON", 4);

    let instance;
    try {
        const wasmAbsPath = path.resolve(wasmPath);
        const wasmBytes = fs.readFileSync(wasmAbsPath);
        instance = new WebAssembly.Instance(
            new WebAssembly.Module(wasmBytes),
            {}
        );
    } catch (e) {
        errorExit(`[runner/cpp] Ошибка загрузки WASM: ${e.message}`, 5);
    }

    const fn = instance.exports && instance.exports[functionName];
    if (typeof fn !== "function")
        errorExit(`[runner/cpp] Функция '${functionName}' не экспортируется wasm-модулем`, 6);

    try {
        fn(...args);
    } catch (e) {
        errorExit(
            `[runner/cpp] Ошибка выполнения функции '${functionName}': ${e.message}`,
            7
        );
    }
}

main();
// scripts/makefiles/testing.mk
TASK_JSON_CPP := $(shell cat test/data/cpp/task.json | tr -d '\n')
HYPERFINE_RESULT_CPP := test/data/cpp/hyperfine.json
test-cpp:
	@echo '=== Benchmarking C++ wasm ==='
	TASK_JSON='$(TASK_JSON_CPP)' \
	hyperfine --warmup 50 --runs 100 'bun ./scripts/runner/cpp.js' \
	--show-output \
	--export-json '$(HYPERFINE_RESULT_CPP)'

.PHONY: test-cpp

Моя очень старая реализация для rust в cmd/bench/main.js
// cmd/bench/main.js
const fs = require("fs/promises");
const path = require("path");

/**
 * BenchmarkRunner для выполнения функций из Rust WebAssembly.
 */
class BenchmarkRunner {
    constructor(wasmPackagePath, argsPath) {
        this.pkgPath = wasmPackagePath;
        this.argsPath = argsPath;
        this.module = null;
    }

    /**
     * Проверка существования файла.
     * @param {string} filePath
     */
    async checkFileExists(filePath) {
        try {
            await fs.access(filePath);
        } catch (err) {
            throw new Error(`Файл не найден: ${filePath}`);
        }
    }

    /**
     * Загрузка JSON-файла.
     * @param {string} filePath
     * @returns {Promise<object>}
     */
    async loadJSON(filePath) {
        await this.checkFileExists(filePath);
        try {
            const data = await fs.readFile(filePath, "utf8");
            return JSON.parse(data);
        } catch (err) {
            throw new Error(`Ошибка чтения или парсинга JSON файла ${filePath}: ${err.message}`);
        }
    }

    /**
     * Динамическая загрузка Rust WebAssembly модуля.
     */
    async loadWasmModule() {
        try {
            const resolvedPath = path.resolve(this.pkgPath);
            this.module = require(resolvedPath);
        } catch (err) {
            throw new Error(`Ошибка загрузки WASM Rust модуля из ${this.pkgPath}: ${err.message}`);
        }
    }

    async run() {
        try {
            await this.checkFileExists(this.pkgPath);
            const argsData = await this.loadJSON(this.argsPath);

            await this.loadWasmModule();

            const functions = argsData.functions || [];
            if (!Array.isArray(functions) || functions.length === 0) {
                throw new Error("Массив 'functions' в JSON пустой или некорректный.");
            }

            for (let i = 0; i < functions.length; i++) {
                const func = functions[i];
                const funcName = func.function;
                const funcArgs = func.args || [];

                if (!funcName) {
                    console.warn(`Функция с индексом ${i} не имеет имени 'function'. Пропуск...`);
                    continue;
                }

                if (typeof this.module[funcName] !== "function") {
                    console.error(`[Ошибка] Функция '${funcName}' не существует в модуле.`);
                    continue;
                }

                await this.runFunction(i, funcName, funcArgs);
            }
        } catch (err) {
            console.error(`[Ошибка] ${err.message}`);
        }
    }

    /**
     * Выполнение отдельной функции из WASM.
     * @param {number} index Индекс функции в массиве
     * @param {string} funcName Название функции
     * @param {Array} funcArgs Аргументы функции
     */
    async runFunction(index, funcName, funcArgs) {
        try {
            const start = performance.now();
            const resultValue = this.module[funcName](...funcArgs); // вызов функции с аргументами
            const end = performance.now();

            this.logResult(index, funcName, funcArgs, resultValue, end - start);
        } catch (err) {
            console.error(`[Ошибка] Не удалось выполнить функцию '${funcName}' с индексом ${index}: ${err.message}`);
        }
    }

    /**
     * Логирование результата выполнения функции.
     * @param {number} index Индекс функции в массиве
     * @param {string} funcName Название функции
     * @param {Array} funcArgs Аргументы функции
     * @param {*} resultValue Результат выполнения
     * @param {number} execTime Время выполнения (мс)
     */
    logResult(index, funcName, funcArgs, resultValue, execTime) {
        console.log(`\n[Функция #${index}]`);
        console.log(`Имя: ${funcName}`);
        console.log(`Аргументы: ${JSON.stringify(funcArgs)}`);
        console.log(`Результат: ${resultValue}`);
        console.log(`Время выполнения: ${execTime.toFixed(2)} мс`);
    }
}

async function main() {
    if (process.argv.length < 4) {
        console.error("Использование: bun cmd/bench/main.js ./pkg/wapa_integrate_rs.js./configs/wapa.json");
        process.exit(1);
    }

    const wasmPackagePath = path.resolve(process.cwd(), process.argv[2]);
    const argsFile = path.resolve(process.cwd(), process.argv[3]);
    const runner = new BenchmarkRunner(wasmPackagePath, argsFile);
    await runner.run();
}

main();

Реализуй scripts/runner/rust.js

В scripts/makefiles/testing.mk добавить test-rust
