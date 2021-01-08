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
func (repositorio Usuarios) Criar(usuario models.Usuario) (uint64, error) {
	statement, erro := repositorio.db.Prepare(
		"INSERT INTO usuarios (nome, nickname, email, senha) VALUES(?,?,?,?)",
	)
	if erro != nil {
		fmt.Println("Erro ao preparar a statement, motivo: ", erro)
		return 0, erro
	}
	defer statement.Close()
	resultado, erro := statement.Exec(usuario.Nome, usuario.Nickname, usuario.Email, usuario.Senha)
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
	nomeOuNickname = fmt.Sprintf("%%%s%%", nomeOuNickname) //%nomeOuNickname

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
		var usuario models.Usuario

		if erro = linhas.Scan(
			&usuario.ID,
			&usuario.Nome,
			&usuario.Nickname,
			&usuario.Email,
			&usuario.CriadoEm,
		); erro != nil {
			return nil, erro
		}
		usuarios = append(usuarios, usuario)
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

	var usuario models.Usuario

	if linhas.Next() {
		if erro = linhas.Scan(
			&usuario.ID,
			&usuario.Nome,
			&usuario.Nickname,
			&usuario.Email,
			&usuario.CriadoEm,
		); erro != nil {
			return models.Usuario{}, erro
		}
	}
	return usuario, nil
}

//Atualizar altera as informações de um usuário no banco de dados
func (repositorio Usuarios) Atualizar(ID uint64, usuario models.Usuario) error {
	statement, erro := repositorio.db.Prepare(
		"UPDATE usuarios SET nome = ?, nickname = ?, email = ? WHERE id = ?",
	)
	if erro != nil {
		return erro
	}
	defer statement.Close()

	if _, erro = statement.Exec(usuario.Nome, usuario.Nickname, usuario.Email, ID); erro != nil {
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

	var usuario models.Usuario

	if linha.Next() {
		if erro = linha.Scan(&usuario.ID, &usuario.Senha); erro != nil {
			return models.Usuario{}, erro
		}
	}
	return usuario, nil

}
