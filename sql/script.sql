DROP TABLE IF EXISTS clientes;

CREATE TABLE clientes(
    id int auto_increment primary key,
    limite int not null,
    saldo_inicial int not null,
    data_criacao timestamp default current_timestamp()
) ENGINE=INNODB;

insert into clientes (limite, saldo_inicial)
values
(100000, 0), 
(80000, 0),
(1000000, 0),
(10000000, 0),
(500000, 0) 