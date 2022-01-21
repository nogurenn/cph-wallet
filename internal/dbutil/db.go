package dbutil

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

func NewDb(config *Config) (*sqlx.DB, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		config.DbUser, config.DbPassword, config.DbHost, config.DbPort, config.DbName)
	return sqlx.Open("pgx", connString)
}

func SwitchSchema(txn *sqlx.Tx, schema string) error {
	_, err := txn.Exec(fmt.Sprintf("SET search_path TO %s", schema))
	return err
}

type Transaction interface {
	sqlx.Execer
	Rollback() error
	Commit() error
	Get(dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
	NamedExec(query string, arg interface{}) (sql.Result, error)
}
