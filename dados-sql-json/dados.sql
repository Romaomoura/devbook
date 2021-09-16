insert into usuarios(nome, nickname, email, senha)
values
("João","João", "João@gmail.com","$2a$10$0iI2ZFBHSfrsVOPkA9DU2ei1EUpRKvcyZGeZxhz4vZeRaVlaBTJmW" ),
("Maria","Maria", "Maria@gmail.com","$2a$10$0iI2ZFBHSfrsVOPkA9DU2ei1EUpRKvcyZGeZxhz4vZeRaVlaBTJmW" ),
("Pedro","Pedro", "Pedro@gmail.com","$2a$10$0iI2ZFBHSfrsVOPkA9DU2ei1EUpRKvcyZGeZxhz4vZeRaVlaBTJmW" ),
("Jose","Jose", "Jose@gmail.com","$2a$10$0iI2ZFBHSfrsVOPkA9DU2ei1EUpRKvcyZGeZxhz4vZeRaVlaBTJmW" ),
("Antonio","Antonio", "Antonio@gmail.com","$2a$10$0iI2ZFBHSfrsVOPkA9DU2ei1EUpRKvcyZGeZxhz4vZeRaVlaBTJmW" );

INSERT INTO seguidores(usuario_id, seguidor_id)
VALUES 
(1, 3), (1, 24), (3, 20), (20, 1), (20, 3),
(22, 23), (23, 30), (31, 25), (3, 25), (24, 1);


select * from usuarios;

select * from seguidores;