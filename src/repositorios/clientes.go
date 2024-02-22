package repositorios

import (
	"database/sql"
	"fmt"
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
		"select id, limite, saldo from clientes where id = $1",
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
		); erro != nil {
			return modelos.Cliente{}, erro
		}
	}

	return cliente, nil
}

func (repositorio Clientes) AtualizarSaldo(clienteID uint64, saldo int64) error {
	tx, err := repositorio.db.Begin()

	if err != nil {
		return err
	}
	defer tx.Rollback()

	rows := tx.QueryRow("SELECT limite, saldo FROM clientes WHERE id = $1 FOR UPDATE;", clienteID)

	var cliente modelos.Cliente

	err = rows.Scan(&cliente.Limite, &cliente.Saldo)
	if err != nil {
		fmt.Println(err)
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("UPDATE clientes SET saldo = $1 WHERE id = $2;", saldo, clienteID)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
