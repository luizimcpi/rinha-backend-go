package controllers

import (
	"encoding/json"
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

	log.Println("Cliente encontrado: " + strconv.FormatUint(cliente.ID, 10))

	var transacaoID uint64
	repositorio := repositorios.NovoRepositorioDeTransacoes(db)
	transacaoID, erro = repositorio.Criar(transacao, clienteID)
	if erro != nil {
		respostas.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	log.Println("TransacaoID is: " + strconv.FormatUint(transacaoID, 10))


	respostas.JSON(w, http.StatusOK, transacao)
}
