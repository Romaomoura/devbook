package seguranca

import (
	"golang.org/x/crypto/bcrypt"
)

//Hash recebe uma senha string e converte em hash
func Hash(senha string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(senha), bcrypt.DefaultCost)
}

//VerrificaSenha compara uma senha e um hash e retorna caso sejam iguais
func VerrificaSenha(senhaHash, senhaString string) error {
	//fmt.Print("Hash>>>  ", senhaHash, "  Senha String>>>  ", senhaString)
	return bcrypt.CompareHashAndPassword([]byte(senhaHash), []byte(senhaString))
}
