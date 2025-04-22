const fs = require("fs/promises");
const { readFileSync } = require("fs");
const path = require("path");

async function main() {
    let task;
    try {
        const taskJson = process.env.TASK_JSON;
        if (!taskJson) throw new Error("Не определена переменная окружения TASK_JSON");
        task = typeof taskJson === "object" ? taskJson : JSON.parse(taskJson);
    } catch (e) {
        process.stderr.write("[runner/cpp] Ошибка разбора TASK_JSON: " + e.message + "\n");
        process.exit(2);
    }

    const { function: functionName, args = [], wasmPath } = task;
    if (!wasmPath) {
        process.stderr.write("[runner/cpp] Не указан путь wasmPath в TASK_JSON\n");
        process.exit(3);
    }

    if (!functionName) {
        process.stderr.write("[runner/cpp] Не указана функция function в TASK_JSON\n");
        process.exit(4);
    }

    let module, instance;
    try {
        const wasmAbsPath = path.resolve(wasmPath);
        const wasmBytes = readFileSync(wasmAbsPath);
        module = await WebAssembly.compile(wasmBytes);
        instance = await WebAssembly.instantiate(module, {});
    } catch (e) {
        process.stderr.write(`[runner/cpp] Ошибка загрузки WASM: ${e.message}\n`);
        process.exit(5);
    }

    if (
        !instance.exports ||
        typeof instance.exports[functionName] !== "function"
    ) {
        process.stderr.write(
            `[runner/cpp] Функция '${functionName}' не экспортируется wasm-модулем\n`
        );
        process.exit(6);
    }

    try {
        instance.exports[functionName](...args);
    } catch (e) {
        process.stderr.write(
            `[runner/cpp] Ошибка выполнения функции '${functionName}': ${e.message}\n`
        );
        process.exit(7);
    }
}

main();
