package banco

import (
	"database/sql"
	"os"

	_ "github.com/go-sql-driver/mysql" // Driver
)

// Conectar abre a conexão com o banco de dados e a retorna
func Conectar() (*sql.DB, error) {
	stringConexao := os.Getenv("DB_STRING_CONEXAO")
	db, erro := sql.Open("mysql", stringConexao)
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
