CREATE DATABASE IF NOT EXISTS devbook;
USE devbook;

DROP TABLE IF EXISTS publicacoes;

DROP TABLE IF EXISTS usuarios;
DROP TABLE IF EXISTS seguidores;

CREATE TABLE usuarios (
    id INTEGER auto_increment primary key,
    nome VARCHAR(50) NOT NULL,
    nickname VARCHAR(50) NOT NULL unique,
    email VARCHAR(50) NOT NULL unique,
    senha VARCHAR(50) NOT NULL unique,
    criadoEm timestamp default current_timestamp
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

    primary key (usuario_idseguidor_id)
)ENGINE=INNODB;


desc seguidores

desc usuarios



CREATE TABLE publicacoes (
    id int auto_increment primary key,
    titulo varchar(50) not null,
    conteudo varchar(300) not null,

    autor_id int not null,
    FOREIGN KEY (autor_id) 
    REFERENCES usuarios(id)
    ON DELETE CASCADE,
    curtidas int default 0,
    criadoEm timestamp default current_timestamp,
)ENGINE=INNODB;