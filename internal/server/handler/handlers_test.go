package handler

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/lookeme/short-url/internal/app/domain/shorten"
	"github.com/lookeme/short-url/internal/configuration"
	"github.com/lookeme/short-url/internal/storage/inmemory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestURLHandlerIndex(t *testing.T) {
	netCfg := configuration.NetworkCfg{
		ServerAddress: ":8080",
		BaseURL:       "http://localhost:8080/",
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
		u, err := url.Parse(string(responseBody))
		require.NoError(t, err)
		key := u.Path[1:len(u.Path)]
		err = res.Body.Close()
		require.NoError(t, err)
		r := httptest.NewRequest(http.MethodGet, "/{id}", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", key)
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
		w = httptest.NewRecorder()
		urlHandler.HandleGet(w, r)
		res = w.Result()
		assert.Equal(t, http.StatusTemporaryRedirect, res.StatusCode)
		assert.Equal(t, requestBody, res.Header.Get("Location"))
		err = res.Body.Close()
		require.NoError(t, err)
	})
}
