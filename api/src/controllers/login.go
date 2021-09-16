package controllers

import (
	"api/src/autentication"
	"api/src/dbconn"
	"api/src/models"
	"api/src/repositories"
	"api/src/responses"
	"api/src/seguranca"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

//Login é responsável por autenticar um usuário na API
func Login(w http.ResponseWriter, r *http.Request) {
	corpoRequisicao, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		responses.Erro(w, http.StatusUnprocessableEntity, erro)
	}

	var usuario models.Usuario
	if erro = json.Unmarshal(corpoRequisicao, &usuario); erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := dbconn.Conectar()
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	//fmt.Println("ID do usuario depois da conexão >>>>  ", usuario.ID)
	repositorio := repositories.NovoRepositorioDeUsuarios(db)
	usuariodb, erro := repositorio.BuscarPorEmail(usuario.Email)
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	//fmt.Println("ID banco para compara>>>>", usuariodb.ID)
	//Compara as senhas
	if erro = seguranca.VerificaSenha(usuariodb.Senha, usuario.Senha); erro != nil {
		responses.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	token, erro := autentication.CriarToken(usuariodb.ID)
	fmt.Println("ID para criar token>>>>", usuariodb.ID)
	if erro == nil && usuariodb.ID == 0 {
		responses.Erro(w, http.StatusInternalServerError, errors.New("Id nulo ou zero"))
		return
	}
	w.Write([]byte(token))
}
