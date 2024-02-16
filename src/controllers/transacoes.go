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

	//log.Println("Cliente encontrado: " + strconv.FormatUint(cliente.ID, 10))

	//var transacaoID uint64
	repositorio := repositorios.NovoRepositorioDeTransacoes(db)

	var somatorioTransacoes int64
	somatorioTransacoes, erro = repositorio.BuscarSomatorio(clienteID)
	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var clienteLimite = int64(cliente.Limite)
	if transacao.Tipo == "d" {
		var limiteNegativo = -clienteLimite
		var saldoComDebito = somatorioTransacoes - int64(transacao.Valor)

		if saldoComDebito < limiteNegativo {
			respostas.Erro(w, http.StatusUnprocessableEntity, errors.New("transação de debito deixará saldo incosistente"))
			return
		}
	}

	if transacao.Tipo == "c" {
		var saldoComCredito = somatorioTransacoes + int64(transacao.Valor)

		if saldoComCredito > clienteLimite {
			respostas.Erro(w, http.StatusUnprocessableEntity, errors.New("transação de credito deixará saldo incosistente"))
			return
		}
	}

	_, erro = repositorio.Criar(transacao, clienteID)
	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var transacaoResponse modelos.TransacaoCriadaResponse
	transacaoResponse.Limite = cliente.Limite

	if transacao.Tipo == "d" {
		transacaoResponse.Saldo = somatorioTransacoes - int64(transacao.Valor)
	} else {
		transacaoResponse.Saldo = somatorioTransacoes + int64(transacao.Valor)
	}

	respostas.JSON(w, http.StatusOK, transacaoResponse)
}
