const fs = require("fs");
const path = require("path");

function errorExit(msg, code) {
    process.stderr.write(msg + "\n");
    process.exit(code);
}

function main() {
    const taskJson = process.env.TASK_JSON;
    if (!taskJson)
        errorExit("[runner/js] Не определена переменная окружения TASK_JSON", 2);

    let task;
    try {
        task = typeof taskJson === "object" ? taskJson : JSON.parse(taskJson);
    } catch (e) {
        errorExit("[runner/js] Ошибка разбора TASK_JSON: " + e.message, 2);
    }

    const { function: functionName, args = [], jsPath } = task;
    if (!jsPath)
        errorExit("[runner/js] Не указан путь jsPath в TASK_JSON", 3);
    if (!functionName)
        errorExit("[runner/js] Не указана функция function в TASK_JSON", 4);

    let mod;
    try {
        const jsAbsPath = path.resolve(jsPath);
        if (!fs.existsSync(jsAbsPath)) {
            errorExit(`[runner/js] JS-модуль не найден по пути: ${jsAbsPath}`, 5);
        }
        mod = require(jsAbsPath);
    } catch (e) {
        errorExit(`[runner/js] Ошибка загрузки JS-модуля: ${e.message}`, 5);
    }

    const fn = mod[functionName];
    if (typeof fn !== "function")
        errorExit(`[runner/js] Функция '${functionName}' не экспортируется JS-модулем`, 6);

    try {
        fn(...args);
    } catch (e) {
        errorExit(`[runner/js] Ошибка выполнения функции '${functionName}': ${e.message}`, 7);
    }
}

main();
