package modelos

import "time"

type Cliente struct {
	ID       uint64    `json:"id,omitempty"`
	Limite   uint64    `json:"limite,omitempty"`
	Saldo    uint64    `json:"saldo_inicial,omitempty"`
	CriadoEm time.Time `json:"data_criacao,omitempty"`
}
