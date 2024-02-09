package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"server/src/respostas"
	"server/src/modelos"
	"github.com/gorilla/mux"
	"strconv"
	"log"
)

func CriarTransacao(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)

	clienteID, erro := strconv.ParseUint(parametros["id"], 10, 64)
	if erro != nil {
		respostas.Erro(w, http.StatusBadRequest, erro)
		return
	}

	log.Println("ClienteID is: " + strconv.FormatUint(clienteID, 10))

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

	respostas.JSON(w, http.StatusOK, transacao)
}
