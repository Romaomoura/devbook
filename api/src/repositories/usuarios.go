package repositories

import (
	"api/src/models"
	"database/sql"
	"fmt"
)

//Usuarios representa um repositorio de usuários
type Usuarios struct {
	db *sql.DB
}

//NovoRepositorioDeUsuarios cria um repositorio de usuários
func NovoRepositorioDeUsuarios(db *sql.DB) *Usuarios {
	return &Usuarios{db}
}

//Criar cria um usuário no banco de dados
func (repositorio Usuarios) Criar(seguindo models.Usuario) (uint64, error) {
	statement, erro := repositorio.db.Prepare(
		"insert into usuarios (nome, nickname, email, senha) values(?, ?, ?, ?)")
	if erro != nil {
		fmt.Println("Erro ao preparar a statement, motivo: ", erro)
		return 0, erro
	}

	defer statement.Close()

	fmt.Println(seguindo)
	resultado, erro := statement.Exec(seguindo.Nome, seguindo.Nickname, seguindo.Email, seguindo.Senha)
	if erro != nil {
		fmt.Println("Erro ao inserir no banco, motivo: ", erro)
		return 0, nil
	}
	ultimoIDInserido, erro := resultado.LastInsertId()
	if erro != nil {
		return 0, nil
	}

	return uint64(ultimoIDInserido), nil
}

//Buscar retorna todos os usuários que atendam ao filtro nome ou nickname
func (repositorio Usuarios) Buscar(nomeOuNickname string) ([]models.Usuario, error) {
	nomeOuNickname = fmt.Sprintf("%%%s%%", nomeOuNickname) //%nomeOuNickname%

	linhas, erro := repositorio.db.Query(
		"SELECT id, nome, nickname, email, criadoEm FROM usuarios WHERE nome LIKE ? OR nickname LIKE ?",
		nomeOuNickname, nomeOuNickname,
	)

	if erro != nil {
		return nil, erro
	}

	defer linhas.Close()

	var usuarios []models.Usuario

	for linhas.Next() {
		var seguindo models.Usuario

		if erro = linhas.Scan(
			&seguindo.ID,
			&seguindo.Nome,
			&seguindo.Nickname,
			&seguindo.Email,
			&seguindo.CriadoEm,
		); erro != nil {
			return nil, erro
		}
		usuarios = append(usuarios, seguindo)
	}
	return usuarios, nil

}

//BuscarPorID retorna um usuário do banco de dados
func (repositorio Usuarios) BuscarPorID(ID uint64) (models.Usuario, error) {
	linhas, erro := repositorio.db.Query(
		"SELECT id, nome, nickname, email, criadoEm FROM usuarios WHERE id = ?",
		ID,
	)
	if erro != nil {
		return models.Usuario{}, nil
	}
	defer linhas.Close()

	var seguindo models.Usuario

	if linhas.Next() {
		if erro = linhas.Scan(
			&seguindo.ID,
			&seguindo.Nome,
			&seguindo.Nickname,
			&seguindo.Email,
			&seguindo.CriadoEm,
		); erro != nil {
			return models.Usuario{}, erro
		}
	}
	return seguindo, nil
}

//Atualizar altera as informações de um usuário no banco de dados
func (repositorio Usuarios) Atualizar(ID uint64, seguindo models.Usuario) error {
	statement, erro := repositorio.db.Prepare(
		"UPDATE usuarios SET nome = ?, nickname = ?, email = ? WHERE id = ?",
	)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(seguindo.Nome, seguindo.Nickname, seguindo.Email, ID); erro != nil {
		return erro
	}

	return nil
}

//Deletar exclui as informações de um usuário no banco de dados
func (repositorio Usuarios) Deletar(ID uint64) error {
	statement, erro := repositorio.db.Prepare(
		"DELETE FROM usuarios WHERE id = ?",
	)
	if erro != nil {
		return erro
	}
	defer statement.Close()
	if _, erro = statement.Exec(ID); erro != nil {
		return erro
	}

	return nil
}

//BuscarPorEmail busca um usuar no banco de dados pelo email e retorna o seu ID e SENHA com hash.
func (repositorio Usuarios) BuscarPorEmail(email string) (models.Usuario, error) {
	linha, erro := repositorio.db.Query("SELECT id, senha FROM usuarios WHERE email = ?", email)
	if erro != nil {
		return models.Usuario{}, erro
	}
	defer linha.Close()

	var seguindo models.Usuario

	if linha.Next() {
		if erro = linha.Scan(&seguindo.ID, &seguindo.Senha); erro != nil {
			return models.Usuario{}, erro
		}
	}
	//fmt.Println(">Busca>>>>>", seguindo.ID)
	return seguindo, nil

}

//Seguir permite que um usuário siga outro
func (repositorio Usuarios) Seguir(usuarioID, seguindoID uint64) error {
	statement, erro := repositorio.db.Prepare(
		"insert ignore into seguindoes (usuario_id, seguindo_id) values (?, ?)",
	)
	if erro != nil {
		return erro
	}
	defer statement.Close()
	if _, erro = statement.Exec(usuarioID, seguindoID); erro != nil {
		return erro
	}

	return nil
}

//DeixarDeSeguir permite que um usuário deixe de seguir outro
func (repositorio Usuarios) DeixarDeSeguir(usuarioID, seguindoID uint64) error {
	statement, erro := repositorio.db.Prepare(
		"delete from seguindoes where usuario_id = ? and seguindo_id = ?",
	)
	if erro != nil {
		return erro
	}
	defer statement.Close()
	if _, erro = statement.Exec(usuarioID, seguindoID); erro != nil {
		return erro
	}

	return nil
}

//BuscarSeguidores retorna todos os seguindoes de um usuário
func (repositorio Usuarios) BuscarSeguidores(usuarioID uint64) ([]models.Usuario, error) {
	linhas, erro := repositorio.db.Query(`
		select u.id, u.nome, u.nickname, u.email, u.criadoEm from usuarios u 
		inner join seguidores s on u.id = s.seguidor_id 
		where s.usuario_id = ? `, usuarioID)

	if erro != nil {
		return nil, erro
	}

	defer linhas.Close()

	var seguidores []models.Usuario

	for linhas.Next() {
		var seguindo models.Usuario

		if erro = linhas.Scan(
			&seguindo.ID,
			&seguindo.Nome,
			&seguindo.Nickname,
			&seguindo.Email,
			&seguindo.CriadoEm,
		); erro != nil {
			return nil, erro
		}
		seguidores = append(seguidores, seguindo)
	}
	return seguidores, nil

}

//BuscarSeguindo retorna todos os usuários que um determinado usuário segue
func (repositorio Usuarios) BuscarSeguindo(usuarioID uint64) ([]models.Usuario, error) {
	linhas, erro := repositorio.db.Query(`
		select u.id, u.nome, u.nickname, u.email, u.criadoEm from usuarios u 
		inner join seguidores s on u.id = s.usuario_id 
		where s.seguidor_id = ? `, usuarioID)

	if erro != nil {
		return nil, erro
	}

	defer linhas.Close()

	var usuarios []models.Usuario

	for linhas.Next() {
		var seguindo models.Usuario

		if erro = linhas.Scan(
			&seguindo.ID,
			&seguindo.Nome,
			&seguindo.Nickname,
			&seguindo.Email,
			&seguindo.CriadoEm,
		); erro != nil {
			return nil, erro
		}
		usuarios = append(usuarios, seguindo)
	}
	return usuarios, nil

}

//BuscarSenha retorna a senha de um usuario salva no banco de dados por ID
func (repositorio Usuarios) BuscarSenha(usuarioID uint64) (string, error) {
	linha, erro := repositorio.db.Query("SELECT senha FROM usuarios WHERE id = ?", usuarioID)
	if erro != nil {
		return "Ops, não foi dessa vez!!!", erro
	}
	defer linha.Close()

	var usuario models.Usuario

	if linha.Next() {
		if erro = linha.Scan(&usuario.Senha); erro != nil {
			return "Ops, não foi dessa vez!!!", erro
		}
	}

	return usuario.Senha, nil
}

//AtualizarSenha atualiza a senha de um usuário
func (repositorio Usuarios) AtualizarSenha(usuarioID uint64, senha string) error {
	statement, erro := repositorio.db.Prepare("update usuarios set senha = ? where id = ?")
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(senha, usuarioID); erro != nil {
		return erro
	}

	return nil
}
