package controllers

import (
	"errors"
	"net/http"
	"server/src/banco"
	"server/src/modelos"
	"server/src/repositorios"
	"server/src/respostas"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func Extrato(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)

	clienteID, erro := strconv.ParseUint(parametros["id"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := banco.Conectar()
	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}
	defer db.Close()

	repositorioTransacoes := repositorios.NovoRepositorioDeTransacoes(db)
	transacoes, erro := repositorioTransacoes.BuscarUltimas(clienteID)
	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

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

	var somatorioTransacoes int64
	somatorioTransacoes, erro = repositorioTransacoes.BuscarSomatorio(clienteID)
	if erro != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var saldoResponse modelos.SaldoResponse
	saldoResponse.Total = somatorioTransacoes
	saldoResponse.DataExtrato = time.Now()
	saldoResponse.Limite = cliente.Limite

	var extrato modelos.Extrato
	extrato.Saldo = saldoResponse
	extrato.UltimasTransacoes = transacoes

	respostas.JSON(w, http.StatusOK, extrato)

}
