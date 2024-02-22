package repositorios

import (
	"database/sql"
	"fmt"
	"server/src/modelos"
)

type Transacoes struct {
	db *sql.DB
}

func NovoRepositorioDeTransacoes(db *sql.DB) *Transacoes {
	return &Transacoes{db}
}

func (repositorio Transacoes) Criar(transacao modelos.Transacao, clienteID uint64) error {
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

	_, err = tx.Exec("INSERT INTO transacoes (valor, tipo, descricao, cliente_id) VALUES ($1, $2, $3, $4);", transacao.Valor, transacao.Tipo, transacao.Descricao, clienteID)
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

func (repositorio Transacoes) BuscarUltimas(clienteID uint64) ([]modelos.TransacaoResponse, error) {
	linhas, erro := repositorio.db.Query(`
	select t.valor, t.tipo, t.descricao, t.realizada_em from transacoes t
	where t.cliente_id = $1
	order by t.realizada_em desc limit 10`,
		clienteID,
	)
	if erro != nil {
		return nil, erro
	}
	defer linhas.Close()

	var transacoes []modelos.TransacaoResponse

	for linhas.Next() {
		var transacao modelos.TransacaoResponse

		if erro = linhas.Scan(
			&transacao.Valor,
			&transacao.Tipo,
			&transacao.Descricao,
			&transacao.RealizadaEm,
		); erro != nil {
			return nil, erro
		}

		transacoes = append(transacoes, transacao)
	}

	return transacoes, nil
}
