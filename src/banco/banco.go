package banco

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Conectar() (*sql.DB, error) {

	//stringConexao := "user:123456@tcp(db:5432)/rinhabank?charset=utf8&parseTime=True&loc=Local"
	stringConexao := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		"db",
		"user",
		"123456",
		"rinhabank",
		"5432",
	)

	db, erro := sql.Open("pgx", stringConexao)
	if erro != nil {
		return nil, erro
	}

	db.SetMaxOpenConns(200)

	if erro = db.Ping(); erro != nil {
		db.Close()
		return nil, erro
	}

	return db, nil

}
