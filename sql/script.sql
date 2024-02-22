CREATE TABLE clientes(
    id SERIAL PRIMARY KEY,
    limite INT NOT NULL,
    saldo INT NOT NULL DEFAULT 0
);

CREATE UNLOGGED TABLE transacoes(
    id SERIAL primary key,
    valor INTEGER NOT NULL,
    tipo CHAR(1) NOT NULL,
    descricao VARCHAR(12) NOT NULL,
    realizada_em TIMESTAMP NOT NULL DEFAULT NOW(),
    cliente_id INTEGER NOT NULL,
    CONSTRAINT fk_transacoes_cliente_id FOREIGN KEY (cliente_id) REFERENCES clientes (id)
);

DO
  $$ BEGIN INSERT INTO clientes (limite) 
VALUES 
  (100000), 
  (80000), 
  (1000000), 
  (10000000), 
  (500000);
END;
$$