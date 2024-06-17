package root_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"vault/internal/models"
	"vault/internal/root"
	"vault/internal/root/mocks"
	"vault/pkg/lib/logger/slogdiscard"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateVault(t *testing.T) {
	tests := []struct {
		TestName  string
		Input     string
		Code      int
		Error     string
		MockError error
	}{
		{
			TestName: "success",
			Input:    `{"name": "test", "data": {"some": "data"}}`,
			Code:     201,
		},
		{
			TestName: "failed name",
			Input:    `{"data": {"some": "data"}}`,
			Code:     422,
			Error:    `{"status":"error","detail":"validation error: field name is a required"}`,
		},
		{
			TestName: "failed data",
			Input:    `{"name": "test"}`,
			Code:     422,
			Error:    `{"status":"error","detail":"validation error: field data is a required"}`,
		},
		{
			TestName: "failed decode",
			Input:    `{"name": "test", "data": ""}`,
			Code:     400,
			Error:    `{"status":"error","detail":"failed to decode model"}`,
		},
		{
			TestName: "failed decode",
			Input:    `{"name": "test", "data": {"": ""}}`,
			Code:     422,
			Error:    `{"status":"error","detail":"key or value length can't be 0"}`,
		},
		{
			TestName:  "mock error",
			Input:     `{"name": "test", "data": {"fd": "fasf"}}`,
			Code:      500,
			Error:     `{"status":"error","detail":"failed to save new vault on database"}`,
			MockError: errors.New("failed to save new vault on database"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.TestName, func(t *testing.T) {
			log := slogdiscard.NewDiscardLogger()
			rootDb := mocks.NewRootDB(t)

			if tt.Error == "" || tt.MockError != nil {
				rootDb.On("CreateVault", context.Background(), log, mock.Anything).
					Return(int(1), tt.MockError).
					Once()
			}

			rootHandlers := root.NewRootHandlerClient(rootDb, log, "")
			handler := rootHandlers.CreateVault(context.Background())

			req, err := http.NewRequest(http.MethodPost, "/create", bytes.NewReader([]byte(tt.Input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			body := rr.Body.String()
			body = strings.ReplaceAll(body, "\n", "")

			require.Equal(t, tt.Code, rr.Code)
			if tt.Error != "" {
				require.Equal(t, tt.Error, body)
			}

		})
	}
}

func TestGetVault(t *testing.T) {
	tests := []struct {
		testName  string
		id        int
		code      int
		output    models.SecretModel
		outputStr string
		errorMsg  string
		mockErr   error
	}{
		{
			testName: "success",
			id:       2,
			code:     200,
			output: models.SecretModel{
				ID:   2,
				Name: "Test",
				Data: map[string]string{
					"test": "some sekret key",
				},
			},
			outputStr: `{"id":2,"name":"Test","data":{"test":"some sekret key"}}`,
		},
		{
			testName:  "not found",
			id:        99,
			code:      404,
			errorMsg:  "vault not found",
			outputStr: `{"status":"error","detail":"vault not found"}`,
		},
		{
			testName:  "failed to get vault",
			id:        2,
			code:      500,
			errorMsg:  "failed to get vault",
			mockErr:   errors.New("failed to get vault"),
			outputStr: `{"status":"error","detail":"failed to get vault"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			log := slogdiscard.NewDiscardLogger()
			rootDb := mocks.NewRootDB(t)

			if tt.errorMsg == "" || tt.mockErr != nil {
				rootDb.On("GetVault", context.Background(), log, tt.id).
					Return(
						models.SecretModel{
							ID:   2,
							Name: "Test",
							Data: map[string]string{
								"test": "some sekret key",
							},
						},
						tt.mockErr,
					).Once()
			}
			if tt.errorMsg == "vault not found" {
				rootDb.On("GetVault", context.Background(), log, tt.id).
					Return(models.SecretModel{}, errors.New(tt.errorMsg)).
					Once()
			}

			rootHandlers := root.NewRootHandlerClient(rootDb, log, "")
			handler := rootHandlers.GetVault(context.Background())

			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/get/%v", tt.id), nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			chiCtx := chi.NewRouteContext()
			reqChi := req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
			chiCtx.URLParams.Add("id", fmt.Sprintf("%v", tt.id))

			handler.ServeHTTP(rr, reqChi)

			body := rr.Body.String()
			body = strings.ReplaceAll(body, "\n", "")

			assert.Equal(t, tt.outputStr, body)
		})
	}
}
