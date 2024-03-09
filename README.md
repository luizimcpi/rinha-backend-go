# rinha-backend-go

## Tecnologias utilizadas
```
go
mux
postgres
nginx
```

## Projeto para atender aos requisitos da rinha de backend 2024
[INSTRUÇÕES DA RINHA](https://github.com/zanfranceschi/rinha-de-backend-2024-q1)

## Como rodar 
```
docker-compose up
```

## Executar os testes
```
go test -> resumido 
go test -v -> verbose mode mostra todos os cenários
```

## Exemplo de request apontando para o Load Balancer (NGINX)

```bash
curl --location 'http://localhost:9999/'
````

## Endpoints
```
use requests.http file with humao.rest-client vscode plugin
```
