package controllers

import (
	"api/src/dbconn"
	"api/src/models"
	"api/src/repositories"
	"api/src/responses"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

//CriarUsuario insere um usuario no banco de dados
func CriarUsuario(w http.ResponseWriter, r *http.Request) {
	corpoRequest, erro := ioutil.ReadAll(r.Body)
	//erro = errors.New("Deu merda!")
	if erro != nil {
		responses.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var usuario models.Usuario
	if erro := json.Unmarshal(corpoRequest, &usuario); erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if erro = usuario.Preparar(); erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := dbconn.Conectar()
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositories.NovoRepositorioDeUsuarios(db)
	usuario.ID, erro = repositorio.Criar(usuario)
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusCreated, usuario)

}

//BuscarUsuarios busca todos os usuarios no banco de dados
func BuscarUsuarios(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("Buscando todos os usuários!"))

}

//BuscarUsuario busca um usuario no banco de dados
func BuscarUsuario(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("Buscando um usuário!"))

}

//AtualizarUsuario atualiza um usuario no banco de dados
func AtualizarUsuario(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("Atualizando usuário!"))

}

//DeletarUsuario deleta um usuario no banco de dados
func DeletarUsuario(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("Deletando usuário!"))

}
