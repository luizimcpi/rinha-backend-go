package repositorios

import (
	"database/sql"
	"server/src/modelos"
)

type Transacoes struct {
	db *sql.DB
}

func NovoRepositorioDeTransacoes(db *sql.DB) *Transacoes {
	return &Transacoes{db}
}

func (repositorio Transacoes) Criar(transacao modelos.Transacao, clienteID uint64) (uint64, error) {
	statement, erro := repositorio.db.Prepare(
		"insert into transacoes (valor, tipo, descricao, cliente_id) values (?, ?, ?, ?)",
	)
	if erro != nil {
		return 0, erro
	}
	defer statement.Close()

	resultado, erro := statement.Exec(transacao.Valor, transacao.Tipo, transacao.Descricao, clienteID)
	if erro != nil {
		return 0, erro
	}

	ultimoIDInserido, erro := resultado.LastInsertId()
	if erro != nil {
		return 0, erro
	}

	return uint64(ultimoIDInserido), nil
}

func (repositorio Transacoes) BuscarUltimas(clienteID uint64) ([]modelos.TransacaoResponse, error) {
	linhas, erro := repositorio.db.Query(`
	select t.valor, t.tipo, t.descricao, t.realizada_em from transacoes t
	where t.cliente_id = ?
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
