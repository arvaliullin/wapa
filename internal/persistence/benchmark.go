package persistence

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/arvaliullin/wapa/internal/domain"
)

type BenchmarkRepositoryContract interface {
	GetBenchmarkResults(metric string, arch string) (domain.BenchmarkResults, error)
	GetAllBenchmarkResults() ([]domain.BenchmarkResults, error)
	GetBenchmarksOnlyMock(metric, arch string) (domain.BenchmarkResults, error)
	GetBenchmarksOnlyNotMock(metric, arch string) (domain.BenchmarkResults, error)
}

type BenchmarkRepository struct {
	DbConnection string
}

func mapMetricView(metric string) (string, error) {
	switch metric {
	case "mean":
		return "composer.v_metric_mean_json", nil
	case "median":
		return "composer.v_metric_median_json", nil
	case "min":
		return "composer.v_metric_min_json", nil
	case "max":
		return "composer.v_metric_max_json", nil
	case "stddev":
		return "composer.v_metric_stddev_json", nil
	default:
		return "", fmt.Errorf("unknown metric: %s", metric)
	}
}

func (repo *BenchmarkRepository) GetBenchmarkResults(metric string, arch string) (domain.BenchmarkResults, error) {
	var results domain.BenchmarkResults

	view, err := mapMetricView(metric)
	if err != nil {
		return results, err
	}

	err = WithConnection(repo.DbConnection, func(conn *sql.DB) error {
		query := fmt.Sprintf(
			`SELECT data FROM %s WHERE data->>'arch' = $1`, view)
		var dataBytes []byte

		err := conn.QueryRow(query, arch).Scan(&dataBytes)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return sql.ErrNoRows
			}
			return err
		}

		if err := json.Unmarshal(dataBytes, &results); err != nil {
			return fmt.Errorf("unmarshal: %w", err)
		}

		return nil
	})

	return results, err
}

func (repo *BenchmarkRepository) GetAllBenchmarkResults() ([]domain.BenchmarkResults, error) {
	metrics := []string{"mean", "median", "min", "max", "stddev"}
	var allResults []domain.BenchmarkResults

	for _, metric := range metrics {
		view, err := mapMetricView(metric)
		if err != nil {
			return nil, err
		}
		err = WithConnection(repo.DbConnection, func(conn *sql.DB) error {
			rows, err := conn.Query(fmt.Sprintf(`SELECT data FROM %s`, view))
			if err != nil {
				return err
			}
			defer rows.Close()

			for rows.Next() {
				var dataBytes []byte
				if err := rows.Scan(&dataBytes); err != nil {
					return err
				}
				var results domain.BenchmarkResults
				if err := json.Unmarshal(dataBytes, &results); err != nil {
					return err
				}
				allResults = append(allResults, results)
			}
			return rows.Err()
		})
		if err != nil {
			return nil, err
		}
	}
	return allResults, nil
}

func (repo *BenchmarkRepository) GetBenchmarksOnlyMock(metric, arch string) (domain.BenchmarkResults, error) {
	full, err := repo.GetBenchmarkResults(metric, arch)
	if err != nil {
		return full, err
	}
	var filtered []domain.BenchmarkCase
	for _, c := range full.Results {
		if strings.HasSuffix(c.Name, "Mock") {
			filtered = append(filtered, c)
		}
	}
	full.Results = filtered
	return full, nil
}

func (repo *BenchmarkRepository) GetBenchmarksOnlyNotMock(metric, arch string) (domain.BenchmarkResults, error) {
	full, err := repo.GetBenchmarkResults(metric, arch)
	if err != nil {
		return full, err
	}
	var filtered []domain.BenchmarkCase
	for _, c := range full.Results {
		if !strings.HasSuffix(c.Name, "Mock") {
			filtered = append(filtered, c)
		}
	}
	full.Results = filtered
	return full, nil
}
