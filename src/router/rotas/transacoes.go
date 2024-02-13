package rotas

import (
	"server/src/controllers"
	"net/http"
)

var rotasTransacoes = []Rota{
	{
		URI:                "/clientes/{id}/transacoes",
		Metodo:             http.MethodPost,
		Funcao:             controllers.CriarTransacao,
	},
}