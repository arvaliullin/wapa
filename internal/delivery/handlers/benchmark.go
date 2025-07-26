package handlers

import (
	"database/sql"
	"errors"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/arvaliullin/wapa/internal/delivery"
	"github.com/arvaliullin/wapa/internal/domain"
	"github.com/arvaliullin/wapa/internal/persistence"
	"github.com/labstack/echo/v4"
)

type BenchmarkHandler struct {
	BenchmarkRepo persistence.BenchmarkRepositoryContract
}

func RegisterBenchmarkHandler(
	httpService delivery.HttpService,
	repo persistence.BenchmarkRepositoryContract,
) {
	handler := &BenchmarkHandler{
		BenchmarkRepo: repo,
	}
	e := httpService.(*delivery.EchoHttpService).Echo

	e.GET("/api/benchmark", handler.GetBenchmarkResults)
	e.GET("/api/benchmark-diff", handler.GetBenchmarkDiff)
	e.GET("/api/benchmark/all", handler.GetAllBenchmarkResults)
	e.GET("/api/benchmark-diff/all", handler.GetAllBenchmarkDiffs)
	e.GET("/api/benchmark/mock", handler.GetBenchmarksOnlyMock)
	e.GET("/api/benchmark/not-mock", handler.GetBenchmarksOnlyNotMock)
	e.GET("/api/benchmark/reliable", handler.GetReliableBenchmarks)
}

// GetBenchmarkResults возвращает результаты бенчмарков для заданной метрики и архитектуры.
//
// @Summary      Получить результаты бенчмарков
// @Description  Получает результаты по заданной метрике (mean, median, stddev, min, max) и архитектуре (например, arm64, amd64)
// @Tags         Benchmark
// @Produce      json
// @Param        metric  query     string true  "Метрика (mean, median, stddev, min, max)"
// @Param        arch    query     string true  "Архитектура (arm64, amd64)"
// @Success      200     {object}  domain.BenchmarkResults
// @Failure      400     {object}  object "Ошибка в запросе из-за отсутствия параметров"
// @Failure      404     {object}  object "Данные не найдены"
// @Failure      500     {object}  object "Внутренняя ошибка сервера"
// @Router       /api/benchmark [get].
func (h *BenchmarkHandler) GetBenchmarkResults(c echo.Context) error {
	metric := c.QueryParam("metric")
	arch := c.QueryParam("arch")

	if metric == "" || arch == "" {
		return c.JSON(http.StatusBadRequest, struct {
			Error string `json:"error"`
		}{Error: "parameters 'metric' and 'arch' are required"})
	}

	results, err := h.BenchmarkRepo.GetBenchmarkResults(metric, arch)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(http.StatusNotFound, struct {
				Error string `json:"error"`
			}{Error: "results not found"})
		}
		return c.JSON(http.StatusInternalServerError, struct {
			Error string `json:"error"`
		}{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, results)
}

// GetBenchmarkDiff возвращает разницу между обычной функцией и её Mock-аналогом.
//
// @Summary      Получить разницу между функцией и Mock
// @Description  Возвращает разницу по каждому языку между функцией и её Mock-версией (например, d_factorize = factorize - factorizeMock)
// @Tags         Benchmark
// @Produce      json
// @Param        metric  query     string true  "Метрика (mean, median, stddev, min, max)"
// @Param        arch    query     string true  "Архитектура (arm64, amd64)"
// @Success      200     {object}  domain.BenchmarkResults
// @Failure      400     {object}  object "Ошибка в запросе из-за отсутствия параметров"
// @Failure      404     {object}  object "Данные не найдены"
// @Failure      500     {object}  object "Внутренняя ошибка сервера"
// @Router       /api/benchmark-diff [get].
func (h *BenchmarkHandler) GetBenchmarkDiff(c echo.Context) error {
	metric := c.QueryParam("metric")
	arch := c.QueryParam("arch")

	if metric == "" || arch == "" {
		return c.JSON(http.StatusBadRequest, struct {
			Error string `json:"error"`
		}{Error: "parameters 'metric' and 'arch' are required"})
	}

	results, err := h.BenchmarkRepo.GetBenchmarkResults(metric, arch)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(http.StatusNotFound, struct {
				Error string `json:"error"`
			}{Error: "results not found"})
		}
		return c.JSON(http.StatusInternalServerError, struct {
			Error string `json:"error"`
		}{Error: err.Error()})
	}

	realCases := make(map[string]domain.BenchmarkCase)
	mockCases := make(map[string]domain.BenchmarkCase)

	for _, item := range results.Results {
		if strings.HasSuffix(item.Name, "Mock") {
			mockCases[strings.TrimSuffix(item.Name, "Mock")] = item
		} else {
			realCases[item.Name] = item
		}
	}

	var diffResults []domain.BenchmarkCase
	for name, real := range realCases {
		mock, ok := mockCases[name]
		if !ok {
			continue
		}
		diffResults = append(diffResults, domain.BenchmarkCase{
			Name:       "d_" + name,
			Go:         math.Abs(real.Go - mock.Go),
			Cpp:        math.Abs(real.Cpp - mock.Cpp),
			Rust:       math.Abs(real.Rust - mock.Rust),
			Javascript: math.Abs(real.Javascript - mock.Javascript),
		})
	}

	diff := domain.BenchmarkResults{
		Arch:    results.Arch,
		Metric:  results.Metric,
		Results: diffResults,
	}

	return c.JSON(http.StatusOK, diff)
}

// GetAllBenchmarkResults возвращает массив всех BenchmarkResults по всем архитектурам и метрикам.
//
// @Summary      Получить результаты бенчмарков по всем архитектурам и метрикам
// @Description  Возвращает массив всех результатов бенчмарков по всем архитектурам (например, amd64, arm64) и всем метрикам (mean, median, stddev, min, max).
// @Tags         Benchmark
// @Produce      json
// @Success      200     {array}   domain.BenchmarkResults
// @Failure      500     {object}  object "Внутренняя ошибка сервера"
// @Router       /api/benchmark/all [get].
func (h *BenchmarkHandler) GetAllBenchmarkResults(c echo.Context) error {
	results, err := h.BenchmarkRepo.GetAllBenchmarkResults()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, struct {
			Error string `json:"error"`
		}{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, results)
}

// GetAllBenchmarkDiffs возвращает все дифференциалы по всем архитектурам и метрикам
//
// @Summary      Получить разницу между функцией и Mock по всем архитектурам и метрикам
// @Description  Возвращает массив разниц по каждому языку между функцией и её Mock-версией для всех архитектур и метрик (например, d_factorize = factorize - factorizeMock).
// @Tags         Benchmark
// @Produce      json
// @Success      200     {array}  domain.BenchmarkResults
// @Failure      500     {object} object "Внутренняя ошибка сервера"
// @Router       /api/benchmark-diff/all [get].
func (h *BenchmarkHandler) GetAllBenchmarkDiffs(c echo.Context) error {
	allResults, err := h.BenchmarkRepo.GetAllBenchmarkResults()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, struct {
			Error string `json:"error"`
		}{Error: err.Error()})
	}

	var allDiffs []domain.BenchmarkResults
	for _, results := range allResults {
		realCases := make(map[string]domain.BenchmarkCase)
		mockCases := make(map[string]domain.BenchmarkCase)
		for _, item := range results.Results {
			if strings.HasSuffix(item.Name, "Mock") {
				mockCases[strings.TrimSuffix(item.Name, "Mock")] = item
			} else {
				realCases[item.Name] = item
			}
		}

		var diffResults []domain.BenchmarkCase
		for name, real := range realCases {
			mock, ok := mockCases[name]
			if !ok {
				continue
			}
			diffResults = append(diffResults, domain.BenchmarkCase{
				Name:       "d_" + name,
				Go:         math.Abs(real.Go - mock.Go),
				Cpp:        math.Abs(real.Cpp - mock.Cpp),
				Rust:       math.Abs(real.Rust - mock.Rust),
				Javascript: math.Abs(real.Javascript - mock.Javascript),
			})
		}
		if len(diffResults) > 0 {
			allDiffs = append(allDiffs, domain.BenchmarkResults{
				Arch:    results.Arch,
				Metric:  results.Metric,
				Results: diffResults,
			})
		}
	}
	return c.JSON(http.StatusOK, allDiffs)
}

// GetBenchmarksOnlyMock возвращает только те кейсы, имя которых заканчивается на "Mock".
//
// @Summary      Получить результаты только Mock-функций
// @Description  Возвращает только бенчмарки с постфиксом Mock в имени функции.
// @Tags         Benchmark
// @Produce      json
// @Param        metric  query     string true  "Метрика (mean, median, stddev, min, max)"
// @Param        arch    query     string true  "Архитектура (arm64, amd64)"
// @Success      200     {object}  domain.BenchmarkResults
// @Failure      400     {object}  object "Ошибка в запросе из-за отсутствия параметров"
// @Failure      404     {object}  object "Данные не найдены"
// @Failure      500     {object}  object "Внутренняя ошибка сервера"
// @Router       /api/benchmark/mock [get].
func (h *BenchmarkHandler) GetBenchmarksOnlyMock(c echo.Context) error {
	metric := c.QueryParam("metric")
	arch := c.QueryParam("arch")

	if metric == "" || arch == "" {
		return c.JSON(http.StatusBadRequest, struct {
			Error string `json:"error"`
		}{Error: "parameters 'metric' and 'arch' are required"})
	}

	results, err := h.BenchmarkRepo.GetBenchmarksOnlyMock(metric, arch)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(http.StatusNotFound, struct {
				Error string `json:"error"`
			}{Error: "results not found"})
		}
		return c.JSON(http.StatusInternalServerError, struct {
			Error string `json:"error"`
		}{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, results)
}

// GetBenchmarksOnlyNotMock возвращает только те кейсы, имя которых НЕ заканчивается на "Mock".
//
// @Summary      Получить результаты только "реальных" функций (без Mock)
// @Description  Возвращает только бенчмарки без постфикса Mock в имени функции.
// @Tags         Benchmark
// @Produce      json
// @Param        metric  query     string true  "Метрика (mean, median, stddev, min, max)"
// @Param        arch    query     string true  "Архитектура (arm64, amd64)"
// @Success      200     {object}  domain.BenchmarkResults
// @Failure      400     {object}  object "Ошибка в запросе из-за отсутствия параметров"
// @Failure      404     {object}  object "Данные не найдены"
// @Failure      500     {object}  object "Внутренняя ошибка сервера"
// @Router       /api/benchmark/not-mock [get].
func (h *BenchmarkHandler) GetBenchmarksOnlyNotMock(c echo.Context) error {
	metric := c.QueryParam("metric")
	arch := c.QueryParam("arch")

	if metric == "" || arch == "" {
		return c.JSON(http.StatusBadRequest, struct {
			Error string `json:"error"`
		}{Error: "parameters 'metric' and 'arch' are required"})
	}

	results, err := h.BenchmarkRepo.GetBenchmarksOnlyNotMock(metric, arch)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(http.StatusNotFound, struct {
				Error string `json:"error"`
			}{Error: "results not found"})
		}
		return c.JSON(http.StatusInternalServerError, struct {
			Error string `json:"error"`
		}{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, results)
}

// GetReliableBenchmarks выбирает надёжные функции для сравнения по коэффициенту вариации (Cv).
//
// @Summary     Получить надёжные функции по коэффициенту вариации
// @Description Возвращает список имён функций, у которых Cv (stddev/mean) для всех языков не превышает заданный порог, по выбранной архитектуре.
// @Tags        Benchmark
// @Produce     json
// @Param       arch         query string  true  "Архитектура (amd64, arm64 и т.д.)"
// @Param       cv-threshold query number false "Максимально допустимое значение Cv (по умолчанию 0.2)"
// @Param       min-mean     query number false "Минимальное значение среднего (по умолчанию 1e-12)"
// @Success     200 {object} ReliableBenchmarksResponse
// @Failure     400 {object} object "Ошибка в запросе"
// @Failure     404 {object} object "Данные не найдены"
// @Failure     500 {object} object "Внутренняя ошибка"
// @Router      /api/benchmark/reliable [get].
func (h *BenchmarkHandler) GetReliableBenchmarks(c echo.Context) error {
	arch := c.QueryParam("arch")
	if arch == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "parameter 'arch' is required"})
	}

	cvThreshold := 0.2
	minMean := 1e-12

	if s := c.QueryParam("cv-threshold"); s != "" {
		val, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "cv-threshold is not a float"})
		}
		cvThreshold = val
	}
	if s := c.QueryParam("min-mean"); s != "" {
		val, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "min-mean is not a float"})
		}
		minMean = val
	}

	allResults, err := h.BenchmarkRepo.GetAllBenchmarkResults()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return pickReliableBenchmarks(c, arch, cvThreshold, minMean, allResults)
}

type ReliableBenchmarksResponse struct {
	Arch  string   `json:"arch"`
	Count int      `json:"count"`
	Names []string `json:"names"`
}

func pickReliableBenchmarks(
	c echo.Context,
	arch string,
	cvThreshold, minMean float64,
	allResults []domain.BenchmarkResults,
) error {
	type LangStat struct {
		Mean   float64
		Stddev float64
	}

	langs := []string{"go", "cpp", "rust", "javascript"}
	stats := map[string]map[string]LangStat{}

	for _, group := range allResults {
		if group.Arch != arch {
			continue
		}
		metric := group.Metric
		for _, bench := range group.Results {
			name := bench.Name
			if _, ok := stats[name]; !ok {
				stats[name] = map[string]LangStat{}
			}
			for _, lang := range langs {
				var value float64
				switch lang {
				case "go":
					value = bench.Go
				case "cpp":
					value = bench.Cpp
				case "rust":
					value = bench.Rust
				case "javascript":
					value = bench.Javascript
				}
				if value == 0 && metric == "mean" {
					continue
				}
				s := stats[name][lang]
				switch metric {
				case "mean":
					s.Mean = value
				case "stddev":
					s.Stddev = value
				}
				stats[name][lang] = s
			}
		}
	}

	reliable := []string{}

FUNC_LOOP:
	for name, perLang := range stats {
		if len(perLang) != len(langs) {
			continue
		}
		for _, lang := range langs {
			s := perLang[lang]
			if s.Mean <= minMean {
				continue FUNC_LOOP
			}
			cv := 0.0
			if s.Mean != 0 {
				cv = s.Stddev / s.Mean
			}
			if cv > cvThreshold || math.IsNaN(cv) {
				continue FUNC_LOOP
			}
		}
		reliable = append(reliable, name)
	}
	sort.Strings(reliable)

	if len(reliable) == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "No reliable functions for this architecture"})
	}

	resp := ReliableBenchmarksResponse{
		Arch:  arch,
		Count: len(reliable),
		Names: reliable,
	}
	return c.JSON(http.StatusOK, resp)
}
