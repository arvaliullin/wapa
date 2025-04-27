package persistence

import (
	"database/sql"
	"errors"

	"github.com/arvaliullin/wapa/internal/domain"
)

type MetricRepositoryContract interface {
	Create(metrics domain.Metrics, functionResultID string) (string, error)
	Update(metrics domain.Metrics, functionResultID string) error
	GetByFunctionResultID(functionResultID string) (domain.Metrics, error)
	Delete(functionResultID string) error
}

type MetricRepository struct {
	DbConnection string
}

func NewMetricRepository(conn string) *MetricRepository {
	return &MetricRepository{DbConnection: conn}
}

func (repo *MetricRepository) Create(metrics domain.Metrics, functionResultID string) (string, error) {
	var id string

	err := WithConnection(repo.DbConnection, func(conn *sql.DB) error {
		query := `
			INSERT INTO composer.metric (function_result_id, mean, stddev, median, user_time, system, min, max)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id
		`
		return conn.QueryRow(query,
			functionResultID,
			metrics.Mean,
			metrics.Stddev,
			metrics.Median,
			metrics.User,
			metrics.System,
			metrics.Min,
			metrics.Max,
		).Scan(&id)
	})

	return id, err
}

func (repo *MetricRepository) Update(metrics domain.Metrics, functionResultID string) error {
	return WithConnection(repo.DbConnection, func(conn *sql.DB) error {
		query := `
			UPDATE composer.metric SET mean=$1, stddev=$2, median=$3, user_time=$4, system=$5, min=$6, max=$7
			WHERE function_result_id=$8
		`
		result, err := conn.Exec(query,
			metrics.Mean,
			metrics.Stddev,
			metrics.Median,
			metrics.User,
			metrics.System,
			metrics.Min,
			metrics.Max,
			functionResultID,
		)
		if err != nil {
			return err
		}
		n, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if n == 0 {
			return errors.New("no metric rows updated")
		}
		return nil
	})
}

func (repo *MetricRepository) GetByFunctionResultID(functionResultID string) (domain.Metrics, error) {
	var m domain.Metrics

	err := WithConnection(repo.DbConnection, func(conn *sql.DB) error {
		query := `
			SELECT mean, stddev, median, user_time, system, min, max
			FROM composer.metric
			WHERE function_result_id = $1
		`
		return conn.QueryRow(query, functionResultID).Scan(
			&m.Mean, &m.Stddev, &m.Median, &m.User, &m.System, &m.Min, &m.Max)
	})
	return m, err
}

func (repo *MetricRepository) Delete(functionResultID string) error {
	return WithConnection(repo.DbConnection, func(conn *sql.DB) error {
		query := `DELETE FROM composer.metric WHERE function_result_id = $1`
		_, err := conn.Exec(query, functionResultID)
		return err
	})
}
