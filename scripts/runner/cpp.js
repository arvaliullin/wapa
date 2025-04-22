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
