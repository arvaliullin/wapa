const fs = require("fs");
const path = require("path");

function errorExit(msg, code) {
    process.stderr.write(msg + "\n");
    process.exit(code);
}

async function main() {
    const taskJson = process.env.TASK_JSON;
    if (!taskJson)
        errorExit("[runner/cpp] Не определена переменная окружения TASK_JSON", 2);

    let task;
    try {
        task = typeof taskJson === "object" ? taskJson : JSON.parse(taskJson);
    } catch (e) {
        errorExit("[runner/cpp] Ошибка разбора TASK_JSON: " + e.message, 2);
    }

    const { function: functionName, args = [], jsPath } = task;
    if (!jsPath)
        errorExit("[runner/cpp] Не указан путь jsPath в TASK_JSON", 3);
    if (!functionName)
        errorExit("[runner/cpp] Не указана функция function в TASK_JSON", 4);

    let mod;
    try {
        const jsAbsPath = path.resolve(jsPath);
        if (!fs.existsSync(jsAbsPath)) {
            errorExit(`[runner/cpp] JS-мост не найден по пути: ${jsAbsPath}`, 5);
        }
        const createModule = require(jsAbsPath);

        mod = createModule();
        if (typeof mod.then === "function") {
            mod = await mod;
        }
    } catch (e) {
        errorExit(`[runner/cpp] Ошибка загрузки JS-моста: ${e.message}`, 5);
    }

    const fn = mod["_" + functionName] ?? mod[functionName];
    if (typeof fn !== "function")
        errorExit(`[runner/cpp] Функция '${functionName}' не экспортируется wasm-модулем (пробовал '${functionName}' и '_${functionName}')`, 6);
    try {
        fn(...args);
    } catch (e) {
        errorExit(`[runner/cpp] Ошибка выполнения функции '${functionName}': ${e.message}`, 7);
    }
}

main();
