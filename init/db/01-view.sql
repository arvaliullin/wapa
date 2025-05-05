CREATE OR REPLACE VIEW composer.vw_func_metrics_by_arch_lang AS
SELECT 
    e.arch,
    fr.function_name AS name,
    d.lang,
    AVG(m.mean) AS mean,
    AVG(m.stddev) AS stddev,
    AVG(m.median) AS median,
    AVG(m.min) AS min,
    AVG(m.max) AS max,
    AVG(m.user_time) AS user_time,
    AVG(m.system) AS system
FROM composer.metric m
JOIN composer.function_result fr ON fr.id = m.function_result_id
JOIN composer.experiment e ON e.id = fr.experiment_id
JOIN composer.design d ON d.id = e.design_id
GROUP BY e.arch, fr.function_name, d.lang;

CREATE OR REPLACE VIEW composer.vw_mock_func_metrics_by_arch_lang AS
SELECT 
    e.arch,
    fr.function_name AS name,
    d.lang,
    AVG(m.mean) AS mean,
    AVG(m.stddev) AS stddev,
    AVG(m.median) AS median,
    AVG(m.min) AS min,
    AVG(m.max) AS max,
    AVG(m.user_time) AS user_time,
    AVG(m.system) AS system
FROM composer.metric m
JOIN composer.function_result fr ON fr.id = m.function_result_id
JOIN composer.experiment e ON e.id = fr.experiment_id
JOIN composer.design d ON d.id = e.design_id
WHERE fr.function_name LIKE '%Mock'
GROUP BY e.arch, fr.function_name, d.lang;

CREATE OR REPLACE VIEW composer.vw_func_delta_metrics_mock_by_arch_lang AS
SELECT
    orig.arch,
    'd' || orig.name || 'Mock' AS name,
    orig.lang,
    ABS(orig.mean - mock.mean) AS d_mean,
    ABS(orig.median - mock.median) AS d_median,
    ABS(orig.stddev - mock.stddev) AS d_stddev,
    ABS(orig.min - mock.min) AS d_min,
    ABS(orig.max - mock.max) AS d_max,
    ABS(orig.user_time - mock.user_time) AS d_user_time,
    ABS(orig.system - mock.system) AS d_system
FROM composer.vw_func_metrics_by_arch_lang orig
JOIN composer.vw_mock_func_metrics_by_arch_lang mock 
  ON orig.arch = mock.arch
     AND orig.lang = mock.lang
     AND orig.name || 'Mock' = mock.name
WHERE orig.name NOT LIKE '%Mock';

CREATE OR REPLACE VIEW composer.vw_fastest_functions_by_arch_lang AS
SELECT *
FROM (
    SELECT
        e.arch,
        fr.function_name,
        d.lang,
        AVG(m.mean) AS mean,
        ROW_NUMBER() OVER (PARTITION BY e.arch, d.lang ORDER BY AVG(m.mean)) AS rn
    FROM composer.metric m
    JOIN composer.function_result fr ON fr.id = m.function_result_id
    JOIN composer.experiment e ON e.id = fr.experiment_id
    JOIN composer.design d ON d.id = e.design_id
    GROUP BY e.arch, fr.function_name, d.lang
) t
WHERE rn <= 5;

CREATE OR REPLACE VIEW composer.vw_func_dispersion_by_arch_lang AS
SELECT
    e.arch,
    fr.function_name,
    d.lang,
    MIN(m.min) AS abs_min,
    MAX(m.max) AS abs_max
FROM composer.metric m
JOIN composer.function_result fr ON fr.id = m.function_result_id
JOIN composer.experiment e ON e.id = fr.experiment_id
JOIN composer.design d ON d.id = e.design_id
GROUP BY e.arch, fr.function_name, d.lang;
