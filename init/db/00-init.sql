CREATE SCHEMA IF NOT EXISTS composer;

CREATE EXTENSION IF NOT EXISTS pgcrypto SCHEMA composer;

-- Таблица для хранения плана эксперимента
CREATE TABLE composer.design (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL UNIQUE,
    lang TEXT NOT NULL CHECK (lang IN ('go', 'javascript', 'rust', 'cpp', 'tinygo')),
    js TEXT NOT NULL,
    wasm TEXT,
    functions JSONB NOT NULL
);

COMMENT ON TABLE composer.design IS 'План эксперимента';
COMMENT ON COLUMN composer.design.id IS 'Идентификатор плана эксперимента';
COMMENT ON COLUMN composer.design.name IS 'Имя эксперимента';
COMMENT ON COLUMN composer.design.lang IS 'Язык реализации модуля';
COMMENT ON COLUMN composer.design.js IS 'Имя файла javascript';
COMMENT ON COLUMN composer.design.wasm IS 'Имя модуля WebAssembly';
COMMENT ON COLUMN composer.design.functions IS 'Массив функций и их аргументов';

-- Таблица для хранения экспериментов
CREATE TABLE composer.experiment (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    design_id UUID REFERENCES composer.design(id) ON DELETE CASCADE,
    hostname TEXT NOT NULL,
    arch TEXT NOT NULL
);

COMMENT ON TABLE composer.experiment IS 'Таблица для хранения информации о выполнении экспериментов';
COMMENT ON COLUMN composer.experiment.id IS 'Уникальный идентификатор эксперимента';
COMMENT ON COLUMN composer.experiment.design_id IS 'План эксперимента';
COMMENT ON COLUMN composer.experiment.hostname IS 'Имя машины';
COMMENT ON COLUMN composer.experiment.arch IS 'Архитектура системы';

CREATE TABLE composer.function_result (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    experiment_id UUID REFERENCES composer.experiment(id) ON DELETE CASCADE,
    function_name TEXT NOT NULL,
    args JSONB NOT NULL,
    repeats INTEGER NOT NULL,
    result JSONB NOT NULL 
);

COMMENT ON TABLE composer.function_result IS 'Результаты выполнения отдельных функций в рамках экспериментов';
COMMENT ON COLUMN composer.function_result.experiment_id IS 'Идентификатор эксперимента';
COMMENT ON COLUMN composer.function_result.function_name IS 'Имя функции';
COMMENT ON COLUMN composer.function_result.args IS 'Аргументы, используемые для запуска функции';
COMMENT ON COLUMN composer.function_result.repeats IS 'Количество повторов выполнения функции';
COMMENT ON COLUMN composer.function_result.result IS 'Результаты выполнения функции (JSON)';

-- Таблица для хранения метрик по функциям
CREATE TABLE composer.metric (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    function_result_id UUID REFERENCES composer.function_result(id) ON DELETE CASCADE,
    mean DOUBLE PRECISION NOT NULL,
    stddev DOUBLE PRECISION NOT NULL,
    median DOUBLE PRECISION NOT NULL,
    user_time DOUBLE PRECISION NOT NULL,
    system DOUBLE PRECISION NOT NULL,
    min DOUBLE PRECISION NOT NULL,
    max DOUBLE PRECISION NOT NULL
);

-- Метрики производительности (времени выполнения), полученные с помощью утилиты Hyperfine для функций
COMMENT ON TABLE composer.metric IS 'Метрики производительности (результаты Hyperfine) для функций';
COMMENT ON COLUMN composer.metric.function_result_id IS 'Внешний ключ на результаты выполнения конкретной функции';
COMMENT ON COLUMN composer.metric.mean IS 'Среднее арифметическое время выполнения, полученное Hyperfine (секунды)';
COMMENT ON COLUMN composer.metric.stddev IS 'Стандартное отклонение времени выполнения, вычисленное Hyperfine (секунды)';
COMMENT ON COLUMN composer.metric.median IS 'Медианное время выполнения (секунды)';
COMMENT ON COLUMN composer.metric.user_time IS 'Среднее пользовательское время (user time), оцененное Hyperfine (секунды)';
COMMENT ON COLUMN composer.metric.system IS 'Среднее системное время (system time), оцененное Hyperfine (секунды)';
COMMENT ON COLUMN composer.metric.min IS 'Минимальное измеренное время выполнения (секунды)';
COMMENT ON COLUMN composer.metric.max IS 'Максимальное измеренное время выполнения (секунды)';
