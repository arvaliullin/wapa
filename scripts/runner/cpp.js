// scripts/runner/cpp.js

// Получаем JSON-данные из переменной окружения
const rawData = process.env.JSON_DATA;

if (!rawData) {
    console.error("Environment variable JSON_DATA is not set.");
    process.exit(1);
}

console.log("Raw input from environment:", rawData);

try {
    // Парсим JSON
    const obj = JSON.parse(rawData);
    console.log("Parsed object:", obj);

    // Функция для задержки
    function sleep(ms) {
        return new Promise(resolve => setTimeout(resolve, ms));
    }

    // Пример асинхронной функции
    async function example() {
        console.log("Start sleeping...");
        await sleep(1000);
        console.log("Finished sleeping!");
    }

    example();
} catch (error) {
    console.error("Invalid JSON in environment variable JSON_DATA:", error.message);
    process.exit(1);
}
