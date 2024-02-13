package modelos

import "time"

type Extrato struct {
	Saldo             SaldoResponse       `json:"saldo"`
	UltimasTransacoes []TransacaoResponse `json:"ultimas_transacoes"`
}

type SaldoResponse struct {
	Total       int64     `json:"total"`
	DataExtrato time.Time `json:"data_extrato"`
	Limite      uint64    `json:"limite"`
}
