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

func (repositorio Transacoes) Criar(transacao modelos.Transacao, clienteID uint64) error {
	statement, erro := repositorio.db.Prepare(
		"insert into transacoes (valor, tipo, descricao, cliente_id) values ($1, $2, $3, $4)",
	)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	_, erro = statement.Exec(transacao.Valor, transacao.Tipo, transacao.Descricao, clienteID)
	if erro != nil {
		return erro
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
