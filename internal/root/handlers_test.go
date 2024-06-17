package root_test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"vault/internal/root"
	"vault/internal/root/mocks"
	"vault/pkg/lib/logger/slogdiscard"

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
