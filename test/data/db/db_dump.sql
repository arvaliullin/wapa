--
-- PostgreSQL database dump
--

-- Dumped from database version 17.5 (Debian 17.5-1.pgdg120+1)
-- Dumped by pg_dump version 17.5 (Debian 17.5-1.pgdg120+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: composer; Type: SCHEMA; Schema: -; Owner: postgres
--

CREATE SCHEMA composer;


ALTER SCHEMA composer OWNER TO postgres;

--
-- Name: pgcrypto; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS pgcrypto WITH SCHEMA composer;


--
-- Name: EXTENSION pgcrypto; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION pgcrypto IS 'cryptographic functions';


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: design; Type: TABLE; Schema: composer; Owner: postgres
--

CREATE TABLE composer.design (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    name text NOT NULL,
    lang text NOT NULL,
    js text NOT NULL,
    wasm text,
    functions jsonb NOT NULL,
    CONSTRAINT design_lang_check CHECK ((lang = ANY (ARRAY['go'::text, 'javascript'::text, 'rust'::text, 'cpp'::text, 'tinygo'::text])))
);


ALTER TABLE composer.design OWNER TO postgres;

--
-- Name: TABLE design; Type: COMMENT; Schema: composer; Owner: postgres
--

COMMENT ON TABLE composer.design IS 'План эксперимента';


--
-- Name: COLUMN design.id; Type: COMMENT; Schema: composer; Owner: postgres
--

COMMENT ON COLUMN composer.design.id IS 'Идентификатор плана эксперимента';


--
-- Name: COLUMN design.name; Type: COMMENT; Schema: composer; Owner: postgres
--

COMMENT ON COLUMN composer.design.name IS 'Имя эксперимента';


--
-- Name: COLUMN design.lang; Type: COMMENT; Schema: composer; Owner: postgres
--

COMMENT ON COLUMN composer.design.lang IS 'Язык реализации модуля';


--
-- Name: COLUMN design.js; Type: COMMENT; Schema: composer; Owner: postgres
--

COMMENT ON COLUMN composer.design.js IS 'Имя файла javascript';


--
-- Name: COLUMN design.wasm; Type: COMMENT; Schema: composer; Owner: postgres
--

COMMENT ON COLUMN composer.design.wasm IS 'Имя модуля WebAssembly';


--
-- Name: COLUMN design.functions; Type: COMMENT; Schema: composer; Owner: postgres
--

COMMENT ON COLUMN composer.design.functions IS 'Массив функций и их аргументов';


--
-- Name: experiment; Type: TABLE; Schema: composer; Owner: postgres
--

CREATE TABLE composer.experiment (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    design_id uuid,
    hostname text NOT NULL,
    arch text NOT NULL
);


ALTER TABLE composer.experiment OWNER TO postgres;

--
-- Name: TABLE experiment; Type: COMMENT; Schema: composer; Owner: postgres
--

COMMENT ON TABLE composer.experiment IS 'Таблица для хранения информации о выполнении экспериментов';


--
-- Name: COLUMN experiment.id; Type: COMMENT; Schema: composer; Owner: postgres
--

COMMENT ON COLUMN composer.experiment.id IS 'Уникальный идентификатор эксперимента';


--
-- Name: COLUMN experiment.design_id; Type: COMMENT; Schema: composer; Owner: postgres
--

COMMENT ON COLUMN composer.experiment.design_id IS 'План эксперимента';


--
-- Name: COLUMN experiment.hostname; Type: COMMENT; Schema: composer; Owner: postgres
--

COMMENT ON COLUMN composer.experiment.hostname IS 'Имя машины';


--
-- Name: COLUMN experiment.arch; Type: COMMENT; Schema: composer; Owner: postgres
--

COMMENT ON COLUMN composer.experiment.arch IS 'Архитектура системы';


--
-- Name: function_result; Type: TABLE; Schema: composer; Owner: postgres
--

CREATE TABLE composer.function_result (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    experiment_id uuid,
    function_name text NOT NULL,
    args jsonb NOT NULL,
    repeats integer NOT NULL,
    result jsonb NOT NULL
);


ALTER TABLE composer.function_result OWNER TO postgres;

--
-- Name: TABLE function_result; Type: COMMENT; Schema: composer; Owner: postgres
--

COMMENT ON TABLE composer.function_result IS 'Результаты выполнения отдельных функций в рамках экспериментов';


--
-- Name: COLUMN function_result.experiment_id; Type: COMMENT; Schema: composer; Owner: postgres
--

COMMENT ON COLUMN composer.function_result.experiment_id IS 'Идентификатор эксперимента';


--
-- Name: COLUMN function_result.function_name; Type: COMMENT; Schema: composer; Owner: postgres
--

COMMENT ON COLUMN composer.function_result.function_name IS 'Имя функции';


--
-- Name: COLUMN function_result.args; Type: COMMENT; Schema: composer; Owner: postgres
--

COMMENT ON COLUMN composer.function_result.args IS 'Аргументы, используемые для запуска функции';


--
-- Name: COLUMN function_result.repeats; Type: COMMENT; Schema: composer; Owner: postgres
--

COMMENT ON COLUMN composer.function_result.repeats IS 'Количество повторов выполнения функции';


--
-- Name: COLUMN function_result.result; Type: COMMENT; Schema: composer; Owner: postgres
--

COMMENT ON COLUMN composer.function_result.result IS 'Результаты выполнения функции (JSON)';


--
-- Name: metric; Type: TABLE; Schema: composer; Owner: postgres
--

CREATE TABLE composer.metric (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    function_result_id uuid,
    mean double precision NOT NULL,
    stddev double precision NOT NULL,
    median double precision NOT NULL,
    user_time double precision NOT NULL,
    system double precision NOT NULL,
    min double precision NOT NULL,
    max double precision NOT NULL
);


ALTER TABLE composer.metric OWNER TO postgres;

--
-- Name: TABLE metric; Type: COMMENT; Schema: composer; Owner: postgres
--

COMMENT ON TABLE composer.metric IS 'Метрики производительности (результаты Hyperfine) для функций';


--
-- Name: COLUMN metric.function_result_id; Type: COMMENT; Schema: composer; Owner: postgres
--

COMMENT ON COLUMN composer.metric.function_result_id IS 'Внешний ключ на результаты выполнения конкретной функции';


--
-- Name: COLUMN metric.mean; Type: COMMENT; Schema: composer; Owner: postgres
--

COMMENT ON COLUMN composer.metric.mean IS 'Среднее арифметическое время выполнения, полученное Hyperfine (секунды)';


--
-- Name: COLUMN metric.stddev; Type: COMMENT; Schema: composer; Owner: postgres
--

COMMENT ON COLUMN composer.metric.stddev IS 'Стандартное отклонение времени выполнения, вычисленное Hyperfine (секунды)';


--
-- Name: COLUMN metric.median; Type: COMMENT; Schema: composer; Owner: postgres
--

COMMENT ON COLUMN composer.metric.median IS 'Медианное время выполнения (секунды)';


--
-- Name: COLUMN metric.user_time; Type: COMMENT; Schema: composer; Owner: postgres
--

COMMENT ON COLUMN composer.metric.user_time IS 'Среднее пользовательское время (user time), оцененное Hyperfine (секунды)';


--
-- Name: COLUMN metric.system; Type: COMMENT; Schema: composer; Owner: postgres
--

COMMENT ON COLUMN composer.metric.system IS 'Среднее системное время (system time), оцененное Hyperfine (секунды)';


--
-- Name: COLUMN metric.min; Type: COMMENT; Schema: composer; Owner: postgres
--

COMMENT ON COLUMN composer.metric.min IS 'Минимальное измеренное время выполнения (секунды)';


--
-- Name: COLUMN metric.max; Type: COMMENT; Schema: composer; Owner: postgres
--

COMMENT ON COLUMN composer.metric.max IS 'Максимальное измеренное время выполнения (секунды)';


--
-- Name: v_metric_max_json; Type: VIEW; Schema: composer; Owner: postgres
--

CREATE VIEW composer.v_metric_max_json AS
 SELECT jsonb_build_object('metric', 'max', 'arch', arch, 'results', jsonb_agg(r)) AS data
   FROM ( SELECT e.arch,
            jsonb_build_object('name', fr.function_name, 'cpp', max(
                CASE
                    WHEN (d.lang = 'cpp'::text) THEN m.max
                    ELSE NULL::double precision
                END), 'go', max(
                CASE
                    WHEN (d.lang = 'go'::text) THEN m.max
                    ELSE NULL::double precision
                END), 'rust', max(
                CASE
                    WHEN (d.lang = 'rust'::text) THEN m.max
                    ELSE NULL::double precision
                END), 'javascript', max(
                CASE
                    WHEN (d.lang = 'javascript'::text) THEN m.max
                    ELSE NULL::double precision
                END)) AS r
           FROM (((composer.function_result fr
             JOIN composer.metric m ON ((m.function_result_id = fr.id)))
             JOIN composer.experiment e ON ((fr.experiment_id = e.id)))
             JOIN composer.design d ON ((e.design_id = d.id)))
          GROUP BY e.arch, fr.function_name) tmp
  GROUP BY arch;


ALTER VIEW composer.v_metric_max_json OWNER TO postgres;

--
-- Name: v_metric_mean_json; Type: VIEW; Schema: composer; Owner: postgres
--

CREATE VIEW composer.v_metric_mean_json AS
 SELECT jsonb_build_object('metric', 'mean', 'arch', arch, 'results', jsonb_agg(r)) AS data
   FROM ( SELECT e.arch,
            jsonb_build_object('name', fr.function_name, 'cpp', max(
                CASE
                    WHEN (d.lang = 'cpp'::text) THEN m.mean
                    ELSE NULL::double precision
                END), 'go', max(
                CASE
                    WHEN (d.lang = 'go'::text) THEN m.mean
                    ELSE NULL::double precision
                END), 'rust', max(
                CASE
                    WHEN (d.lang = 'rust'::text) THEN m.mean
                    ELSE NULL::double precision
                END), 'javascript', max(
                CASE
                    WHEN (d.lang = 'javascript'::text) THEN m.mean
                    ELSE NULL::double precision
                END)) AS r
           FROM (((composer.function_result fr
             JOIN composer.metric m ON ((m.function_result_id = fr.id)))
             JOIN composer.experiment e ON ((fr.experiment_id = e.id)))
             JOIN composer.design d ON ((e.design_id = d.id)))
          GROUP BY e.arch, fr.function_name) tmp
  GROUP BY arch;


ALTER VIEW composer.v_metric_mean_json OWNER TO postgres;

--
-- Name: v_metric_median_json; Type: VIEW; Schema: composer; Owner: postgres
--

CREATE VIEW composer.v_metric_median_json AS
 SELECT jsonb_build_object('metric', 'median', 'arch', arch, 'results', jsonb_agg(r)) AS data
   FROM ( SELECT e.arch,
            jsonb_build_object('name', fr.function_name, 'cpp', max(
                CASE
                    WHEN (d.lang = 'cpp'::text) THEN m.median
                    ELSE NULL::double precision
                END), 'go', max(
                CASE
                    WHEN (d.lang = 'go'::text) THEN m.median
                    ELSE NULL::double precision
                END), 'rust', max(
                CASE
                    WHEN (d.lang = 'rust'::text) THEN m.median
                    ELSE NULL::double precision
                END), 'javascript', max(
                CASE
                    WHEN (d.lang = 'javascript'::text) THEN m.median
                    ELSE NULL::double precision
                END)) AS r
           FROM (((composer.function_result fr
             JOIN composer.metric m ON ((m.function_result_id = fr.id)))
             JOIN composer.experiment e ON ((fr.experiment_id = e.id)))
             JOIN composer.design d ON ((e.design_id = d.id)))
          GROUP BY e.arch, fr.function_name) tmp
  GROUP BY arch;


ALTER VIEW composer.v_metric_median_json OWNER TO postgres;

--
-- Name: v_metric_min_json; Type: VIEW; Schema: composer; Owner: postgres
--

CREATE VIEW composer.v_metric_min_json AS
 SELECT jsonb_build_object('metric', 'min', 'arch', arch, 'results', jsonb_agg(r)) AS data
   FROM ( SELECT e.arch,
            jsonb_build_object('name', fr.function_name, 'cpp', max(
                CASE
                    WHEN (d.lang = 'cpp'::text) THEN m.min
                    ELSE NULL::double precision
                END), 'go', max(
                CASE
                    WHEN (d.lang = 'go'::text) THEN m.min
                    ELSE NULL::double precision
                END), 'rust', max(
                CASE
                    WHEN (d.lang = 'rust'::text) THEN m.min
                    ELSE NULL::double precision
                END), 'javascript', max(
                CASE
                    WHEN (d.lang = 'javascript'::text) THEN m.min
                    ELSE NULL::double precision
                END)) AS r
           FROM (((composer.function_result fr
             JOIN composer.metric m ON ((m.function_result_id = fr.id)))
             JOIN composer.experiment e ON ((fr.experiment_id = e.id)))
             JOIN composer.design d ON ((e.design_id = d.id)))
          GROUP BY e.arch, fr.function_name) tmp
  GROUP BY arch;


ALTER VIEW composer.v_metric_min_json OWNER TO postgres;

--
-- Name: v_metric_stddev_json; Type: VIEW; Schema: composer; Owner: postgres
--

CREATE VIEW composer.v_metric_stddev_json AS
 SELECT jsonb_build_object('metric', 'stddev', 'arch', arch, 'results', jsonb_agg(r)) AS data
   FROM ( SELECT e.arch,
            jsonb_build_object('name', fr.function_name, 'cpp', max(
                CASE
                    WHEN (d.lang = 'cpp'::text) THEN m.stddev
                    ELSE NULL::double precision
                END), 'go', max(
                CASE
                    WHEN (d.lang = 'go'::text) THEN m.stddev
                    ELSE NULL::double precision
                END), 'rust', max(
                CASE
                    WHEN (d.lang = 'rust'::text) THEN m.stddev
                    ELSE NULL::double precision
                END), 'javascript', max(
                CASE
                    WHEN (d.lang = 'javascript'::text) THEN m.stddev
                    ELSE NULL::double precision
                END)) AS r
           FROM (((composer.function_result fr
             JOIN composer.metric m ON ((m.function_result_id = fr.id)))
             JOIN composer.experiment e ON ((fr.experiment_id = e.id)))
             JOIN composer.design d ON ((e.design_id = d.id)))
          GROUP BY e.arch, fr.function_name) tmp
  GROUP BY arch;


ALTER VIEW composer.v_metric_stddev_json OWNER TO postgres;

--
-- Data for Name: design; Type: TABLE DATA; Schema: composer; Owner: postgres
--

COPY composer.design (id, name, lang, js, wasm, functions) FROM stdin;
6ce1e232-f6c2-48bd-8ef0-55271378b4b4	js_exps	javascript	lib_js.out.js		[{"args": [0, 100, 10000], "function": "x2Integrate"}, {"args": [0, 100, 10000], "function": "x2IntegrateMock"}, {"args": [35], "function": "fibonacciRecursive"}, {"args": [35], "function": "fibonacciRecursiveMock"}, {"args": [35], "function": "fibonacciIterative"}, {"args": [35], "function": "fibonacciIterativeMock"}, {"args": [1000000], "function": "multiply"}, {"args": [1000000], "function": "multiplyMock"}, {"args": [10000], "function": "multiplyVector"}, {"args": [10000], "function": "multiplyVectorMock"}, {"args": [56], "function": "factorize"}, {"args": [56], "function": "factorizeMock"}]
6d23fcd5-d040-4191-982e-f557d08f720f	cpp_exps	cpp	lib_cpp.out.js	lib_cpp.out.wasm	[{"args": [0, 100, 10000], "function": "x2Integrate"}, {"args": [0, 100, 10000], "function": "x2IntegrateMock"}, {"args": [35], "function": "fibonacciRecursive"}, {"args": [35], "function": "fibonacciRecursiveMock"}, {"args": [35], "function": "fibonacciIterative"}, {"args": [35], "function": "fibonacciIterativeMock"}, {"args": [1000000], "function": "multiply"}, {"args": [1000000], "function": "multiplyMock"}, {"args": [10000], "function": "multiplyVector"}, {"args": [10000], "function": "multiplyVectorMock"}, {"args": [56], "function": "factorize"}, {"args": [56], "function": "factorizeMock"}]
b742f7a8-1d53-4f7a-9f36-cdee60865463	go_exps	go		integrate.wasm	[{"args": [0, 100, 10000], "function": "x2Integrate"}, {"args": [0, 100, 10000], "function": "x2IntegrateMock"}, {"args": [35], "function": "fibonacciRecursive"}, {"args": [35], "function": "fibonacciRecursiveMock"}, {"args": [35], "function": "fibonacciIterative"}, {"args": [35], "function": "fibonacciIterativeMock"}, {"args": [1000000], "function": "multiply"}, {"args": [1000000], "function": "multiplyMock"}, {"args": [10000], "function": "multiplyVector"}, {"args": [10000], "function": "multiplyVectorMock"}, {"args": [56], "function": "factorize"}, {"args": [56], "function": "factorizeMock"}]
e805a22a-1c18-437c-b04c-70f10e4bb7c2	rs_exps	rust	rs_exps.js	rs_exps_bg.wasm	[{"args": [0, 100, 10000], "function": "x2Integrate"}, {"args": [0, 100, 10000], "function": "x2IntegrateMock"}, {"args": [35], "function": "fibonacciRecursive"}, {"args": [35], "function": "fibonacciRecursiveMock"}, {"args": [35], "function": "fibonacciIterative"}, {"args": [35], "function": "fibonacciIterativeMock"}, {"args": [1000000], "function": "multiply"}, {"args": [1000000], "function": "multiplyMock"}, {"args": [10000], "function": "multiplyVector"}, {"args": [10000], "function": "multiplyVectorMock"}, {"args": [56], "function": "factorize"}, {"args": [56], "function": "factorizeMock"}]
\.


--
-- Data for Name: experiment; Type: TABLE DATA; Schema: composer; Owner: postgres
--

COPY composer.experiment (id, design_id, hostname, arch) FROM stdin;
82ab2b61-8e40-443f-9db7-9ae7bde370e8	6ce1e232-f6c2-48bd-8ef0-55271378b4b4	51991e72636b	amd64
a6fa7715-afc4-4228-a8cd-9ac1fb535065	6ce1e232-f6c2-48bd-8ef0-55271378b4b4	rpi3	arm64
9bdc8bf7-8e77-40c2-8f55-6bf927a2c80b	6d23fcd5-d040-4191-982e-f557d08f720f	51991e72636b	amd64
d02506eb-4d0a-48ac-a07f-415b41e89019	6d23fcd5-d040-4191-982e-f557d08f720f	rpi3	arm64
5865d7b3-1952-40c1-bb42-33a3f9dc8afc	b742f7a8-1d53-4f7a-9f36-cdee60865463	51991e72636b	amd64
18003139-c722-4abc-b4f4-f2b82d09ef90	b742f7a8-1d53-4f7a-9f36-cdee60865463	rpi3	arm64
72ff8a2c-a0d5-4022-8a24-d9cd4e7cde45	e805a22a-1c18-437c-b04c-70f10e4bb7c2	51991e72636b	amd64
b69732a4-e070-4381-9655-32883ebdfadd	e805a22a-1c18-437c-b04c-70f10e4bb7c2	rpi3	arm64
\.


--
-- Data for Name: function_result; Type: TABLE DATA; Schema: composer; Owner: postgres
--

COPY composer.function_result (id, experiment_id, function_name, args, repeats, result) FROM stdin;
64645a3c-fd52-4df6-abdf-7beecc668a17	82ab2b61-8e40-443f-9db7-9ae7bde370e8	x2Integrate	[0, 100, 10000]	1000	null
19b87222-918e-461b-a0e4-1b8525817988	82ab2b61-8e40-443f-9db7-9ae7bde370e8	x2IntegrateMock	[0, 100, 10000]	1000	null
416c71c5-04e0-4725-94d3-c9d529ec668c	82ab2b61-8e40-443f-9db7-9ae7bde370e8	fibonacciRecursive	[35]	1000	null
f932695b-b59a-4d83-83fd-7b1953ee27c7	82ab2b61-8e40-443f-9db7-9ae7bde370e8	fibonacciRecursiveMock	[35]	1000	null
3157ed1f-5aca-4a62-845f-99b1232e80ff	82ab2b61-8e40-443f-9db7-9ae7bde370e8	fibonacciIterative	[35]	1000	null
83a4b7f6-6aec-4c1e-a143-dbb256969bed	82ab2b61-8e40-443f-9db7-9ae7bde370e8	fibonacciIterativeMock	[35]	1000	null
e27dcf72-9c89-491e-9a62-f4356f4d1c33	82ab2b61-8e40-443f-9db7-9ae7bde370e8	multiply	[1000000]	1000	null
ab08fb6d-0ad8-4856-998a-32fc65d1090f	82ab2b61-8e40-443f-9db7-9ae7bde370e8	multiplyMock	[1000000]	1000	null
76b35674-82b8-494d-b383-118436e19b0c	82ab2b61-8e40-443f-9db7-9ae7bde370e8	multiplyVector	[10000]	1000	null
a7133fcf-c92b-4dc8-9df7-2a9e788763fd	82ab2b61-8e40-443f-9db7-9ae7bde370e8	multiplyVectorMock	[10000]	1000	null
88a25685-2c96-44f0-a37e-414f3f69c33d	82ab2b61-8e40-443f-9db7-9ae7bde370e8	factorize	[56]	1000	null
a5472a47-5470-4447-9f40-3694bdf40a32	82ab2b61-8e40-443f-9db7-9ae7bde370e8	factorizeMock	[56]	1000	null
324a3b62-552b-456e-9787-5685bff57798	a6fa7715-afc4-4228-a8cd-9ac1fb535065	x2Integrate	[0, 100, 10000]	1000	null
2a018988-8a93-43b8-9276-36ae715f8229	a6fa7715-afc4-4228-a8cd-9ac1fb535065	x2IntegrateMock	[0, 100, 10000]	1000	null
9a742ae7-dadb-45b6-80b5-90561d450bb7	a6fa7715-afc4-4228-a8cd-9ac1fb535065	fibonacciRecursive	[35]	1000	null
86c0e78c-9cae-4bb7-a440-ad95a381fbfc	a6fa7715-afc4-4228-a8cd-9ac1fb535065	fibonacciRecursiveMock	[35]	1000	null
01b15ee0-eb0e-47be-9e8c-822b7cd732f3	a6fa7715-afc4-4228-a8cd-9ac1fb535065	fibonacciIterative	[35]	1000	null
bd8350d2-7c2d-4b0a-88f2-4dd2a9c11457	a6fa7715-afc4-4228-a8cd-9ac1fb535065	fibonacciIterativeMock	[35]	1000	null
8e327d51-36bf-4218-9c91-5af4030f1ffd	a6fa7715-afc4-4228-a8cd-9ac1fb535065	multiply	[1000000]	1000	null
2f260f4f-e75c-4e06-b8d9-75ce5c4342e8	a6fa7715-afc4-4228-a8cd-9ac1fb535065	multiplyMock	[1000000]	1000	null
69047899-f08c-4628-9d19-2d722dc301f4	a6fa7715-afc4-4228-a8cd-9ac1fb535065	multiplyVector	[10000]	1000	null
14286f69-9c4e-4ed8-aa01-419fd0796126	a6fa7715-afc4-4228-a8cd-9ac1fb535065	multiplyVectorMock	[10000]	1000	null
8e4f367e-5374-45c4-a043-eb0518b6d90d	a6fa7715-afc4-4228-a8cd-9ac1fb535065	factorize	[56]	1000	null
3df7ffdd-9d60-4463-91db-e91e8939a4a5	a6fa7715-afc4-4228-a8cd-9ac1fb535065	factorizeMock	[56]	1000	null
edbf9bcc-0b19-49bb-bfb8-37e967498754	9bdc8bf7-8e77-40c2-8f55-6bf927a2c80b	x2Integrate	[0, 100, 10000]	1000	null
d6870154-49a3-43e5-80b6-bc715989fec2	9bdc8bf7-8e77-40c2-8f55-6bf927a2c80b	x2IntegrateMock	[0, 100, 10000]	1000	null
c23146ec-93f5-44ae-aecc-2a93a43b8a4d	9bdc8bf7-8e77-40c2-8f55-6bf927a2c80b	fibonacciRecursive	[35]	1000	null
44316852-7e3a-42f1-bf5d-abe55e8eb580	9bdc8bf7-8e77-40c2-8f55-6bf927a2c80b	fibonacciRecursiveMock	[35]	1000	null
c30672b8-6b5f-4f5e-816b-4a12aeccefaa	9bdc8bf7-8e77-40c2-8f55-6bf927a2c80b	fibonacciIterative	[35]	1000	null
16809cde-9f68-4616-89a7-752dcb272017	9bdc8bf7-8e77-40c2-8f55-6bf927a2c80b	fibonacciIterativeMock	[35]	1000	null
57206da7-3a01-4ec4-99b2-5ea032719031	9bdc8bf7-8e77-40c2-8f55-6bf927a2c80b	multiply	[1000000]	1000	null
15b8c6d8-3c29-4c53-8a18-1cd5293577fb	9bdc8bf7-8e77-40c2-8f55-6bf927a2c80b	multiplyMock	[1000000]	1000	null
521b196b-1404-4c75-b5f3-051ab4bcace2	9bdc8bf7-8e77-40c2-8f55-6bf927a2c80b	multiplyVector	[10000]	1000	null
7243a8fb-4649-4b6a-b640-5e9c03b4b329	9bdc8bf7-8e77-40c2-8f55-6bf927a2c80b	multiplyVectorMock	[10000]	1000	null
b5cbba5b-3322-46f4-b037-f8bbc01d2acb	9bdc8bf7-8e77-40c2-8f55-6bf927a2c80b	factorize	[56]	1000	null
5b0f5bf3-85e4-4669-9627-28e7d677fdce	9bdc8bf7-8e77-40c2-8f55-6bf927a2c80b	factorizeMock	[56]	1000	null
ed935f54-6066-4a1a-9113-967472d72c76	d02506eb-4d0a-48ac-a07f-415b41e89019	x2Integrate	[0, 100, 10000]	1000	null
ccb72251-93e4-449e-b3a6-b5f60a16f5ca	d02506eb-4d0a-48ac-a07f-415b41e89019	x2IntegrateMock	[0, 100, 10000]	1000	null
3ed03fcf-bd8b-4687-9051-bd4b54dfbaab	d02506eb-4d0a-48ac-a07f-415b41e89019	fibonacciRecursive	[35]	1000	null
7cbb57db-005f-4bda-bbcf-00e87839fafa	d02506eb-4d0a-48ac-a07f-415b41e89019	fibonacciRecursiveMock	[35]	1000	null
6476cae9-9018-46c1-aeac-73ea2ee6a3c4	d02506eb-4d0a-48ac-a07f-415b41e89019	fibonacciIterative	[35]	1000	null
3914164d-4fa5-4784-ad50-604183aadf27	d02506eb-4d0a-48ac-a07f-415b41e89019	fibonacciIterativeMock	[35]	1000	null
01e2b74c-ae85-4a78-bc78-400d11194187	d02506eb-4d0a-48ac-a07f-415b41e89019	multiply	[1000000]	1000	null
116b307c-6036-4214-8785-94e7a7bedd5e	d02506eb-4d0a-48ac-a07f-415b41e89019	multiplyMock	[1000000]	1000	null
27e2b384-2269-4a79-bf6b-985855139fad	d02506eb-4d0a-48ac-a07f-415b41e89019	multiplyVector	[10000]	1000	null
797b2630-821a-4a8a-8427-79117496fcdb	d02506eb-4d0a-48ac-a07f-415b41e89019	multiplyVectorMock	[10000]	1000	null
902b0f41-6bf2-4d42-aff8-4a5b373ac1a8	d02506eb-4d0a-48ac-a07f-415b41e89019	factorize	[56]	1000	null
520e646f-888f-48a2-9357-35ace3a1f865	d02506eb-4d0a-48ac-a07f-415b41e89019	factorizeMock	[56]	1000	null
6081cb38-7bac-4641-82db-c6a0e2d75595	5865d7b3-1952-40c1-bb42-33a3f9dc8afc	x2Integrate	[0, 100, 10000]	1000	null
dc634997-face-4f3d-8036-1241351a0532	5865d7b3-1952-40c1-bb42-33a3f9dc8afc	x2IntegrateMock	[0, 100, 10000]	1000	null
0b0171ba-044a-4bca-8045-8ed662bcb2b1	5865d7b3-1952-40c1-bb42-33a3f9dc8afc	fibonacciRecursive	[35]	1000	null
b0eef869-cb0b-4fdc-98d0-667d8b4911f1	5865d7b3-1952-40c1-bb42-33a3f9dc8afc	fibonacciRecursiveMock	[35]	1000	null
d338831a-23d6-4ea2-a34b-e30a649c7512	5865d7b3-1952-40c1-bb42-33a3f9dc8afc	fibonacciIterative	[35]	1000	null
12010680-8624-4deb-9090-f20059462df9	5865d7b3-1952-40c1-bb42-33a3f9dc8afc	fibonacciIterativeMock	[35]	1000	null
c011f022-e3ad-4144-bd84-7d60c78cff38	5865d7b3-1952-40c1-bb42-33a3f9dc8afc	multiply	[1000000]	1000	null
c607e8aa-d751-4df4-a425-82d2cb86bb6a	5865d7b3-1952-40c1-bb42-33a3f9dc8afc	multiplyMock	[1000000]	1000	null
f177a862-db67-4a50-a2b0-115b2fc58f18	5865d7b3-1952-40c1-bb42-33a3f9dc8afc	multiplyVector	[10000]	1000	null
a025e64d-48be-42f0-889c-6503c2d50b90	5865d7b3-1952-40c1-bb42-33a3f9dc8afc	multiplyVectorMock	[10000]	1000	null
9c00e530-5b2e-487b-ab7f-f2e59f78d89e	5865d7b3-1952-40c1-bb42-33a3f9dc8afc	factorize	[56]	1000	null
cc975bdc-d797-4e53-a36a-379a603d87e4	5865d7b3-1952-40c1-bb42-33a3f9dc8afc	factorizeMock	[56]	1000	null
8fc1a994-8a54-49da-8d95-83094e6350a4	18003139-c722-4abc-b4f4-f2b82d09ef90	x2Integrate	[0, 100, 10000]	1000	null
aa04004b-c48f-44ad-8202-52cb6ff289a8	18003139-c722-4abc-b4f4-f2b82d09ef90	x2IntegrateMock	[0, 100, 10000]	1000	null
36c439f7-5c4a-40bc-a74e-54639d541e16	18003139-c722-4abc-b4f4-f2b82d09ef90	fibonacciRecursive	[35]	1000	null
226d6a76-f692-4415-8445-c653f62c9622	18003139-c722-4abc-b4f4-f2b82d09ef90	fibonacciRecursiveMock	[35]	1000	null
e7c5665b-d1f5-4765-b7a5-b1013465a42d	18003139-c722-4abc-b4f4-f2b82d09ef90	fibonacciIterative	[35]	1000	null
61269776-f5ac-469c-999f-59d1a90e9125	18003139-c722-4abc-b4f4-f2b82d09ef90	fibonacciIterativeMock	[35]	1000	null
f668ed2a-6cf8-4fa7-981f-966169ca7ff1	18003139-c722-4abc-b4f4-f2b82d09ef90	multiply	[1000000]	1000	null
5d781ef6-346e-48c9-a6c0-781cbf862254	18003139-c722-4abc-b4f4-f2b82d09ef90	multiplyMock	[1000000]	1000	null
8e07c902-c2f4-440f-90bd-49844621dbe7	18003139-c722-4abc-b4f4-f2b82d09ef90	multiplyVector	[10000]	1000	null
07bfa258-99b7-4d1a-afbd-19653c2fb6b3	18003139-c722-4abc-b4f4-f2b82d09ef90	multiplyVectorMock	[10000]	1000	null
eb20fd65-92e8-439c-a800-27da832d0a66	18003139-c722-4abc-b4f4-f2b82d09ef90	factorize	[56]	1000	null
148c760b-a496-4e6b-add2-3d7831275ad1	18003139-c722-4abc-b4f4-f2b82d09ef90	factorizeMock	[56]	1000	null
639e6b39-f3f5-4bbb-a166-bfd861fe1b99	72ff8a2c-a0d5-4022-8a24-d9cd4e7cde45	x2Integrate	[0, 100, 10000]	1000	null
09504118-d771-4c38-a5b2-68a1fdecdd65	72ff8a2c-a0d5-4022-8a24-d9cd4e7cde45	x2IntegrateMock	[0, 100, 10000]	1000	null
e171f1ca-080e-4303-bc3d-f8673e0c61a7	72ff8a2c-a0d5-4022-8a24-d9cd4e7cde45	fibonacciRecursive	[35]	1000	null
e3665be8-11ec-492c-8082-b358d7436f7e	72ff8a2c-a0d5-4022-8a24-d9cd4e7cde45	fibonacciRecursiveMock	[35]	1000	null
170d80e7-d1eb-4be6-b5f7-ef98b97bfbf4	72ff8a2c-a0d5-4022-8a24-d9cd4e7cde45	fibonacciIterative	[35]	1000	null
6f0daad2-68fe-4525-a890-4aebdcd21ad6	72ff8a2c-a0d5-4022-8a24-d9cd4e7cde45	fibonacciIterativeMock	[35]	1000	null
496842b8-82ce-4f30-b284-34f9c38d917f	72ff8a2c-a0d5-4022-8a24-d9cd4e7cde45	multiply	[1000000]	1000	null
c9379bf9-f4c4-4ce7-8456-ecc0439e088c	72ff8a2c-a0d5-4022-8a24-d9cd4e7cde45	multiplyMock	[1000000]	1000	null
5b0ff550-6c2f-49f7-8f48-1089ca9301da	72ff8a2c-a0d5-4022-8a24-d9cd4e7cde45	multiplyVector	[10000]	1000	null
ce2d9cbc-a834-4cfc-b233-01fdf7d5a85c	72ff8a2c-a0d5-4022-8a24-d9cd4e7cde45	multiplyVectorMock	[10000]	1000	null
e7e15c58-1517-4937-869b-9e59e529a508	72ff8a2c-a0d5-4022-8a24-d9cd4e7cde45	factorize	[56]	1000	null
238d8aa9-b899-4edc-a131-5a88bd1d5588	72ff8a2c-a0d5-4022-8a24-d9cd4e7cde45	factorizeMock	[56]	1000	null
b6adace7-8860-4532-92d5-7d77c2dbf96b	b69732a4-e070-4381-9655-32883ebdfadd	x2Integrate	[0, 100, 10000]	1000	null
475995f3-b9e9-4198-b2fc-34a5228a3a3d	b69732a4-e070-4381-9655-32883ebdfadd	x2IntegrateMock	[0, 100, 10000]	1000	null
cd854eff-1094-40b2-9517-501979df5a4f	b69732a4-e070-4381-9655-32883ebdfadd	fibonacciRecursive	[35]	1000	null
51dc038c-32b6-4643-a8c5-94ab38d1b81a	b69732a4-e070-4381-9655-32883ebdfadd	fibonacciRecursiveMock	[35]	1000	null
45e77e3e-3049-436c-b5e4-616e4500e1cd	b69732a4-e070-4381-9655-32883ebdfadd	fibonacciIterative	[35]	1000	null
9421225d-feed-4451-b9d0-9a06b34071cc	b69732a4-e070-4381-9655-32883ebdfadd	fibonacciIterativeMock	[35]	1000	null
ba4dcd7c-19bb-4cad-a8da-81e4f92c9b76	b69732a4-e070-4381-9655-32883ebdfadd	multiply	[1000000]	1000	null
b9d40f82-f2dc-4f39-98a2-ca405d2595ba	b69732a4-e070-4381-9655-32883ebdfadd	multiplyMock	[1000000]	1000	null
93ff8948-abcc-40ec-93e3-cff24bb1db74	b69732a4-e070-4381-9655-32883ebdfadd	multiplyVector	[10000]	1000	null
fc4da5ee-fee6-45f5-b41a-cd48a2973f21	b69732a4-e070-4381-9655-32883ebdfadd	multiplyVectorMock	[10000]	1000	null
dbbba38e-906a-4406-ae43-108b04f2c11c	b69732a4-e070-4381-9655-32883ebdfadd	factorize	[56]	1000	null
550e1eba-ced2-4a3e-9100-f35ce77830d0	b69732a4-e070-4381-9655-32883ebdfadd	factorizeMock	[56]	1000	null
\.


--
-- Data for Name: metric; Type: TABLE DATA; Schema: composer; Owner: postgres
--

COPY composer.metric (id, function_result_id, mean, stddev, median, user_time, system, min, max) FROM stdin;
4c527d64-e406-43e4-adab-9b4b3131a475	64645a3c-fd52-4df6-abdf-7beecc668a17	0.029207259469999976	0.0012063769249202505	0.029034497240000002	0.016719014520000006	0.01936222090000001	0.02658148274	0.036985416740000006
a2670ead-6cb6-484f-82e0-2b0f104bb4c0	19b87222-918e-461b-a0e4-1b8525817988	0.02935470104699998	0.0024192768373991883	0.029034703120000005	0.016307013519999986	0.01802901312000001	0.026037016620000002	0.07308457762
2d131fcf-6497-4a94-b573-90170c310122	416c71c5-04e0-4725-94d3-c9d529ec668c	0.10639012809400011	0.0025799150094971135	0.10597613500000001	0.09716891500000018	0.018719754920000012	0.10198165350000002	0.13247458550000002
43710a7f-1db2-4def-8de8-c1eab7072ee1	f932695b-b59a-4d83-83fd-7b1953ee27c7	0.03116311083499996	0.007329046722698101	0.029528123360000002	0.01886491852	0.017479523519999998	0.02568363136	0.14944420736000003
a9409897-cf17-42d6-890b-f4392c2d1a48	3157ed1f-5aca-4a62-845f-99b1232e80ff	0.02989776367799998	0.001796155137037582	0.029664773900000004	0.01727778248	0.01760855243999999	0.025859929400000002	0.049244729400000006
1c852c47-6107-4c01-b991-3a0cd568bb0c	83a4b7f6-6aec-4c1e-a143-dbb256969bed	0.029449184226999988	0.0021838342372419942	0.02911207038	0.01745297352000003	0.017082504720000002	0.02571407938	0.06533020838
a04bd975-db58-4170-b05f-5baeeec5cf39	e27dcf72-9c89-491e-9a62-f4356f4d1c33	0.030371661768000012	0.001184780512854467	0.03022306214	0.019213073679999984	0.01948533099999998	0.027481008140000002	0.03782681314
686526a8-80da-4ff7-9565-3702938432a7	ab08fb6d-0ad8-4856-998a-32fc65d1090f	0.028509755349000014	0.001145467744620013	0.0284475925	0.016673699719999966	0.016650534080000032	0.0252927155	0.036151392500000004
b1f5bf59-8923-438b-98f4-b9a97b993998	76b35674-82b8-494d-b383-118436e19b0c	0.029502381616999963	0.003029674333817694	0.029218107700000003	0.016976483960000007	0.018048506399999997	0.026154935700000003	0.1064255347
26df8dfa-95b3-4eee-9460-a98310f2dd04	a7133fcf-c92b-4dc8-9df7-2a9e788763fd	0.029503899656000056	0.0017134142407997164	0.029342388840000004	0.016489464999999995	0.017785730599999997	0.026185042840000002	0.04995877684
79d5fa84-20af-45b3-8e49-434c91fe1324	88a25685-2c96-44f0-a37e-414f3f69c33d	0.029668615115999997	0.0016040111504591133	0.029500787	0.017010923459999994	0.017692771999999996	0.025773692	0.046259732000000005
8560dbb4-9f12-4975-b1b1-11669ae3cfaf	a5472a47-5470-4447-9f40-3694bdf40a32	0.02967104025200005	0.001447385662482104	0.029556367120000003	0.016943429879999982	0.01778159878000001	0.025924546120000003	0.03759191212
038fdc87-c174-4cd9-b3fd-600041fc3bc3	324a3b62-552b-456e-9787-5685bff57798	0.1398554765420001	0.0035887158184773164	0.13920337186	0.12311236599999992	0.03689579899999994	0.13774484886000002	0.16957611886
6e83fba2-81dd-461c-b230-338c87070152	2a018988-8a93-43b8-9276-36ae715f8229	0.13592766150999988	0.004439849191418966	0.13507335822	0.11260711500000001	0.03652365799999997	0.13328714472	0.16960995172
6dccc561-4cef-484b-8a60-2a0bb1e57c6a	9a742ae7-dadb-45b6-80b5-90561d450bb7	1.0528682867849997	0.025916446145189405	1.05760179474	1.0399381239999999	0.039212112999999986	1.02506202874	1.12763880274
86fdcab4-1b2c-4492-a47a-37cb68616097	86c0e78c-9cae-4bb7-a440-ad95a381fbfc	0.13635719157000029	0.004556244197583257	0.1352213543	0.11240158799999991	0.037149857000000015	0.13361026580000002	0.1673637078
a9002f68-fe78-457e-a19d-0b39254f398d	01b15ee0-eb0e-47be-9e8c-822b7cd732f3	0.1361039866629999	0.004126900076559657	0.13534517748000002	0.1130180780000001	0.03629289599999999	0.13381677298000003	0.17033261798000002
9e0f393f-1203-47cf-9ed2-973b1ca04eb6	bd8350d2-7c2d-4b0a-88f2-4dd2a9c11457	0.1359710627320001	0.004039234632809856	0.13521342446	0.11262496900000009	0.036567227000000035	0.13376903346000002	0.17075294246
4e323e65-c2f4-4c12-affc-bb73ba90bbce	8e327d51-36bf-4218-9c91-5af4030f1ffd	0.1442619576130001	0.00391694163743002	0.14355936466000002	0.12784353799999984	0.03872002499999999	0.14213108766000002	0.17564500666000002
f621c6ca-63a4-4c3d-a93f-186befe16fea	2f260f4f-e75c-4e06-b8d9-75ce5c4342e8	0.1360153095269999	0.004116513316031009	0.13521041278	0.11277714399999994	0.036491691999999985	0.13375294027999998	0.16593573028
6bf5e2f0-0ba0-4195-bf4a-77dba1dc8ac8	69047899-f08c-4628-9d19-2d722dc301f4	0.14016044684000012	0.003961610192812905	0.139412996	0.11901451400000002	0.03796196400000002	0.137909847	0.170238097
af36d6e8-3a99-4ca2-8cca-31c2fdeb4faa	14286f69-9c4e-4ed8-aa01-419fd0796126	0.13591037956100013	0.003672411555999857	0.13518003104000004	0.11257966200000019	0.03656073800000001	0.13348154504	0.16394110204
d37918a1-859f-418c-8495-330125998b1b	8e4f367e-5374-45c4-a043-eb0518b6d90d	0.13606616088200008	0.0036526938951649676	0.13540663028	0.11294422200000016	0.036336872000000006	0.13387215628000002	0.16546191328
8680e41d-c3da-45d3-a7ee-3739eecc1bab	3df7ffdd-9d60-4463-91db-e91e8939a4a5	0.13596194828200003	0.003997753360667843	0.1352492957	0.11247702700000001	0.036692163000000035	0.1336212507	0.1680212097
b6b291ee-6350-4bd7-bc13-f48732c0d05b	edbf9bcc-0b19-49bb-bfb8-37e967498754	0.03714727722400007	0.0038537202916483103	0.03621965244	0.02395094802000003	0.024524896999999983	0.03241073844	0.11504417844
2ad718ff-cb21-4cdb-b281-c9070eba23eb	d6870154-49a3-43e5-80b6-bc715989fec2	0.03652020717300005	0.0034541409222849065	0.035576511719999995	0.024258508640000023	0.02367671278000004	0.03237112522	0.08989814722000002
7fbeec48-c785-44e6-865c-00c90895b514	c23146ec-93f5-44ae-aecc-2a93a43b8a4d	0.14793502860599994	0.00802991516564117	0.14675774248	0.13901339199999976	0.029996174159999996	0.13918697348	0.33267724548000005
8b6f7fe5-478e-4603-8fdc-9bcdec971c91	44316852-7e3a-42f1-bf5d-abe55e8eb580	0.03811340182	0.002544091198872624	0.037459644300000004	0.024923137999999994	0.024853768519999996	0.033827621300000006	0.05871417030000001
4be4f8dc-cbad-4496-8ae4-ec9c04efe9d2	c30672b8-6b5f-4f5e-816b-4a12aeccefaa	0.03781298474499991	0.002591491782823671	0.03718086738	0.02369270925999994	0.02578862155999999	0.03370394438	0.05910713438
0f735efb-a8f6-4325-b14d-951e19b0f221	16809cde-9f68-4616-89a7-752dcb272017	0.037649373919999965	0.0024042421710216188	0.037038169680000003	0.024419992959999987	0.024840328479999993	0.033393784180000004	0.053010591180000004
9dc12cb4-d766-41c7-af22-da80d19e1365	57206da7-3a01-4ec4-99b2-5ea032719031	0.041622922183000004	0.0026349952878032385	0.040938943059999996	0.027395978319999983	0.026650495419999993	0.03729048106	0.06529500406000001
24d97e84-a0fa-4fbb-abd1-0b25f8e2f65f	15b8c6d8-3c29-4c53-8a18-1cd5293577fb	0.03820417232000001	0.005220109132822453	0.037298736799999996	0.024447664200000017	0.02570879257999998	0.033762208800000006	0.1614654048
0c740ea3-f278-4b6c-969b-4bd1fb229ec6	521b196b-1404-4c75-b5f3-051ab4bcace2	0.039326270041000005	0.003617845286315784	0.03849776532	0.024789858979999996	0.02656518916000001	0.03417243232	0.10459012732
69a7a386-b09a-44f0-966a-b3a7b08984c9	7243a8fb-4649-4b6a-b640-5e9c03b4b329	0.03761175965699999	0.002396182889691267	0.037045865659999996	0.025085505100000037	0.024293097840000007	0.03299065066	0.05151484566
ae4aa3b3-e519-4c67-b1ce-5fc664db54ec	b5cbba5b-3322-46f4-b037-f8bbc01d2acb	0.03767401333199999	0.00242729482844709	0.037063425080000004	0.024508579919999957	0.024968676879999988	0.03354529608	0.053668715080000005
fa44570e-0704-4f97-9218-edc67c94f6a5	5b0f5bf3-85e4-4669-9627-28e7d677fdce	0.03764390267200001	0.002609197363775081	0.037040843860000006	0.025286479379999983	0.02399839140000001	0.03313851986000001	0.06019866086
a6742271-dcdb-4b27-b287-693a2fa326b3	ed935f54-6066-4a1a-9113-967472d72c76	0.17494761727500008	0.0046061613452980455	0.1740279301	0.15556070699999996	0.048600061	0.17209294360000002	0.2049842866
c1d35634-1741-4921-94bf-4c7f08e8579a	ccb72251-93e4-449e-b3a6-b5f60a16f5ca	0.17166467292000037	0.004730688520173889	0.17072749114000002	0.14963159099999973	0.04871071100000009	0.16867808564	0.20217744964
d36e99a7-6789-4efe-933b-ad75462c0cd0	3ed03fcf-bd8b-4687-9051-bd4b54dfbaab	0.8616657371759997	0.008589794392131337	0.85888697386	0.8485090039999996	0.050576527999999926	0.85684785886	0.8954887118600001
58dfe7e8-591c-4cf7-a7c9-062a9b3c8ee0	7cbb57db-005f-4bda-bbcf-00e87839fafa	0.17167162904100006	0.004585007050371829	0.17071916108000001	0.1502247549999998	0.04815494600000001	0.16885856508	0.20311509508
89b9e30b-a5c0-4d8e-bb37-6bf0bff8c809	6476cae9-9018-46c1-aeac-73ea2ee6a3c4	0.1717878393469999	0.004610743278543807	0.17075983032	0.15035964399999993	0.04814981399999997	0.16912410632	0.20096246032
52f69658-1376-4866-b3c5-47029e7ea8f2	3914164d-4fa5-4784-ad50-604183aadf27	0.17176964733400005	0.004558199458348077	0.17084659870000002	0.14991179300000015	0.04855419099999992	0.16863299620000002	0.20151262520000002
772f8f84-5dd6-4945-a3bc-30cbae35269f	01e2b74c-ae85-4a78-bc78-400d11194187	0.19175972769199984	0.005206647792356486	0.19063349274000002	0.17271870199999992	0.05038876999999998	0.18895226624	0.22215957624000002
8312b6af-d29b-4af0-932b-13d32c13907d	116b307c-6036-4214-8785-94e7a7bedd5e	0.17172193376199987	0.004696457748751807	0.17081972520000002	0.14955601899999998	0.04881282499999996	0.16880175270000003	0.20222034270000003
f0a48903-59d3-4ae3-8530-b9ed85bb3ce7	27e2b384-2269-4a79-bf6b-985855139fad	0.17641313479000015	0.0044588373311406455	0.17557088494	0.1557467660000001	0.04855636800000007	0.17343643244000004	0.20641220844000002
43da243a-6c50-4b05-9b15-648cbcb45c0f	797b2630-821a-4a8a-8427-79117496fcdb	0.17170487930400016	0.004641876111533619	0.17084104836	0.14982689100000002	0.04860426799999999	0.16902422986	0.20612776086
3aa40115-6344-4ad1-97a8-8f09733828d0	902b0f41-6bf2-4d42-aff8-4a5b373ac1a8	0.17180620954700004	0.004921084306770653	0.17080408438000003	0.14988081300000006	0.04864966000000003	0.16896316788000001	0.20730699588
0d04f37d-b943-444f-9082-dc6692fbe2d0	520e646f-888f-48a2-9357-35ace3a1f865	0.17239287506199988	0.014475608939918236	0.16993394636000003	0.1507140510000002	0.04866203800000004	0.16817480836	0.31209053836
d1ade698-df81-4adb-9607-8c22d78e0599	6081cb38-7bac-4641-82db-c6a0e2d75595	0.07725081867600012	0.004577741654500598	0.07651706327999999	0.090348549	0.08410523400000004	0.07017491478	0.16948315278
7e263799-6d36-421d-9175-0b11e9e6c463	dc634997-face-4f3d-8036-1241351a0532	0.07970069478199998	0.007076795217518216	0.07857731904	0.095544175	0.0837122199999999	0.07212971204	0.23852047504
06bc88d5-dc41-4afb-bd14-3fcf43ed864c	0b0171ba-044a-4bca-8045-8ed662bcb2b1	0.24362573154900014	0.006845253510957448	0.24174143818000002	0.26803747299999986	0.08211759699999996	0.23429118018	0.30227808718000004
5c6e746b-9193-4d40-9153-82f67b9f71a2	b0eef869-cb0b-4fdc-98d0-667d8b4911f1	0.08043675731000012	0.0029411040216356242	0.07993420434000001	0.09800832300000005	0.08353863299999993	0.07414276584000001	0.11024088684
93c59e04-8970-4d06-929f-59131f57c066	d338831a-23d6-4ea2-a34b-e30a649c7512	0.08096209115899991	0.003707048993127933	0.08023569638	0.09941956800000006	0.08323571700000015	0.07512265388	0.12571596788
46f4fc55-558c-47fc-b04b-de03e114db56	12010680-8624-4deb-9090-f20059462df9	0.08105506485399994	0.0038295786200146513	0.08059664622000001	0.09955541900000008	0.08332654500000002	0.07509707972	0.15229713072
30b5cf97-1a05-4faa-8aba-3cef19460841	c011f022-e3ad-4144-bd84-7d60c78cff38	0.08323743507100002	0.0034888157167983717	0.08263478954	0.1026647580000001	0.08538580099999991	0.07767915704	0.12750531204
cd3c6d5c-4a37-410e-96cc-c6fa9c8411fb	c607e8aa-d751-4df4-a425-82d2cb86bb6a	0.07913518954300006	0.003157783909093645	0.07891066348	0.09298704699999985	0.08567257800000004	0.07247673698000001	0.11897489498000001
c4bc5938-1289-466a-961c-a56341e3381e	f177a862-db67-4a50-a2b0-115b2fc58f18	0.08174802744100003	0.00757336116340774	0.07976126868	0.0955047870000001	0.08741941899999997	0.07381520718	0.20217589118
47185122-0118-43a9-bd45-958b4cfe351e	a025e64d-48be-42f0-889c-6503c2d50b90	0.08100415084099996	0.0043745995951528895	0.08030776896	0.09872173400000006	0.08394646900000009	0.07325381496	0.16030131396
34d9b944-dcf6-4cf6-886c-fe3c474cbe5b	9c00e530-5b2e-487b-ab7f-f2e59f78d89e	0.08056287110200007	0.0027559225761302464	0.080093952	0.09879190300000007	0.08303356899999997	0.074820456	0.117397732
da883971-61f0-4b6e-bbd6-08586c79a38a	cc975bdc-d797-4e53-a36a-379a603d87e4	0.08103647170199994	0.0037052332926424333	0.0802934069	0.09817476500000004	0.08484216100000005	0.0753082999	0.12080943090000001
b333227e-0519-49bf-b5a8-243707e63199	8fc1a994-8a54-49da-8d95-83094e6350a4	0.3860726873000003	0.021338020541928293	0.38857782778000005	0.6460757059999985	0.13295806599999993	0.35533375627999997	0.6558626642800001
0f2306e3-0bd6-4c3e-9552-98ddd154b332	aa04004b-c48f-44ad-8202-52cb6ff289a8	0.39722281192699976	0.005290469641681883	0.39717504056	0.6576949319999995	0.13680237700000006	0.38169468856000005	0.42138484156000006
a15ab578-ae73-49bb-a6ca-525fd32cbc3a	36c439f7-5c4a-40bc-a74e-54639d541e16	2.1177988495980027	0.00937122708232194	2.11400074298	2.367198597999997	0.129072556	2.10873144748	2.16555824748
e9c3ce60-5a5c-4e37-ad99-773cef77d3f6	226d6a76-f692-4415-8445-c653f62c9622	0.39284202560599984	0.01214319731633698	0.3958820578800001	0.6493329209999997	0.13654585800000005	0.35429341038000006	0.42142472638000006
3e26931e-02d3-452d-93a2-c3afd2cbb22c	e7c5665b-d1f5-4765-b7a5-b1013465a42d	0.4039069999489995	0.017026154401865175	0.4028281788	0.6673854910000002	0.13918468300000014	0.3566909968	0.6388817708000001
e87ba688-6dfa-49e2-ad85-3f55c0c1a193	61269776-f5ac-469c-999f-59d1a90e9125	0.40713544284299985	0.005354570609815562	0.40680128696000006	0.6724030419999998	0.1399316760000001	0.39186035846	0.42280147446000005
0b9bab8c-b965-44b3-a2eb-30e7e3a6e419	f668ed2a-6cf8-4fa7-981f-966169ca7ff1	0.5363713010719999	0.10927737710923939	0.45004424746000005	0.8893181129999992	0.17761059799999993	0.38084744946000004	0.6926816304600001
e2d91233-f055-4739-8748-5b7f922d7a79	5d781ef6-346e-48c9-a6c0-781cbf862254	0.5834997247060001	0.07424902545997591	0.61857925802	0.9580546070000008	0.20090573599999995	0.35148648402	0.65018868102
9a0ba998-1d9d-4009-8017-a44e22670fa4	8e07c902-c2f4-440f-90bd-49844621dbe7	0.5767520446900004	0.09065636615967133	0.6263675301	0.9547738119999991	0.19999342299999995	0.3572942406	0.6431490166
f67663e6-3c19-486f-a134-65736ae344a8	07bfa258-99b7-4d1a-afbd-19653c2fb6b3	0.5625896133580008	0.09696042601801705	0.6183457079200001	0.9244266800000006	0.1934726880000002	0.35169447442	0.63839548542
f66194b1-59a1-44dc-99b9-6f1aefec3343	eb20fd65-92e8-439c-a800-27da832d0a66	0.5622190921489995	0.09828044057601931	0.6186239835	0.9241084749999979	0.19323895699999988	0.352249614	0.6867090480000001
14642ee5-7881-4690-baf1-d91ebcff5282	148c760b-a496-4e6b-add2-3d7831275ad1	0.5610667798639999	0.09813574546785703	0.6179826051800001	0.9200379160000007	0.195065548	0.35122688868	0.6770328266800001
6e19cfb5-009d-4c8d-9414-53048c4511d5	639e6b39-f3f5-4bbb-a166-bfd861fe1b99	0.031126629311000003	0.002422982963046422	0.030861302120000003	0.017954515180000016	0.0222869798	0.02796328962	0.09161354562000001
7150fb83-253b-4274-9d84-bf0cc609ac79	09504118-d771-4c38-a5b2-68a1fdecdd65	0.03213149006300001	0.005733771766744344	0.03159353294	0.018932361760000002	0.022827652379999997	0.027977543440000002	0.19188756044000002
b0c63001-a2d3-4381-bab3-c8ffee0260cd	e171f1ca-080e-4303-bc3d-f8673e0c61a7	0.05427886148299998	0.004129502950802804	0.05332819626	0.040862465000000014	0.025014922400000007	0.049143964760000006	0.09088161276000001
13a05dd7-0aa4-4e99-b43b-dff177ebe083	e3665be8-11ec-492c-8082-b358d7436f7e	0.032134121171999985	0.0015445605270179342	0.031927803840000005	0.019119673199999996	0.022196649079999988	0.02903526284	0.051364559840000004
902fdd0e-a06e-40f8-87ad-ee60fb410058	170d80e7-d1eb-4be6-b5f7-ef98b97bfbf4	0.032338049960999996	0.0014091039527139642	0.032199613780000004	0.01846636630000002	0.02302780643999999	0.028844359780000002	0.039848958780000006
a512ec16-07aa-4108-9645-63b564019fe2	6f0daad2-68fe-4525-a890-4aebdcd21ad6	0.03230894527399998	0.0018046693223397014	0.03211737326	0.018715526919999986	0.022774474240000012	0.029068576260000003	0.06144223026
685c301e-3182-4461-bd24-75dc135ee905	496842b8-82ce-4f30-b284-34f9c38d917f	0.03392093226700007	0.0015876306575137332	0.03378274482	0.02060700168000004	0.023092771160000015	0.030539921820000003	0.05071896682
9df66771-6e59-436f-ae2a-134c4228cb84	c9379bf9-f4c4-4ce7-8456-ecc0439e088c	0.03230957057999999	0.001918767597672419	0.03202331828	0.018873159599999995	0.022784436560000013	0.02920655528	0.05577927628
156b59f4-5f25-44ef-be6c-648fc0fd68ae	5b0ff550-6c2f-49f7-8f48-1089ca9301da	0.03330534314600005	0.002405223634587773	0.033037543360000005	0.019315635800000015	0.023243211760000007	0.029696585360000003	0.08915533336
ebdc35aa-150f-491e-b3c8-904c0f0c98b2	ce2d9cbc-a834-4cfc-b233-01fdf7d5a85c	0.03202782564700003	0.0022914635044362485	0.03175242392	0.018208791959999988	0.02310057164000004	0.02879918592	0.06854470092
26ed973e-acf0-4cc6-9498-52b997ce09eb	e7e15c58-1517-4937-869b-9e59e529a508	0.03217651201699996	0.0015518792826782022	0.031999956659999995	0.01855710680000001	0.02281124913999997	0.02935072316	0.047900884160000004
4ddb1545-029a-43ea-a5d8-ea205ac8855e	238d8aa9-b899-4edc-a131-5a88bd1d5588	0.03204335861500003	0.0013031186326837997	0.031931998280000005	0.019381046320000043	0.021903841879999983	0.02905242978	0.03783814378000001
098a0c51-7ab3-4169-b528-0af80dea31fc	b6adace7-8860-4532-92d5-7d77c2dbf96b	0.14186561948200005	0.00454308106800912	0.14096104040000002	0.11756468699999983	0.04189891899999994	0.13948388139999998	0.1719679724
949eddcb-5209-4657-aee7-b1f60f003af3	475995f3-b9e9-4198-b2fc-34a5228a3a3d	0.1415378697659999	0.00435043365414894	0.1406919702	0.11665275499999998	0.041734107000000034	0.1391548392	0.17859463920000002
5e06384a-6c78-4420-a6d1-0cd057c7e8b7	cd854eff-1094-40b2-9517-501979df5a4f	0.41189253862099995	0.008052831261672761	0.40948933476000005	0.39062892199999977	0.04462335099999995	0.40780462376000004	0.44917108176000003
5ba5d240-36f0-4f13-b0cf-67b162918faf	51dc038c-32b6-4643-a8c5-94ab38d1b81a	0.14139182194899994	0.003957231635857542	0.140641516	0.1164094939999999	0.041866644999999994	0.139184958	0.171811126
0d2b2dc2-dbb1-4a82-8fdd-a80f94670451	45e77e3e-3049-436c-b5e4-616e4500e1cd	0.14151991544800005	0.004168299555119829	0.1407501859	0.11634152800000003	0.04203470400000002	0.1389465184	0.1765450184
4167b063-4924-4b65-a0d6-e78be5f069a2	9421225d-feed-4451-b9d0-9a06b34071cc	0.14149746310300004	0.004070910483021819	0.14073317302	0.11603593200000008	0.042295265000000026	0.13910203052	0.17173340752000002
aa289a90-377e-4b4a-8af5-aff415245204	ba4dcd7c-19bb-4cad-a8da-81e4f92c9b76	0.14917619768099996	0.004468285178670759	0.1482605233	0.125537482	0.04415236199999999	0.1463972228	0.1820307998
2ad1868b-eb9d-4679-9706-08de7219527a	b9d40f82-f2dc-4f39-98a2-ca405d2595ba	0.14150146544199996	0.0040373994380114445	0.14075485597999998	0.11630298500000016	0.04205416399999993	0.13917682898	0.17600742998
8f2de99f-6bb6-459b-9738-9f100fe78ab7	93ff8948-abcc-40ec-93e3-cff24bb1db74	0.144881210291	0.004176414600182868	0.14406599620000002	0.11971061000000004	0.04335897800000001	0.14273544470000002	0.17794694670000002
645f0b45-cbc7-4bbe-9751-b6df0ca3e2c7	fc4da5ee-fee6-45f5-b41a-cd48a2973f21	0.14154989551499994	0.00429696686190094	0.1407308126	0.11631278299999984	0.04208787100000003	0.13883962260000002	0.17252746260000001
656e42a9-4211-4e6c-bfa4-d0d285ca6b55	dbbba38e-906a-4406-ae43-108b04f2c11c	0.14153999276699983	0.004070995191307745	0.14074752906	0.11630693999999976	0.04206791300000004	0.13915698356	0.17380611956
894421fb-115f-4633-bd37-44b71d2cbde6	550e1eba-ced2-4a3e-9100-f35ce77830d0	0.141591751701	0.004293063195979355	0.14069717312000002	0.11609549300000005	0.04236084400000001	0.13898902662	0.17314546062000002
\.


--
-- Name: design design_name_key; Type: CONSTRAINT; Schema: composer; Owner: postgres
--

ALTER TABLE ONLY composer.design
    ADD CONSTRAINT design_name_key UNIQUE (name);


--
-- Name: design design_pkey; Type: CONSTRAINT; Schema: composer; Owner: postgres
--

ALTER TABLE ONLY composer.design
    ADD CONSTRAINT design_pkey PRIMARY KEY (id);


--
-- Name: experiment experiment_pkey; Type: CONSTRAINT; Schema: composer; Owner: postgres
--

ALTER TABLE ONLY composer.experiment
    ADD CONSTRAINT experiment_pkey PRIMARY KEY (id);


--
-- Name: function_result function_result_pkey; Type: CONSTRAINT; Schema: composer; Owner: postgres
--

ALTER TABLE ONLY composer.function_result
    ADD CONSTRAINT function_result_pkey PRIMARY KEY (id);


--
-- Name: metric metric_pkey; Type: CONSTRAINT; Schema: composer; Owner: postgres
--

ALTER TABLE ONLY composer.metric
    ADD CONSTRAINT metric_pkey PRIMARY KEY (id);


--
-- Name: experiment experiment_design_id_fkey; Type: FK CONSTRAINT; Schema: composer; Owner: postgres
--

ALTER TABLE ONLY composer.experiment
    ADD CONSTRAINT experiment_design_id_fkey FOREIGN KEY (design_id) REFERENCES composer.design(id) ON DELETE CASCADE;


--
-- Name: function_result function_result_experiment_id_fkey; Type: FK CONSTRAINT; Schema: composer; Owner: postgres
--

ALTER TABLE ONLY composer.function_result
    ADD CONSTRAINT function_result_experiment_id_fkey FOREIGN KEY (experiment_id) REFERENCES composer.experiment(id) ON DELETE CASCADE;


--
-- Name: metric metric_function_result_id_fkey; Type: FK CONSTRAINT; Schema: composer; Owner: postgres
--

ALTER TABLE ONLY composer.metric
    ADD CONSTRAINT metric_function_result_id_fkey FOREIGN KEY (function_result_id) REFERENCES composer.function_result(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

