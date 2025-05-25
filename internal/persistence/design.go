package persistence

import (
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/arvaliullin/wapa/internal/domain"
)

type DesignRepositoryContract interface {
	Create(design domain.Design) (string, error)
	Update(design domain.Design) error
	Delete(id string) error
	GetAll() ([]domain.Design, error)
	GetByID(id string) (domain.Design, error)
}

type DesignRepository struct {
	DbConnection string
}

func (repo *DesignRepository) Create(design domain.Design) (string, error) {
	var id string

	err := WithConnection(repo.DbConnection, func(conn *sql.DB) error {
		query := `
			INSERT INTO composer.design (name, lang, js, wasm, functions) 
			VALUES ($1, $2, $3, $4, $5) 
			ON CONFLICT (name) DO UPDATE
			SET
				lang = EXCLUDED.lang,
				js = EXCLUDED.js,
				wasm = EXCLUDED.wasm,
				functions = EXCLUDED.functions
			RETURNING id
		`

		functionsJSON, err := json.Marshal(design.Functions)
		if err != nil {
			return err
		}

		err = conn.QueryRow(query,
			design.Name,
			design.Lang,
			design.JS,
			design.Wasm,
			functionsJSON).Scan(&id)
		return err
	})

	if err != nil {
		return "", err
	}

	return id, nil
}

func (repo *DesignRepository) Update(design domain.Design) error {
	return WithConnection(repo.DbConnection, func(conn *sql.DB) error {
		query := `
			UPDATE composer.design
			SET name = $1, lang = $2, js = $3, wasm = $4, functions = $5
			WHERE id = $6
		`

		functionsJSON, err := json.Marshal(design.Functions)
		if err != nil {
			return err
		}

		result, err := conn.Exec(query,
			design.Name,
			design.Lang,
			design.JS,
			design.Wasm,
			functionsJSON,
			design.ID)
		if err != nil {
			return err
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return err
		}

		if rowsAffected == 0 {
			return errors.New("no rows were updated, possibly invalid ID")
		}

		return nil
	})
}

func (repo *DesignRepository) GetAll() ([]domain.Design, error) {
	var designs []domain.Design

	err := WithConnection(repo.DbConnection, func(conn *sql.DB) error {
		query := `
			SELECT id, name, lang, js, wasm, functions 
			FROM composer.design
		`

		rows, err := conn.Query(query)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var design domain.Design
			var functionsJSON []byte

			err := rows.Scan(&design.ID,
				&design.Name,
				&design.Lang,
				&design.JS,
				&design.Wasm,
				&functionsJSON)
			if err != nil {
				return err
			}

			if err := json.Unmarshal(functionsJSON, &design.Functions); err != nil {
				return err
			}

			designs = append(designs, design)
		}

		return rows.Err()
	})

	if err != nil {
		return nil, err
	}

	return designs, nil
}

func (repo *DesignRepository) GetByID(id string) (domain.Design, error) {

	var design domain.Design

	err := WithConnection(repo.DbConnection, func(conn *sql.DB) error {
		query := `
			SELECT id, name, lang, js, wasm, functions 
			FROM composer.design
			WHERE id = $1
		`

		var functionsJSON []byte

		err := conn.QueryRow(query, id).Scan(&design.ID,
			&design.Name,
			&design.Lang,
			&design.JS,
			&design.Wasm,
			&functionsJSON)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil
			}
			return err
		}

		if err := json.Unmarshal(functionsJSON, &design.Functions); err != nil {
			return err
		}

		return nil
	})

	return design, err
}

func (repo *DesignRepository) Delete(id string) error {
	return WithConnection(repo.DbConnection, func(conn *sql.DB) error {
		query := `DELETE FROM composer.design WHERE id = $1`
		_, err := conn.Exec(query, id)
		return err
	})
}
