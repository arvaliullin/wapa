-- Среднее значение (mean)
CREATE OR REPLACE VIEW composer.v_metric_mean_json AS
SELECT jsonb_build_object(
    'metric', 'mean',
    'arch', arch,
    'results', jsonb_agg(r)
) AS data
FROM (
    SELECT
        e.arch,
        jsonb_build_object(
            'name', fr.function_name,
            'cpp',   MAX(CASE WHEN d.lang = 'cpp'        THEN m.mean END),
            'go',    MAX(CASE WHEN d.lang = 'go'         THEN m.mean END),
            'rust',  MAX(CASE WHEN d.lang = 'rust'       THEN m.mean END),
            'javascript', MAX(CASE WHEN d.lang = 'javascript' THEN m.mean END)
        ) AS r
    FROM
        composer.function_result fr
        INNER JOIN composer.metric m ON m.function_result_id = fr.id
        INNER JOIN composer.experiment e ON fr.experiment_id = e.id
        INNER JOIN composer.design d ON e.design_id = d.id
    GROUP BY e.arch, fr.function_name
) tmp
GROUP BY arch;

-- Стандартное отклонение (stddev)
CREATE OR REPLACE VIEW composer.v_metric_stddev_json AS
SELECT jsonb_build_object(
    'metric', 'stddev',
    'arch', arch,
    'results', jsonb_agg(r)
) AS data
FROM (
    SELECT
        e.arch,
        jsonb_build_object(
            'name', fr.function_name,
            'cpp',   MAX(CASE WHEN d.lang = 'cpp'        THEN m.stddev END),
            'go',    MAX(CASE WHEN d.lang = 'go'         THEN m.stddev END),
            'rust',  MAX(CASE WHEN d.lang = 'rust'       THEN m.stddev END),
            'javascript', MAX(CASE WHEN d.lang = 'javascript' THEN m.stddev END)
        ) AS r
    FROM
        composer.function_result fr
        INNER JOIN composer.metric m ON m.function_result_id = fr.id
        INNER JOIN composer.experiment e ON fr.experiment_id = e.id
        INNER JOIN composer.design d ON e.design_id = d.id
    GROUP BY e.arch, fr.function_name
) tmp
GROUP BY arch;

-- Медиана (median)
CREATE OR REPLACE VIEW composer.v_metric_median_json AS
SELECT jsonb_build_object(
    'metric', 'median',
    'arch', arch,
    'results', jsonb_agg(r)
) AS data
FROM (
    SELECT
        e.arch,
        jsonb_build_object(
            'name', fr.function_name,
            'cpp',   MAX(CASE WHEN d.lang = 'cpp'        THEN m.median END),
            'go',    MAX(CASE WHEN d.lang = 'go'         THEN m.median END),
            'rust',  MAX(CASE WHEN d.lang = 'rust'       THEN m.median END),
            'javascript', MAX(CASE WHEN d.lang = 'javascript' THEN m.median END)
        ) AS r
    FROM
        composer.function_result fr
        INNER JOIN composer.metric m ON m.function_result_id = fr.id
        INNER JOIN composer.experiment e ON fr.experiment_id = e.id
        INNER JOIN composer.design d ON e.design_id = d.id
    GROUP BY e.arch, fr.function_name
) tmp
GROUP BY arch;

-- Минимальные значения (min)
CREATE OR REPLACE VIEW composer.v_metric_min_json AS
SELECT jsonb_build_object(
    'metric', 'min',
    'arch', arch,
    'results', jsonb_agg(r)
) AS data
FROM (
    SELECT
        e.arch,
        jsonb_build_object(
            'name', fr.function_name,
            'cpp',   MAX(CASE WHEN d.lang = 'cpp'        THEN m.min END),
            'go',    MAX(CASE WHEN d.lang = 'go'         THEN m.min END),
            'rust',  MAX(CASE WHEN d.lang = 'rust'       THEN m.min END),
            'javascript', MAX(CASE WHEN d.lang = 'javascript' THEN m.min END)
        ) AS r
    FROM
        composer.function_result fr
        INNER JOIN composer.metric m ON m.function_result_id = fr.id
        INNER JOIN composer.experiment e ON fr.experiment_id = e.id
        INNER JOIN composer.design d ON e.design_id = d.id
    GROUP BY e.arch, fr.function_name
) tmp
GROUP BY arch;

-- Максимальные значения (max)
CREATE OR REPLACE VIEW composer.v_metric_max_json AS
SELECT jsonb_build_object(
    'metric', 'max',
    'arch', arch,
    'results', jsonb_agg(r)
) AS data
FROM (
    SELECT
        e.arch,
        jsonb_build_object(
            'name', fr.function_name,
            'cpp',   MAX(CASE WHEN d.lang = 'cpp'        THEN m.max END),
            'go',    MAX(CASE WHEN d.lang = 'go'         THEN m.max END),
            'rust',  MAX(CASE WHEN d.lang = 'rust'       THEN m.max END),
            'javascript', MAX(CASE WHEN d.lang = 'javascript' THEN m.max END)
        ) AS r
    FROM
        composer.function_result fr
        INNER JOIN composer.metric m ON m.function_result_id = fr.id
        INNER JOIN composer.experiment e ON fr.experiment_id = e.id
        INNER JOIN composer.design d ON e.design_id = d.id
    GROUP BY e.arch, fr.function_name
) tmp
GROUP BY arch;
