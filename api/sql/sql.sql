CREATE DATABASE IF NOT EXISTS devbook;
USE devbook;

DROP TABLE IF EXISTS usuarios;

CREATE TABLE usuarios (
    id INTEGER auto_increment primary key,
    nome VARCHAR(50) NOT NULL,
    nickname VARCHAR(50) NOT NULL unique,
    email VARCHAR(50) NOT NULL unique,
    senha VARCHAR(50) NOT NULL unique,
    criadoEm timestampt default current_timestamp()
) ENGINE=INNODB;