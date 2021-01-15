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
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

//CriarUsuario insere um usuario no banco de dados
func CriarUsuario(w http.ResponseWriter, r *http.Request) {
	corpoRequisicao, erro := ioutil.ReadAll(r.Body)
	//erro = errors.New("Deu merda!")
	if erro != nil {
		responses.Erro(w, http.StatusUnprocessableEntity, erro)
		return
	}

	var usuario models.Usuario
	if erro := json.Unmarshal(corpoRequisicao, &usuario); erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if erro = usuario.Preparar("cadastro"); erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := dbconn.Conectar()
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, errors.New("Não foi possivel conectar ao banco de dados"))
		return
	}
	defer db.Close()

	repositorio := repositories.NovoRepositorioDeUsuarios(db)
	usuarioID, erro := repositorio.Criar(usuario)
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusCreated, usuarioID)

}

//BuscarUsuarios busca todos os usuarios no banco de dados
func BuscarUsuarios(w http.ResponseWriter, r *http.Request) {

	nomeOuNickname := strings.ToLower(r.URL.Query().Get("q"))

	db, erro := dbconn.Conectar()
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, errors.New("Não foi possivel conectar ao banco de dados"))
		return
	}
	defer db.Close()

	repositorio := repositories.NovoRepositorioDeUsuarios(db)
	usuarios, erro := repositorio.Buscar(nomeOuNickname)
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusOK, usuarios)

}

//BuscarUsuario busca um usuario no banco de dados
func BuscarUsuario(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	usuarioID, erro := strconv.ParseUint(params["usuarioId"], 10, 64)
	if erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := dbconn.Conectar()
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, errors.New("Não foi possivel conectar ao banco de dados"))
		return
	}
	defer db.Close()

	repositorio := repositories.NovoRepositorioDeUsuarios(db)
	usuario, erro := repositorio.BuscarPorID(usuarioID)
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	responses.JSON(w, http.StatusOK, usuario)
}

//AtualizarUsuario atualiza um usuario no banco de dados
func AtualizarUsuario(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(params["usuarioId"], 10, 64)
	if erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	usuarioIDToken, erro := autentication.ExtrairUsuarioID(r)
	if erro != nil {
		responses.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	if usuarioID != usuarioIDToken {
		responses.Erro(w, http.StatusForbidden, errors.New("Não é possivel atualizar um usuário que não seja o seu"))
		return
	}

	corpoRequisicao, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		responses.Erro(w, http.StatusUnprocessableEntity, erro)
	}

	var usuario models.Usuario
	if erro = json.Unmarshal(corpoRequisicao, &usuario); erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}
	if erro = usuario.Preparar("edicao"); erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}
	db, erro := dbconn.Conectar()
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, errors.New("Não foi possivel conectar ao banco de dados"))
		return
	}
	defer db.Close()

	repositorio := repositories.NovoRepositorioDeUsuarios(db)
	if erro = repositorio.Atualizar(usuarioID, usuario); erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	responses.JSON(w, http.StatusNoContent, nil)
}

//DeletarUsuario exclui um usuario no banco de dados
func DeletarUsuario(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(params["usuarioId"], 10, 64)
	if erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}
	usuarioIDToken, erro := autentication.ExtrairUsuarioID(r)
	if erro != nil {
		responses.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	if usuarioID != usuarioIDToken {
		responses.Erro(w, http.StatusForbidden, errors.New("Não é possivel deletar um usuário que não seja o seu"))
		return
	}

	db, erro := dbconn.Conectar()
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	repositorio := repositories.NovoRepositorioDeUsuarios(db)
	if erro = repositorio.Deletar(usuarioID); erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}
	responses.JSON(w, http.StatusNoContent, nil)
}

//SeguirUsuario permite que um usuário siga outro
func SeguirUsuario(w http.ResponseWriter, r *http.Request) {
	seguidorID, erro := autentication.ExtrairUsuarioID(r)
	if erro != nil {
		responses.Erro(w, http.StatusUnauthorized, erro)
		return
	}
	params := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(params["usuarioId"], 10, 64)
	if erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if seguidorID == usuarioID {
		responses.Erro(w, http.StatusBadRequest, errors.New("Caramba, você é muito seu fã, mas dessa vez não"))
		return
	}

	db, erro := dbconn.Conectar()
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, errors.New("Não foi possivel conectar ao banco de dados"))
		return
	}
	defer db.Close()

	repositorio := repositories.NovoRepositorioDeUsuarios(db)
	if erro = repositorio.Seguir(usuarioID, seguidorID); erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)
}

//DeixarDeSeguirUsuario permite que um usuário deixe de seguir outro
func DeixarDeSeguirUsuario(w http.ResponseWriter, r *http.Request) {
	seguidorID, erro := autentication.ExtrairUsuarioID(r)
	if erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}
	params := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(params["usuarioId"], 10, 64)
	if erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if seguidorID == usuarioID {
		responses.Erro(w, http.StatusBadRequest, errors.New("Você está condenado a sua propria companhia, ha ha ha"))
		return
	}

	db, erro := dbconn.Conectar()
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, errors.New("Não foi possivel conectar ao banco de dados"))
		return
	}
	defer db.Close()

	repositorio := repositories.NovoRepositorioDeUsuarios(db)
	if erro = repositorio.DeixarDeSeguir(usuarioID, seguidorID); erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)

}

//BuscarSeguidores traz todos os seguidores de um usuário
func BuscarSeguidores(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(params["usuarioId"], 10, 64)
	if erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := dbconn.Conectar()
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, errors.New("Não foi possivel conectar ao banco de dados"))
		return
	}
	defer db.Close()

	repositorio := repositories.NovoRepositorioDeUsuarios(db)
	seguidores, erro := repositorio.BuscarSeguidores(usuarioID)
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusOK, seguidores)
}

//BuscarSeguindo traz todos os usuários que um usuário está seguindo
func BuscarSeguindo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(params["usuarioId"], 10, 64)
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

	repositorio := repositories.NovoRepositorioDeUsuarios(db)
	usuarios, erro := repositorio.BuscarSeguindo(usuarioID)
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, erro)
		return
	}

	responses.JSON(w, http.StatusOK, usuarios)
}

//AtualizarSenha permite atualizar a senha de um usuário
func AtualizarSenha(w http.ResponseWriter, r *http.Request) {
	usuarioIDToken, erro := autentication.ExtrairUsuarioID(r)
	if erro != nil {
		responses.Erro(w, http.StatusUnauthorized, erro)
		return
	}

	params := mux.Vars(r)
	usuarioID, erro := strconv.ParseUint(params["usuarioId"], 10, 64)
	if erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if usuarioIDToken != usuarioID {
		responses.Erro(w, http.StatusForbidden, errors.New("Ai você já quer demais, atualizar senha de outro não pode"))
		return
	}

	corpoRequisicao, erro := ioutil.ReadAll(r.Body)

	var senha models.Senha
	if erro = json.Unmarshal(corpoRequisicao, &senha); erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := dbconn.Conectar()
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, errors.New("Não foi possivel conectar ao banco de dados"))
		return
	}

	defer db.Close()

	repositorio := repositories.NovoRepositorioDeUsuarios(db)
	senhadb, erro := repositorio.BuscarSenha(usuarioID)
	if erro != nil {
		responses.Erro(w, http.StatusInternalServerError, errors.New("Não foi possivel buscar a senha solicitada"))
		return
	}

	if erro = seguranca.VerificaSenha(senhadb, senha.Atual); erro != nil {
		responses.Erro(w, http.StatusUnauthorized, errors.New("Senha não são as mesmas"))
		return
	}

	senhaHash, erro := seguranca.Hash(senha.Nova)
	if erro != nil {
		responses.Erro(w, http.StatusBadRequest, erro)
		return
	}

	if erro = repositorio.AtualizarSenha(usuarioID, string(senhaHash)); erro != nil {
		responses.Erro(w, http.StatusInternalServerError, errors.New("Não foi possivel atualizar a senha"))
		return
	}

	responses.JSON(w, http.StatusNoContent, nil)

}
