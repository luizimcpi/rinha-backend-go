package repositorios

import (
	"database/sql"
	"server/src/modelos"
)

type Clientes struct {
	db *sql.DB
}

func NovoRepositorioDeClientes(db *sql.DB) *Clientes {
	return &Clientes{db}
}

func (repositorio Clientes) BuscarPorID(ID uint64) (modelos.Cliente, error) {
	linhas, erro := repositorio.db.Query(
		"select id, limite, saldo from clientes where id = ? FOR UPDATE",
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

func (repositorio Clientes) AtualizarSaldo(clienteID uint64, saldo int64) error {
	statement, erro := repositorio.db.Prepare("update clientes set saldo = ? where id = ?")
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(saldo, clienteID); erro != nil {
		return erro
	}

	return nil
}
