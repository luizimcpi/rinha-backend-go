package controllers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"server/src/banco"
	"server/src/modelos"
	"server/src/repositorios"
	"server/src/respostas"
	"strconv"

	"github.com/gorilla/mux"
)

func CriarTransacao(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)

	clienteID, erro := strconv.ParseUint(parametros["id"], 10, 64)

	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	corpoRequest, erro := io.ReadAll(r.Body)
	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var transacao modelos.Transacao
	if erro = json.Unmarshal(corpoRequest, &transacao); erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	if erro = transacao.Preparar(); erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}
	defer db.Close()

	repositorioCliente := repositorios.NovoRepositorioDeClientes(db)
	cliente, erro := repositorioCliente.BuscarPorID(clienteID)
	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	if (modelos.Cliente{}) == cliente {
		respostas.Erro(w, http.StatusNotFound, errors.New("cliente não existe na base"))
		return
	}

	var saldoAtual = cliente.Saldo

	if transacao.Tipo == "d" {

		if transacao.Valor > cliente.Limite+uint64(saldoAtual) {
			respostas.Erro(w, http.StatusUnprocessableEntity, errors.New("transação de debito deixará saldo incosistente"))
			return
		}
	}

	repositorioTransacoes := repositorios.NovoRepositorioDeTransacoes(db)

	erro = repositorioTransacoes.Criar(transacao, clienteID)
	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var transacaoResponse modelos.TransacaoCriadaResponse
	transacaoResponse.Limite = cliente.Limite

	if transacao.Tipo == "d" {
		transacaoResponse.Saldo = saldoAtual - int64(transacao.Valor)
	} else {
		transacaoResponse.Saldo = saldoAtual + int64(transacao.Valor)
	}

	if erro = repositorioCliente.AtualizarSaldo(clienteID, transacaoResponse.Saldo); erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	respostas.JSON(w, http.StatusOK, transacaoResponse)
}
