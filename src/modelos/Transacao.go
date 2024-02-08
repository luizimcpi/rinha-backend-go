package modelos

type Transacao struct {
	Valor     uint64    `json:"valor,omitempty"`
	Tipo      string    `json:"tipo,omitempty"`
	Descricao string    `json:"descricao,omitempty"`
}