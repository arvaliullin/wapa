-- 1. Зависимость времени выполнения от языка реализации
CREATE OR REPLACE VIEW composer.v_mean_time_by_lang AS
SELECT
    d.lang,
    AVG(m.mean)   AS avg_mean_time,
    AVG(m.median) AS avg_median_time,
    COUNT(*)      AS num_experiments
FROM composer.metric m
JOIN composer.function_result fr ON m.function_result_id = fr.id
JOIN composer.experiment e ON fr.experiment_id = e.id
JOIN composer.design d ON e.design_id = d.id
GROUP BY d.lang;

-- 2. Зависимость времени выполнения от архитектуры платформы
CREATE OR REPLACE VIEW composer.v_mean_time_by_arch AS
SELECT
    e.arch,
    AVG(m.mean)   AS avg_mean_time,
    AVG(m.median) AS avg_median_time,
    COUNT(*)      AS num_experiments
FROM composer.metric m
JOIN composer.function_result fr ON m.function_result_id = fr.id
JOIN composer.experiment e ON fr.experiment_id = e.id
GROUP BY e.arch;

-- 3. Зависимость времени выполнения от типа реализации (WASM или JS)
CREATE OR REPLACE VIEW composer.v_mean_time_by_js_wasm AS
SELECT
    d.lang,
    CASE 
        WHEN d.wasm IS NULL THEN 'js'
        ELSE 'wasm'
    END AS implementation_type,
    AVG(m.mean)   AS avg_mean_time,
    AVG(m.median) AS avg_median_time,
    COUNT(*)      AS num_experiments
FROM composer.metric m
JOIN composer.function_result fr ON m.function_result_id = fr.id
JOIN composer.experiment e ON fr.experiment_id = e.id
JOIN composer.design d ON e.design_id = d.id
GROUP BY d.lang, implementation_type;

CREATE OR REPLACE VIEW composer.v_js_wasm_speedup AS
WITH perf AS (
    SELECT
        fr.function_name,
        fr.args,
        CASE WHEN d.wasm IS NULL THEN 'js' ELSE 'wasm' END AS implementation_type,
        AVG(m.mean) AS avg_mean_time
    FROM composer.metric m
    JOIN composer.function_result fr ON m.function_result_id = fr.id
    JOIN composer.experiment e ON fr.experiment_id = e.id
    JOIN composer.design d ON e.design_id = d.id
    GROUP BY fr.function_name, fr.args, implementation_type
)
SELECT
    js.function_name,
    js.args,
    js.avg_mean_time AS js_mean_time,
    wasm.avg_mean_time AS wasm_mean_time,
    CASE WHEN wasm.avg_mean_time IS NOT NULL AND wasm.avg_mean_time != 0
         THEN js.avg_mean_time / wasm.avg_mean_time
         ELSE NULL
    END AS speedup
FROM perf js
JOIN perf wasm
    ON js.function_name = wasm.function_name
    AND js.args = wasm.args
    AND js.implementation_type = 'js'
    AND wasm.implementation_type = 'wasm'
ORDER BY js.function_name, js.args;

-- 5. Совместное влияние языка и архитектуры
CREATE OR REPLACE VIEW composer.v_time_by_lang_arch AS
SELECT
    d.lang,
    e.arch,
    AVG(m.mean)   AS avg_mean_time,
    AVG(m.median) AS avg_median_time,
    COUNT(*)      AS num_experiments
FROM composer.metric m
JOIN composer.function_result fr ON m.function_result_id = fr.id
JOIN composer.experiment e ON fr.experiment_id = e.id
JOIN composer.design d ON e.design_id = d.id
GROUP BY d.lang, e.arch
ORDER BY d.lang, e.arch;

-- 6. Стабильность (разброс значений во времени)
CREATE OR REPLACE VIEW composer.v_func_stability AS
SELECT
    d.name AS experiment_design,
    fr.function_name,
    AVG(m.stddev) AS avg_stddev,
    AVG(m.mean)   AS avg_mean,
    (AVG(m.stddev) / NULLIF(AVG(m.mean),0)) AS rel_variation, -- относительная вариация (чем меньше, тем стабильнее)
    COUNT(*) AS num_measurements
FROM composer.metric m
JOIN composer.function_result fr ON m.function_result_id = fr.id
JOIN composer.experiment e ON fr.experiment_id = e.id
JOIN composer.design d ON e.design_id = d.id
GROUP BY d.name, fr.function_name
ORDER BY rel_variation;

-- 7. Сравнение производительности разных языков на одной задаче и одних аргументах
CREATE OR REPLACE VIEW composer.v_lang_compare_per_task AS
SELECT
    fr.function_name,
    fr.args,
    d.lang,
    AVG(m.mean) AS avg_mean_time
FROM composer.metric m
JOIN composer.function_result fr ON m.function_result_id = fr.id
JOIN composer.experiment e ON fr.experiment_id = e.id
JOIN composer.design d ON e.design_id = d.id
GROUP BY fr.function_name, fr.args, d.lang
ORDER BY fr.function_name, fr.args, avg_mean_time;

CREATE OR REPLACE VIEW composer.function_result_full AS
SELECT
    fr.id AS function_result_id,
    fr.function_name,
    fr.args,
    fr.repeats,
    fr.result,
    e.id AS experiment_id,
    d.id AS design_id,
    d.name AS design_name,
    d.lang AS design_lang,
    d.js AS design_js,
    d.wasm AS design_wasm,
    d.functions AS design_functions,
    e.hostname AS experiment_hostname,
    e.arch AS experiment_arch
FROM
    composer.function_result fr
    INNER JOIN composer.experiment e ON fr.experiment_id = e.id
    INNER JOIN composer.design d ON e.design_id = d.id;
