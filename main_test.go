package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"server/src/controllers"
	"testing"

	"github.com/gorilla/mux"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

type databaseContainer struct {
	*postgres.PostgresContainer
	connectionString string
}

func setupContainer(ctx context.Context) (*databaseContainer, error) {
	postgresContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("docker.io/postgres:15.2-alpine"),
		postgres.WithDatabase("rinhabank"),
		postgres.WithUsername("user"),
		postgres.WithPassword("123456"),
		postgres.WithInitScripts("./sql/script.sql"),
		postgres.WithConfigFile("postgresql.conf"),
	)

	if err != nil {
		log.Fatalf("Could not start postgres container: %s", err)
	}

	connString, err := postgresContainer.ConnectionString(ctx, "parseTime=true")

	if err != nil {
		log.Fatalf("Could not get connection string postgres: %s", err)
	}

	return &databaseContainer{
		PostgresContainer: postgresContainer,
		connectionString:  connString,
	}, nil
}

func TestMain(m *testing.M) {
	ctx := context.Background()
	container, err := setupContainer(ctx)

	if err != nil {
		log.Fatalf("could not start postgresContainer")
	}

	log.Println("connection string: " + container.connectionString)

	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_USER", "user")
	os.Setenv("DB_PASSWORD", "123456")
	os.Setenv("DB_DATABASE", "rinhabank")
	os.Setenv("DB_PORT", container.connectionString)

	defer func() {
		if err := container.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %s", err)
		}
	}()

	code := m.Run()
	os.Exit(code)
}

func TestCriarTransacaoRetornaErroQuandoRecebeIdInvalido(t *testing.T) {

	req, _ := http.NewRequest("POST", "/clientes/a/transacoes", nil)
	response := executeTransacaoRequest(req, "/clientes/{id}/transacoes")

	checkResponseCode(t, http.StatusUnprocessableEntity, response.Code)
}

func TestCriarTransacaoRetornaErroQuandoRecebeIdInexistente(t *testing.T) {

	var jsonStr = []byte(`{"valor": 1000, "tipo": "c", "descricao": "teste"}`)
	req, _ := http.NewRequest("POST", "/clientes/6/transacoes", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeTransacaoRequest(req, "/clientes/{id}/transacoes")

	checkResponseCode(t, http.StatusNotFound, response.Code)
}

func TestCriarTransacaoRetornaErroQuandoRecebeIdValidoComBodyComCampoValorNegativo(t *testing.T) {

	var jsonStr = []byte(`{"valor": -1000, "tipo": "c", "descricao": "teste"}`)
	req, _ := http.NewRequest("POST", "/clientes/1/transacoes", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeTransacaoRequest(req, "/clientes/{id}/transacoes")

	checkResponseCode(t, http.StatusUnprocessableEntity, response.Code)
}

func TestCriarTransacaoRetornaErroQuandoRecebeIdValidoComBodyComCampoValorDecimal(t *testing.T) {

	var jsonStr = []byte(`{"valor": 1.2, "tipo": "c", "descricao": "teste"}`)
	req, _ := http.NewRequest("POST", "/clientes/1/transacoes", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeTransacaoRequest(req, "/clientes/{id}/transacoes")

	checkResponseCode(t, http.StatusUnprocessableEntity, response.Code)
}

func TestCriarTransacaoRetornaErroQuandoRecebeIdValidoComBodySemCampoValor(t *testing.T) {

	var jsonStr = []byte(`{"tipo":"d", "descricao": "teste"}`)
	req, _ := http.NewRequest("POST", "/clientes/1/transacoes", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeTransacaoRequest(req, "/clientes/{id}/transacoes")

	checkResponseCode(t, http.StatusUnprocessableEntity, response.Code)
}

func TestCriarTransacaoRetornaErroQuandoRecebeIdValidoComBodySemCampoTipo(t *testing.T) {

	var jsonStr = []byte(`{"valor": 1000, "descricao": "teste"}`)
	req, _ := http.NewRequest("POST", "/clientes/1/transacoes", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeTransacaoRequest(req, "/clientes/{id}/transacoes")

	checkResponseCode(t, http.StatusUnprocessableEntity, response.Code)

	var body map[string]string
	var expected = "o campo tipo é obrigatório e não pode estar em branco"
	json.Unmarshal(response.Body.Bytes(), &body)
	if body["erro"] != expected {
		t.Errorf("wrong response body for param tipo: got %v want %v",
			response.Body.String(), expected)
	}
}

func TestCriarTransacaoRetornaErroQuandoRecebeIdValidoComBodyCampoTipoInvalido(t *testing.T) {

	var jsonStr = []byte(`{"valor": 1000, "tipo": "z", "descricao": "teste"}`)
	req, _ := http.NewRequest("POST", "/clientes/1/transacoes", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeTransacaoRequest(req, "/clientes/{id}/transacoes")

	checkResponseCode(t, http.StatusUnprocessableEntity, response.Code)

	var body map[string]string
	var expected = "o campo tipo deve ser 'd' para débito ou 'c' para crédito"
	json.Unmarshal(response.Body.Bytes(), &body)
	if body["erro"] != expected {
		t.Errorf("wrong response body for param tipo: got %v want %v",
			response.Body.String(), expected)
	}
}

func TestCriarTransacaoRetornaErroQuandoRecebeIdValidoComBodySemCampoDescricao(t *testing.T) {

	var jsonStr = []byte(`{"valor": 1000, "tipo": "c"}`)
	req, _ := http.NewRequest("POST", "/clientes/1/transacoes", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeTransacaoRequest(req, "/clientes/{id}/transacoes")

	checkResponseCode(t, http.StatusUnprocessableEntity, response.Code)

	var body map[string]string
	var expected = "o campo descrição é obrigatório e não pode estar em branco"
	json.Unmarshal(response.Body.Bytes(), &body)
	if body["erro"] != expected {
		t.Errorf("wrong response body for param descrição: got %v want %v",
			response.Body.String(), expected)
	}
}

func TestCriarTransacaoRetornaErroQuandoRecebeIdValidoComBodyECampoDescricaoMaiorQuePermitido(t *testing.T) {

	var jsonStr = []byte(`{"valor": 1000, "tipo": "c", "descricao": "abcdefghijkl"}`)
	req, _ := http.NewRequest("POST", "/clientes/1/transacoes", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeTransacaoRequest(req, "/clientes/{id}/transacoes")

	checkResponseCode(t, http.StatusUnprocessableEntity, response.Code)

	var body map[string]string
	var expected = "o campo descrição não pode conter mais que 10 caracteres"
	json.Unmarshal(response.Body.Bytes(), &body)
	if body["erro"] != expected {
		t.Errorf("wrong response body for param descrição: got %v want %v",
			response.Body.String(), expected)
	}
}

func TestCriarTransacaoRetornaSucessoQuandoRecebeIdValidoComBodyValidoParaTipoCredito(t *testing.T) {

	var jsonStr = []byte(`{"valor": 1000, "tipo": "c", "descricao": "teste"}`)
	req, _ := http.NewRequest("POST", "/clientes/1/transacoes", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeTransacaoRequest(req, "/clientes/{id}/transacoes")

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestCriarTransacaoRetornaSucessoQuandoRecebeIdValidoComBodyValidoParaTipoDebito(t *testing.T) {

	var jsonStr = []byte(`{"valor": 1000, "tipo": "d", "descricao": "teste"}`)
	req, _ := http.NewRequest("POST", "/clientes/1/transacoes", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeTransacaoRequest(req, "/clientes/{id}/transacoes")

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestRecuperarExtratoRetornaErroQuandoRecebeIdInvalido(t *testing.T) {

	req, _ := http.NewRequest("GET", "/clientes/a/extrato", nil)
	response := executeExtratoRequest(req, "/clientes/{id}/extrato")

	checkResponseCode(t, http.StatusUnprocessableEntity, response.Code)
}

func TestRecuperarExtratoRetornaErroQuandoRecebeIdInexistenteNaBase(t *testing.T) {

	req, _ := http.NewRequest("GET", "/clientes/6/extrato", nil)
	response := executeExtratoRequest(req, "/clientes/{id}/extrato")

	checkResponseCode(t, http.StatusNotFound, response.Code)
}

func TestRecuperarExtratoRetornaSucessoQuandoRecebeIdValido(t *testing.T) {

	req, _ := http.NewRequest("GET", "/clientes/1/extrato", nil)
	response := executeExtratoRequest(req, "/clientes/{id}/extrato")

	checkResponseCode(t, http.StatusOK, response.Code)
}

func executeExtratoRequest(req *http.Request, url string) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc(url, controllers.Extrato).Methods("GET")
	router.ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func executeTransacaoRequest(req *http.Request, url string) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc(url, controllers.CriarTransacao).Methods("POST")
	router.ServeHTTP(rr, req)
	return rr
}
