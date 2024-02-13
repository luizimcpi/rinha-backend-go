package controllers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
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
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	corpoRequest, erro := io.ReadAll(r.Body)
	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var transacao modelos.Transacao
	if erro = json.Unmarshal(corpoRequest, &transacao); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if erro = transacao.Preparar(); erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorioCliente := repositorios.NovoRepositorioDeClientes(db)
	cliente, erro := repositorioCliente.BuscarPorID(clienteID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	if (modelos.Cliente{}) == cliente {
		respostas.Erro(w, http.StatusNotFound, errors.New("Cliente não existe na base"))
		return
	}

	//log.Println("Cliente encontrado: " + strconv.FormatUint(cliente.ID, 10))

	var transacaoID uint64
	repositorio := repositorios.NovoRepositorioDeTransacoes(db)

	var somatorioTransacoes int64
	somatorioTransacoes, erro = repositorio.BuscarSomatorio(clienteID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	if transacao.Tipo == "d" {
		var limiteNegativo = -cliente.Limite
		var saldoAtual = somatorioTransacoes - int64(transacao.Valor)

		//log.Println("Limite negativo is: " + strconv.FormatInt(limiteNegativo, 10))

		if saldoAtual < limiteNegativo {
			respostas.Erro(w, http.StatusUnprocessableEntity, errors.New("Transação de debito deixará saldo incosistente"))
			return
		}
	}

	//log.Println("Somatorio transacoes is: " + strconv.FormatInt(somatorioTransacoes, 10))

	transacaoID, erro = repositorio.Criar(transacao, clienteID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	log.Println("TransacaoID is: " + strconv.FormatUint(transacaoID, 10))

	var transacaoResponse modelos.TransacaoCriadaResponse
	transacaoResponse.Limite = uint64(cliente.Limite)
	if transacao.Tipo == "d" {
		transacaoResponse.Saldo = somatorioTransacoes - int64(transacao.Valor)
	} else {
		transacaoResponse.Saldo = somatorioTransacoes + int64(transacao.Valor)
	}

	respostas.JSON(w, http.StatusOK, transacaoResponse)
}
