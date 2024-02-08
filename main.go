package main

import (
	"fmt"
	"log"
	"net/http"
	"server/src/router"
)

func main() {
	r := router.Gerar()

	fmt.Printf("Escutando na porta %d\n", 8080)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", 8080), r))
}
