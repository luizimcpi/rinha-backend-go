package banco

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql" // Driver
)

// Conectar abre a conexão com o banco de dados e a retorna
func Conectar() (*sql.DB, error) {
	//docker
	stringConexao := "user:123456@tcp(mysqldocker:3306)/rinhabank?charset=utf8&parseTime=True&loc=Local"
	//local mysql
	//stringConexao := "root:123456@tcp(localhost:3306)/rinhabank?charset=utf8&parseTime=True&loc=Local"
	db, erro := sql.Open("mysql", stringConexao)
	if erro != nil {
		return nil, erro
	}

	db.SetConnMaxLifetime(time.Second * 10)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	if erro = db.Ping(); erro != nil {
		db.Close()
		return nil, erro
	}

	return db, nil

}
