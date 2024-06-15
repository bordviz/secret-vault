package models_test

import (
	"errors"
	"testing"
	"vault/internal/models"

	"github.com/stretchr/testify/assert"
)

var (
	ErrName       = errors.New("validation error: [field name is a required]")
	ErrData       = errors.New("validation error: [field data is a required]")
	ErrKeyOrValue = errors.New("key or value length can't be 0")
)

func TestSecretCreateModel(t *testing.T) {
	tests := []struct {
		Name   string
		Model  models.SecretCreateModel
		ErrMsg error
	}{
		{
			Name: "success",
			Model: models.SecretCreateModel{
				Name: "test1",
				Data: map[string]string{
					"test": "test",
				},
			},
			ErrMsg: nil,
		},
		{
			Name: "failed name",
			Model: models.SecretCreateModel{
				Data: map[string]string{
					"test": "test",
				},
			},
			ErrMsg: ErrName,
		},
		{
			Name: "failed name 2",
			Model: models.SecretCreateModel{
				Name: "",
				Data: map[string]string{
					"test": "test",
				},
			},
			ErrMsg: ErrName,
		},
		{
			Name: "failed data",
			Model: models.SecretCreateModel{
				Name: "test3",
			},
			ErrMsg: ErrData,
		},
		{
			Name: "failed data 2",
			Model: models.SecretCreateModel{
				Name: "test4",
				Data: map[string]string{
					"": "test",
				},
			},
			ErrMsg: ErrKeyOrValue,
		},
		{
			Name: "failed data 3",
			Model: models.SecretCreateModel{
				Name: "test5",
				Data: map[string]string{
					"test": "              ",
				},
			},
			ErrMsg: ErrKeyOrValue,
		},
	}

	for _, tt := range tests {
		model := tt.Model
		t.Run(tt.Name, func(t *testing.T) {
			err := model.Validate()
			if err != nil {
				assert.Equal(t, tt.ErrMsg, err)
			}
			if err != nil && tt.ErrMsg == nil {
				t.Fatal("unexpected error")
			}
		})
	}

}
