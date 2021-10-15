package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"api/src/autentication"
	"api/src/dbconn"
	"api/src/models"
	"api/src/repositories"
	"api/src/responses"
)

//CriarPublicacao adiciona uma nova publicação no banco de dados
func CriarPublicacao(w http.ResponseWriter, r *http.Request) {
	usuarioID, erro := autentication.ExtrairUsuarioID(r)
	if erro != nil {
		responses.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	corpoRequisicao, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		responses.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var publicacao models.Publicacao
	if erro = json.Unmarshal(corpoRequisicao, &publicacao); erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	publicacao.AutorID = usuarioID

	if erro = publicacao.Preparar(); erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := dbconn.Conectar()
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositories.NovoRepositorioDePublicacao(db)
	publicacao.ID, erro = repositorio.Criar(publicacao)
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusAccepted, publicacao)

}

//BuscarPublicacoes traz as publicações que apareceriam no feed do usuario
func BuscarPublicacoes(w http.ResponseWriter, r *http.Request) {
	
}

//BuscarPublicacao traz uma unica publicação
func BuscarPublicacao(w http.ResponseWriter, r *http.Request) {
	parametros := mux.Vars(r)
	publicacaoID, erro := strconv.ParseUint(parametros["publicacaoId"], 10, 64)
	if erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}
	db, erro := dbconn.Conectar()
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositories.NovoRepositorioDePublicacao(db)
	publicação, erro := repositorio.BuscarPorID(publicacaoID)
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusOK, publicação)
}

//AtualizarPublicacao altera uma publicação
func AtualizarPublicacao(w http.ResponseWriter, r *http.Request) {

}

//DeletarPublicacao exclui uma publicação
func DeletarPublicacao(w http.ResponseWriter, r *http.Request) {

}
