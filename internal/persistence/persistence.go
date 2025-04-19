package persistence

import "database/sql"

type Action func(conn *sql.DB) error

func WithConnection(DbConnection string, action Action) error {
	db, err := sql.Open("postgres", DbConnection)
	if err != nil {
		return err
	}
	defer db.Close()

	return action(db)
}
