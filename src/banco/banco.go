package banco

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func Conectar() (*sql.DB, error) {

	stringConexao := "user:123456@tcp(mysqldocker:3306)/rinhabank?charset=utf8&parseTime=True&loc=Local"

	db, erro := sql.Open("mysql", stringConexao)
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
