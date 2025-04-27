package persistence

import (
	"database/sql"
	"encoding/json"

	"github.com/arvaliullin/wapa/internal/domain"
)

type FunctionResultRepositoryContract interface {
	Create(fr domain.FunctionResult, experimentID string) (string, error)
	GetByExperimentID(experimentID string) ([]domain.FunctionResult, error)
	GetByID(id string) (domain.FunctionResult, error)
	Delete(id string) error
}

type FunctionResultRepository struct {
	DbConnection string
	MetricsRepo  *MetricRepository
}

func NewFunctionResultRepository(conn string) *FunctionResultRepository {
	return &FunctionResultRepository{
		DbConnection: conn,
		MetricsRepo:  NewMetricRepository(conn),
	}
}

func (repo *FunctionResultRepository) Create(fr domain.FunctionResult, experimentID string) (string, error) {
	var id string
	err := WithConnection(repo.DbConnection, func(conn *sql.DB) error {
		argsJson, err := json.Marshal(fr.Args)
		if err != nil {
			return err
		}
		resultJson, err := json.Marshal(fr.Result)
		if err != nil {
			return err
		}

		query := `
		INSERT INTO composer.function_result
		(experiment_id, function_name, args, repeats, result)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
		`

		err = conn.QueryRow(query,
			experimentID,
			fr.FunctionName,
			argsJson,
			fr.Repeats,
			resultJson,
		).Scan(&id)
		if err != nil {
			return err
		}

		if repo.MetricsRepo != nil {
			_, err = repo.MetricsRepo.Create(fr.Metrics, id)
		}
		return err
	})
	return id, err
}

func (repo *FunctionResultRepository) GetByExperimentID(experimentID string) ([]domain.FunctionResult, error) {
	var results []domain.FunctionResult

	err := WithConnection(repo.DbConnection, func(conn *sql.DB) error {
		query := `
			SELECT id, function_name, args, repeats, result
			FROM composer.function_result WHERE experiment_id = $1
		`
		rows, err := conn.Query(query, experimentID)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var fr domain.FunctionResult
			var argsJson []byte
			var resultJson []byte

			if err := rows.Scan(&fr.ID, &fr.FunctionName, &argsJson, &fr.Repeats, &resultJson); err != nil {
				return err
			}
			if err := json.Unmarshal(argsJson, &fr.Args); err != nil {
				return err
			}
			if err := json.Unmarshal(resultJson, &fr.Result); err != nil {
				return err
			}

			if repo.MetricsRepo != nil {
				m, _ := repo.MetricsRepo.GetByFunctionResultID(fr.ID)
				fr.Metrics = m
			}
			results = append(results, fr)
		}
		return rows.Err()
	})

	return results, err
}

func (repo *FunctionResultRepository) GetByID(id string) (domain.FunctionResult, error) {
	var fr domain.FunctionResult
	err := WithConnection(repo.DbConnection, func(conn *sql.DB) error {
		query := `
			SELECT experiment_id, function_name, args, repeats, result
			FROM composer.function_result WHERE id = $1
		`
		var argsJson, resultJson []byte
		err := conn.QueryRow(query, id).Scan(
			&fr.ExperimentID, &fr.FunctionName, &argsJson, &fr.Repeats, &resultJson,
		)
		if err != nil {
			return err
		}
		fr.ID = id
		if err := json.Unmarshal(argsJson, &fr.Args); err != nil {
			return err
		}
		if err := json.Unmarshal(resultJson, &fr.Result); err != nil {
			return err
		}
		if repo.MetricsRepo != nil {
			m, _ := repo.MetricsRepo.GetByFunctionResultID(fr.ID)
			fr.Metrics = m
		}
		return nil
	})
	return fr, err
}

func (repo *FunctionResultRepository) Delete(id string) error {
	return WithConnection(repo.DbConnection, func(conn *sql.DB) error {
		query := `DELETE FROM composer.function_result WHERE id = $1`
		_, err := conn.Exec(query, id)
		return err
	})
}
