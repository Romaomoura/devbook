package repositories

import (
	"database/sql"

	"api/src/models"
)

//Publicacoes representa um repositório de publicações
type Publicacoes struct {
	db *sql.DB
}

//NovoRepositorioDePublicacao cria um repositório de publicação
func NovoRepositorioDePublicacao(db *sql.DB) *Publicacoes {
	return &Publicacoes{db}
}

//Criar insere uma publicação no banco de dados
func (repositories Publicacoes) Criar(publicacao models.Publicacao) (uint64, error) {
	statement, erro := repositories.db.Prepare(
		"insert into publicacoes (titulo, conteudo, autor_id) values(?, ?, ?)",
	)
	if erro != nil {
		return 0, erro
	}
	defer statement.Close()

	resultado, erro := statement.Exec(publicacao.Titulo, publicacao.Conteudo, publicacao.AutorID)
	if erro != nil {
		return 0, erro
	}

	ultimoIDInserido, erro := resultado.LastInsertId()
	if erro != nil {
		return 0, erro
	}

	return uint64(ultimoIDInserido), nil
}
