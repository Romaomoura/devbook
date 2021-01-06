package dbconn

import (
	"api/src/config"
	"database/sql"

	_ "github.com/go-sql-driver/mysql" //Driver de conexão
)

//Conectar abre a conexão com o banco de dados e a retornar
func Conectar() (*sql.DB, error) {
	db, erro := sql.Open("mysql", config.StringConexaoDB)
	if erro != nil {
		return nil, erro
	}
	if erro := db.Ping(); erro != nil {
		return nil, erro
	}

	return db, nil
}
