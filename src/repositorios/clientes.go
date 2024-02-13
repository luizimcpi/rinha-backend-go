package repositorios

import (
	"database/sql"
	"server/src/modelos"
)

// Clientes representa um repositório de clientes
type Clientes struct {
	db *sql.DB
}

// NovoRepositorioDeClientes cria um repositório de clientes
func NovoRepositorioDeClientes(db *sql.DB) *Clientes {
	return &Clientes{db}
}

// BuscarPorID traz um cliente do banco de dados
func (repositorio Clientes) BuscarPorID(ID uint64) (modelos.Cliente, error) {
	linhas, erro := repositorio.db.Query(
		"select id, limite, saldo_inicial, data_criacao from clientes where id = ?",
		ID,
	)
	if erro != nil {
		return modelos.Cliente{}, erro
	}
	defer linhas.Close()

	var cliente modelos.Cliente

	if linhas.Next() {
		if erro = linhas.Scan(
			&cliente.ID,
			&cliente.Limite,
			&cliente.Saldo,
			&cliente.CriadoEm,
		); erro != nil {
			return modelos.Cliente{}, erro
		}
	}

	return cliente, nil
}
