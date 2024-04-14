package handler

import (
	_ "encoding/json"
	"github.com/lookeme/short-url/internal/app/domain/shorten"
	"github.com/lookeme/short-url/internal/configuration"
	_ "github.com/lookeme/short-url/internal/configuration"
	"github.com/lookeme/short-url/internal/storage/inmemory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestURLHandlerIndex(t *testing.T) {
	netCfg := configuration.NetworkCfg{
		StartAddress: ":8080",
		BaseAddress:  "http://localhost:8080",
	}
	cfg := configuration.Config{
		Network: &netCfg,
	}
	storage := inmemory.NewStorage()
	urlService := shorten.NewURLService(storage)
	urlHandler := NewURLHandler(urlService, &cfg)
	requestBody := "https://practicum.yandex.ru/"
	bodyReader := strings.NewReader(requestBody)
	t.Run("handler test #1", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", bodyReader)
		w := httptest.NewRecorder()
		urlHandler.HandlePOST(w, req)
		res := w.Result()
		assert.Equal(t, http.StatusCreated, res.StatusCode)
		assert.Equal(t, "text/plain", res.Header.Get("Content-Type"))
		responseBody, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		assert.True(t, len(string(responseBody)) > 0)
		url := strings.Split(string(responseBody), "/")
		key := url[len(url)-1]
		err = res.Body.Close()
		require.NoError(t, err)
		req = httptest.NewRequest(http.MethodGet, "/"+key, nil)
		w = httptest.NewRecorder()
		urlHandler.HandleGet(w, req)
		res = w.Result()
		assert.Equal(t, http.StatusTemporaryRedirect, res.StatusCode)
		assert.Equal(t, requestBody, res.Header.Get("Location"))
		err = res.Body.Close()
		require.NoError(t, err)
	})
}
