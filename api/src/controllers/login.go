package controllers

import (
	"api/src/dbconn"
	"api/src/models"
	"api/src/repositories"
	"api/src/responses"
	"api/src/seguranca"
	"encoding/json"
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
	println()
	println(usuario.Email)

	db, erro := dbconn.Conectar()
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositories.NovoRepositorioDeUsuarios(db)
	usuariodb, erro := repositorio.BuscarPorEmail(usuario.Email)
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	//Compara as senhas
	if erro = seguranca.VerrificaSenha(usuariodb.Senha, usuario.Senha); erro != nil {
		responses.Erro(w, http.StatusUnauthorized, erro)
		return
	}
	w.Write([]byte("Parabéns Você está logado!"))
}
