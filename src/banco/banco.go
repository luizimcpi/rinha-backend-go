package banco

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Conectar() (*sql.DB, error) {

	stringConexao := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_DATABASE"),
		os.Getenv("DB_PORT"),
	)

	db, erro := sql.Open("pgx", stringConexao)
	if erro != nil {
		return nil, erro
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	if erro = db.Ping(); erro != nil {
		db.Close()
		return nil, erro
	}

	return db, nil

}
