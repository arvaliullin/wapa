package persistence

import (
	"database/sql"

	"github.com/arvaliullin/wapa/internal/domain"
)

type ExperimentRepositoryContract interface {
	Create(exp domain.Experiment) (string, error)
	GetByID(id string) (domain.Experiment, error)
	GetAll() ([]domain.Experiment, error)
	Delete(id string) error
}

type ExperimentRepository struct {
	DbConnection       string
	FunctionResultRepo *FunctionResultRepository
}

func NewExperimentRepository(conn string) *ExperimentRepository {
	return &ExperimentRepository{
		DbConnection:       conn,
		FunctionResultRepo: NewFunctionResultRepository(conn),
	}
}

func (repo *ExperimentRepository) Create(exp domain.Experiment) (string, error) {
	var id string
	err := WithConnection(repo.DbConnection, func(conn *sql.DB) error {
		query := `
			INSERT INTO composer.experiment (design_id, hostname, arch)
			VALUES ($1, $2, $3)
			RETURNING id
		`
		err := conn.QueryRow(query, exp.DesignID, exp.Hostname, exp.Arch).Scan(&id)
		if err != nil {
			return err
		}

		for _, fr := range exp.FunctionResults {
			if repo.FunctionResultRepo != nil {
				_, err := repo.FunctionResultRepo.Create(fr, id)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})

	return id, err
}

func (repo *ExperimentRepository) GetByID(id string) (domain.Experiment, error) {
	var exp domain.Experiment
	err := WithConnection(repo.DbConnection, func(conn *sql.DB) error {
		query := `
            SELECT id, design_id, hostname, arch
            FROM composer.experiment WHERE id = $1
        `
		err := conn.QueryRow(query, id).Scan(
			&exp.ID, &exp.DesignID, &exp.Hostname, &exp.Arch)
		if err != nil {
			return err
		}

		if repo.FunctionResultRepo != nil {
			list, ferr := repo.FunctionResultRepo.GetByExperimentID(id)
			if ferr == nil {
				exp.FunctionResults = list
			}
		}
		return nil
	})
	return exp, err
}

func (repo *ExperimentRepository) GetAll() ([]domain.Experiment, error) {
	var experiments []domain.Experiment

	err := WithConnection(repo.DbConnection, func(conn *sql.DB) error {
		query := `
            SELECT id, design_id, hostname, arch
            FROM composer.experiment
        `
		rows, err := conn.Query(query)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var exp domain.Experiment
			if err := rows.Scan(&exp.ID, &exp.DesignID, &exp.Hostname, &exp.Arch); err != nil {
				return err
			}

			if repo.FunctionResultRepo != nil {
				list, _ := repo.FunctionResultRepo.GetByExperimentID(exp.ID)
				exp.FunctionResults = list
			}
			experiments = append(experiments, exp)
		}
		return rows.Err()
	})

	return experiments, err
}

func (repo *ExperimentRepository) Delete(id string) error {
	return WithConnection(repo.DbConnection, func(conn *sql.DB) error {
		query := `DELETE FROM composer.experiment WHERE id = $1`
		_, err := conn.Exec(query, id)
		return err
	})
}
