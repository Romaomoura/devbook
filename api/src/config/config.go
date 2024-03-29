package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	//StringConexaoDB é a string de conexão com o banco
	StringConexaoDB = ""

	//Porta é o local onde a API estará rodando
	Porta = 0

	//Secretkey é a chave usada para assinar o token
	Secretkey []byte
)

//Carregar vai inicializar as variaveis de ambiente
func Carregar() {
	var erro error

	if erro = godotenv.Load(); erro != nil {
		log.Fatal(erro)
	}

	Porta, erro = strconv.Atoi(os.Getenv("API_PORT"))
	if erro != nil {
		Porta = 9000
	}

	StringConexaoDB = fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DB_USUARIO"),
		os.Getenv("DB_SENHA"),
		os.Getenv("DB_NOME"),
	)
	fmt.Println(StringConexaoDB)
	Secretkey = []byte(os.Getenv("SECRET_KEY"))
}
