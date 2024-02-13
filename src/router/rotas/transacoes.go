package rotas

import (
	"net/http"
	"server/src/controllers"
)

var rotasTransacoes = []Rota{
	{
		URI:    "/clientes/{id}/transacoes",
		Metodo: http.MethodPost,
		Funcao: controllers.CriarTransacao,
	},
}
