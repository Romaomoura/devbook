package models

import (
	"api/src/seguranca"
	"errors"
	"strings"
	"time"

	"github.com/badoux/checkmail" // Validação de email
)

//Usuario representa um usuario na rede social
type Usuario struct {
	ID       uint64    `json:"id,omitempty"`
	Nome     string    `json:"nome,omitempty"`
	Nickname string    `json:"nickname,omitempty"`
	Email    string    `json:"email,omitempty"`
	Senha    string    `json:"senha,omitempty"`
	CriadoEm time.Time `json:"criadoEm,omitempty"`
}

//Preparar irá chamar os métodos para validar e formatar o usuário recebido
func (usuario *Usuario) Preparar(etapa string) error {
	if erro := usuario.validar(etapa); erro != nil {
		return erro
	}
	if erro := usuario.formatar(etapa); erro != nil {
		return erro
	}
	return nil
}

func (usuario *Usuario) validar(etapa string) error {

	if usuario.Nome == "" {
		return errors.New("O nome é obrigatório e não pode ser em branco")
	}

	if usuario.Nickname == "" {
		return errors.New("O nickname é obrigatório e não pode ser em branco")
	}

	if usuario.Email == "" {
		return errors.New("O email é obrigatório e não pode ser em branco")
	}

	if erro := checkmail.ValidateFormat(usuario.Email); erro != nil {
		return errors.New("O email inserido é inválido")
	}

	if etapa == "cadastro" && usuario.Senha == "" {
		return errors.New("A senha é obrigatório e não pode ser em branco")
	}

	return nil
}

func (usuario *Usuario) formatar(etapa string) error {
	usuario.Nome = strings.TrimSpace(usuario.Nome)
	usuario.Nickname = strings.TrimSpace(usuario.Nickname)
	usuario.Email = strings.TrimSpace(usuario.Email)

	if etapa == "cadastro" {
		senhaComHash, erro := seguranca.Hash(usuario.Senha)
		if erro != nil {
			return erro
		}
		usuario.Senha = string(senhaComHash)
	}
	return nil
}
