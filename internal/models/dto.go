package models

import (
	"fmt"
	"strings"
	"time"
	"vault/pkg/validator"
)

type SecretCreateModel struct {
	Name string            `json:"name" validate:"required"`
	Data map[string]string `json:"data" validate:"required"`
}

type SecretCreateDTO struct {
	VaultDTO
	Data []ValueDTO
}

type VaultDTO struct {
	Name string `json:"name" validate:"required"`
}

type ValueDTO struct {
	Key   string `json:"key" validate:"required"`
	Value string `json:"value" validate:"required"`
}

type CreateVaultTokenDTO struct {
	VaultID int           `json:"vault_id" validate:"required"`
	Expires time.Duration `json:"expires" validate:"required"`
}

func (s *SecretCreateModel) Validate() error {
	if err := validator.Validate(s); err != nil {
		return fmt.Errorf("validation error: %s", err)
	}
	for k, v := range s.Data {
		if strings.ReplaceAll(k, " ", "") == "" || strings.ReplaceAll(v, " ", "") == "" {
			return fmt.Errorf("key or value length can't be 0")
		}
	}
	return nil
}

func (s *SecretCreateModel) ConvertToDTO() SecretCreateDTO {
	data := make([]ValueDTO, 0)

	for k, v := range s.Data {
		value := ValueDTO{
			Key:   k,
			Value: v,
		}
		data = append(data, value)
	}

	return SecretCreateDTO{
		VaultDTO: VaultDTO{Name: s.Name},
		Data:     data,
	}
}

func (c *CreateVaultTokenDTO) Validate() error {
	if err := validator.Validate(c); err != nil {
		return fmt.Errorf("validation error: %s", err)
	}
	return nil
}
