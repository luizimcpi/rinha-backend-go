# rinha-backend-go

## Tecnologias utilizadas
```
go
mux
```

## Projeto para atender aos requisitos da rinha de backend 2024
[INSTRUÇÕES DA RINHA](https://github.com/zanfranceschi/rinha-de-backend-2024-q1)

## Como rodar 
```
docker-compose up
```

## Exemplo de request apontando para o Load Balancer (NGINX)

```bash
curl --location 'http://localhost:9999/'
````

## Mysql docker instructions
```
1- docker exec -it rinha-backend-go_database_1 bash

2- mysql -u user -p
123456

3- use rinhabank;

4- select * from clientes;
```

## Endpoints


### Create transaction
```
curl --location --request POST 'localhost:9999/clientes/1/transacoes' \
--header 'Content-Type: application/json' \
--data-raw '{
    "valor": 1000,
    "tipo" : "c",
    "descricao" : "descricao"
}'
````