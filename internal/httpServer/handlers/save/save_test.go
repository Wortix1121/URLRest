package save_test

import (
	"api/internal/httpServer/handlers/save"
	"api/internal/httpServer/handlers/save/mocks"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestSaveHandler(t *testing.T) {
	cases := []struct {
		name      string
		alias     string
		url       string
		respError string
		mockError error
	}{
		{
			name:  "Success",
			alias: "test_alias",
			url:   "https://google.com",
		},
		// {
		// 	name:  "Empty alias",
		// 	alias: "",
		// 	url:   "https://google.com",
		// },
		// {
		// 	name:      "Empty URL",
		// 	url:       "",
		// 	alias:     "some_alias",
		// 	respError: "field URL is a required field",
		// },
		// {
		// 	name:      "Invalid URL",
		// 	url:       "some invalid URL",
		// 	alias:     "some_alias",
		// 	respError: "field URL is not a valid URL",
		// },
		// {
		// 	name:      "SaveURL Error",
		// 	alias:     "test_alias",
		// 	url:       "https://google.com",
		// 	respError: "failed to add url",
		// 	mockError: errors.New("unexpected error"),
		// },
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			// t.Parallel()

			urlSaverMock := mocks.NewURLSaver(t)

			if tc.respError == "" || tc.mockError != nil {
				urlSaverMock.On("SaveUrl", tc.url, mock.AnythingOfType("string")).
					Return(int64(1), tc.mockError).
					Once()
			}

			var log *zap.Logger

			handler := save.SaveHand(log, urlSaverMock)

			input := fmt.Sprintf(`{"url": "%s", "alias": "%s"}`, tc.url, tc.alias)

			req, err := http.NewRequest(http.MethodPost, "/url", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			// req := httptest.NewRequest("post", "url", bytes.NewReader([]byte(input)))
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()

			var resp save.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)

			// TODO: add more checks
		})
	}
}
