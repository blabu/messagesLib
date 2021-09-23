package dto

import (
	"crypto/sha256"
	"strings"
)

/*
Набор поддерживаемых команд
протого моста между клиентами
Обработка всех команд происходит в Write методе
*/
const (
	ErrorCOMMAND      uint16 = 1
	PingCOMMAND       uint16 = 2
	RegisterCOMMAND   uint16 = 3
	GenerateCOMMAND   uint16 = 4
	AuthCOMMAND       uint16 = 5
	DataCOMMAND       uint16 = 6
	SaveDataCOMMAND   uint16 = 7
	PropertiesCOMMAND uint16 = 8
)

//CalculateSignature - generate signature
func CalculateSignature(name, salt, token string) [32]byte {
	var cred strings.Builder
	cred.WriteString(name)
	cred.WriteString(salt)
	cred.WriteString(token)
	return sha256.Sum256([]byte(cred.String()))
}
