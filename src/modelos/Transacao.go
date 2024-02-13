package modelos

import (
	"errors"
	"strings"
	"log"
	"time"
)

type Transacao struct {
	Valor     uint64 `json:"valor,omitempty"`
	Tipo      string `json:"tipo,omitempty"`
	Descricao string `json:"descricao,omitempty"`
}

type TransacaoCriadaResponse struct {
	Limite     uint64 `json:"limite,omitempty"`
	Saldo      int64 `json:"saldo"`
}

type TransacaoResponse struct {
	Valor     uint64 `json:"valor,omitempty"`
	Tipo      string `json:"tipo,omitempty"`
	Descricao string `json:"descricao,omitempty"`
	RealizadaEm  time.Time `json:"realizada_em,omitempty"`
}

func (transacao *Transacao) Preparar() error {
	if erro := transacao.validar(); erro != nil {
		return erro
	}

	if erro := transacao.formatar(); erro != nil {
		return erro
	}

	return nil
}

func (transacao *Transacao) validar() error {
	if transacao.Valor == 0 {
		return errors.New("O campo valor é obrigatório e não pode ser 0")
	}

	if transacao.Tipo == "" {
		return errors.New("O campo tipo é obrigatório e não pode estar em branco")
	}

	if transacao.Tipo == "c" || transacao.Tipo == "d" {
		log.Println("Tipo transacao: " + transacao.Tipo)
	} else {
		return errors.New("O campo tipo deve ser 'd' para débito ou 'c' para crédito")
	}

	if transacao.Descricao == "" {
		return errors.New("O campo descrição é obrigatório e não pode estar em branco")
	}

	if(len(transacao.Descricao) > 10){
		return errors.New("O campo descrição não pode conter mais que 10 caracteres")
	}

	return nil
}

func (transacao *Transacao) formatar() error {
	transacao.Tipo = strings.TrimSpace(transacao.Tipo)
	transacao.Descricao = strings.TrimSpace(transacao.Descricao)

	return nil
}
