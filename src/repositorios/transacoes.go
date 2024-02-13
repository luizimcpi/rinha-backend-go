package repositorios

import (
	"database/sql"
	"server/src/modelos"
)

// Transacoes representa um repositório de transações
type Transacoes struct {
	db *sql.DB
}

// NovoRepositorioDeTransacoes cria um repositório de transacoes
func NovoRepositorioDeTransacoes(db *sql.DB) *Transacoes {
	return &Transacoes{db}
}

// Criar insere uma transação no banco de dados
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

func (repositorio Transacoes) BuscarSomatorio(clienteID uint64) (int64, error) {
	linha, erro := repositorio.db.Query(`
	select sum(CASE WHEN tipo = 'c' then valor WHEN tipo = 'd' then -valor END) as saldo
	from transacoes 
	where cliente_id = ?`,
		clienteID,
	)
	if erro != nil {
		return 0, erro
	}

	defer linha.Close()

	var saldo sql.NullInt64

	if linha.Next() {
		if erro = linha.Scan(&saldo); erro != nil {
			return 0, erro
		}
	}
	return saldo.Int64, nil
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
