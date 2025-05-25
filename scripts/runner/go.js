const fs = require("fs");
const path = require("path");
require("./go/wasm_exec.js");

function errorExit(msg, code) {
    process.stderr.write(msg + "\n");
    process.exit(code);
}

function main() {
    let taskJson = process.env.TASK_JSON;
    if (!taskJson)
        errorExit("[runner/go] Не определена переменная окружения TASK_JSON", 2);

    let task;
    try {
        task = typeof taskJson === "object" ? taskJson : JSON.parse(taskJson);
    } catch (e) {
        errorExit("[runner/go] Ошибка разбора TASK_JSON: " + e.message, 2);
    }

    const { function: functionName, args = [], wasmPath } = task;
    if (!wasmPath)
        errorExit("[runner/go] Не указан путь wasmPath в TASK_JSON", 3);
    if (!functionName)
        errorExit("[runner/go] Не указана функция function в TASK_JSON", 4);

    let wasmBytes;
    try {
        const wasmAbsPath = path.resolve(wasmPath);
        wasmBytes = fs.readFileSync(wasmAbsPath);
    } catch (e) {
        errorExit(`[runner/go] Ошибка загрузки wasm-файла: ${e.message}`, 5);
    }

    const go = new global.Go();

    WebAssembly.instantiate(wasmBytes, go.importObject)
        .then(result => {
            go.run(result.instance);

            const fn = globalThis[functionName];
            if (typeof fn !== "function")
                errorExit(`[runner/go] Функция '${functionName}' не экспортируется wasm-модулем`, 6);

            try {
                fn(...args);
            } catch (e) {
                errorExit(
                    `[runner/go] Ошибка выполнения функции '${functionName}': ${e.message}`,
                    7
                );
            }
        })
        .catch(e => {
            errorExit(`[runner/go] Ошибка инициализации WASM: ${e.message}`, 5);
        });
}

main();
