package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type SecretModel struct {
	ID   int               `json:"id"`
	Name string            `json:"name"`
	Data map[string]string `json:"data"`
}

type VaultModel struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type TokenModel struct {
	jwt.RegisteredClaims
	ID      int           `json:"vault_id"`
	Expires time.Duration `json:"expires"`
}

func ConvertDTOToSecretModel(vault VaultModel, data []ValueDTO) SecretModel {
	modelData := map[string]string{}
	for _, v := range data {
		modelData[v.Key] = v.Value
	}
	return SecretModel{
		ID:   vault.ID,
		Name: vault.Name,
		Data: modelData,
	}
}
