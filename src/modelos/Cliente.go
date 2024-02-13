package modelos

import "time"

type Cliente struct {
	ID       uint64    `json:"id,omitempty"`
	Limite   int64     `json:"limite,omitempty"`
	Saldo    int64     `json:"saldo_inicial,omitempty"`
	CriadoEm time.Time `json:"data_criacao,omitempty"`
}
