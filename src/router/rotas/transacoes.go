package rotas

import (
	"net/http"
	"server/src/controllers"
)

var rotaTransacoes = Rota{
	URI:    "/clientes/{id}/transacoes",
	Metodo: http.MethodPost,
	Funcao: controllers.CriarTransacao,
}
