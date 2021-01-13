CREATE DATABASE IF NOT EXISTS devbook;
USE devbook;

DROP TABLE IF EXISTS usuarios;
DROP TABLE IF EXISTS seguidores;

CREATE TABLE usuarios (
    id INTEGER auto_increment primary key,
    nome VARCHAR(50) NOT NULL,
    nickname VARCHAR(50) NOT NULL unique,
    email VARCHAR(50) NOT NULL unique,
    senha VARCHAR(50) NOT NULL unique,
    criadoEm timestampt default current_timestamp
) ENGINE=INNODB;

CREATE Table seguidores(
    usuario_id int not null,
    FOREIGN KEY (usuario_id) 
    REFERENCES usuarios(id)
    ON DELETE CASCADE,

    seguidor_id int not null,
    FOREIGN KEY (seguidor_id) 
    REFERENCES usuarios(id)
    ON DELETE CASCADE,

    primary key (usuario_id, seguidor_id)
)ENGINE=INNODB;