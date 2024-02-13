package rotas

import (
	"net/http"
	"server/src/controllers"
)

var rotaExtrato = Rota{
	URI:    "/clientes/{id}/extrato",
	Metodo: http.MethodGet,
	Funcao: controllers.Extrato,
}
