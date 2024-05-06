package connector

import (
	"database/sql"
)

func GetPostgresConnector(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	// с этим параметром можно поиграться, обычно хорошей практикой является
	// количество ядер процессора умноженное в 2-3 раза
	db.SetMaxOpenConns(0)
	return db, nil
}
