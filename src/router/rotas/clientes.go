package rotas

import (
	"server/src/controllers"
	"net/http"
)

var rotasClientes = []Rota{
	{
		URI:                "/clientes/{id}/transacoes",
		Metodo:             http.MethodPost,
		Funcao:             controllers.CriarTransacao,
	},
}