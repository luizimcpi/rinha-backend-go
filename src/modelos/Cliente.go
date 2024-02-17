package modelos

import "time"

type Cliente struct {
	ID       uint64    `json:"id,omitempty"`
	Limite   uint64    `json:"limite"`
	Saldo    int64     `json:"saldo"`
	CriadoEm time.Time `json:"data_criacao,omitempty"`
}
