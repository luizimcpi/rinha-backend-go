CREATE TABLE public.clientes
(
  id INT PRIMARY KEY,
  limite INT NOT NULL,
  saldo INT NOT NULL DEFAULT 0
);

insert into public.clientes (id, limite)
values
(1, 100000),
(2, 80000),
(3, 1000000),
(4, 10000000),
(5, 500000);

CREATE TABLE public.transacoes
(
  id SERIAL PRIMARY KEY,
  valor INT not null,
  tipo CHAR(1) not null,
  descricao character varying(11) not null,
  realizada_em TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  cliente_id INTEGER REFERENCES public.clientes (id) NOT NULL
);

CREATE INDEX transacao_idx ON transacoes USING btree (cliente_id, realizada_em);