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