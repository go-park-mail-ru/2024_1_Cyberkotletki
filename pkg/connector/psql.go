package connector

import "github.com/jmoiron/sqlx"

func GetPostgresConnector(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(0)
	return db, nil
}
