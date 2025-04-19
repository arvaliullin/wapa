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

-- Таблица для хранения информации о системе
CREATE TABLE composer.system_info (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hostname TEXT NOT NULL UNIQUE,
    platform TEXT NOT NULL,
    arch TEXT NOT NULL,
    cpu_count INTEGER NOT NULL,
    node_version TEXT NOT NULL
);

COMMENT ON TABLE composer.system_info IS 'Системная информация для узлов выполнения экспериментов';
COMMENT ON COLUMN composer.system_info.hostname IS 'Уникальное имя хоста';
COMMENT ON COLUMN composer.system_info.platform IS 'Платформа операционной системы';
COMMENT ON COLUMN composer.system_info.arch IS 'Архитектура процессора';
COMMENT ON COLUMN composer.system_info.cpu_count IS 'Количество доступных процессорных ядер';
COMMENT ON COLUMN composer.system_info.node_version IS 'Версия Node.js среды выполнения';

-- Таблица для хранения экспериментов
CREATE TABLE composer.experiment (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    design_id UUID REFERENCES composer.design(id) ON DELETE CASCADE,
    system_id UUID REFERENCES composer.system_info(id),
    execution_time BIGINT NOT NULL
);

COMMENT ON TABLE composer.experiment IS 'Таблица для хранения информации о выполнении экспериментов';
COMMENT ON COLUMN composer.experiment.id IS 'Уникальный идентификатор эксперимента';
COMMENT ON COLUMN composer.experiment.design_id IS 'План эксперимента';
COMMENT ON COLUMN composer.experiment.system_id IS 'Ссылка на системную информацию';
COMMENT ON COLUMN composer.experiment.execution_time IS 'Время выполнения кода';

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
    median DOUBLE PRECISION NOT NULL,
    min DOUBLE PRECISION NOT NULL,
    max DOUBLE PRECISION NOT NULL,
    variance DOUBLE PRECISION,
    std_deviation DOUBLE PRECISION
);

COMMENT ON TABLE composer.metric IS 'Метрики производительности для функций';
COMMENT ON COLUMN composer.metric.function_result_id IS 'Ссылка на результаты выполнения функции';
COMMENT ON COLUMN composer.metric.mean IS 'Среднее время выполнения';
COMMENT ON COLUMN composer.metric.median IS 'Медианное время выполнения';
COMMENT ON COLUMN composer.metric.min IS 'Минимальное время выполнения';
COMMENT ON COLUMN composer.metric.max IS 'Максимальное время выполнения';
COMMENT ON COLUMN composer.metric.variance IS 'Дисперсия';
COMMENT ON COLUMN composer.metric.std_deviation IS 'Стандартное отклонение';
