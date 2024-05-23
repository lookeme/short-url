package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/lookeme/short-url/internal/app/domain/user"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/lookeme/short-url/internal/app/domain/shorten"
	"github.com/lookeme/short-url/internal/configuration"
	"github.com/lookeme/short-url/internal/logger"
	"github.com/lookeme/short-url/internal/models"
	"github.com/lookeme/short-url/internal/storage/inmemory"
)

func TestURLHandlerIndex(t *testing.T) {
	netCfg := configuration.NetworkCfg{
		ServerAddress: ":8080",
		BaseURL:       "http://localhost:8080/",
	}
	cfg := configuration.Config{
		Network: &netCfg,
	}

	stCfg := configuration.Storage{
		FileStoragePath: "/tmp/short-url-db.json",
	}

	log, _ := zap.NewDevelopment()
	zlog := logger.Logger{
		Log: log,
	}
	storageURL, err := inmemory.NewInMemShortenStorage(&stCfg, &zlog)
	require.NoError(t, err)
	usrStorage, err := inmemory.NewInMemUserStorage(&zlog)
	require.NoError(t, err)
	urlService := shorten.NewURLService(storageURL, &zlog, &cfg)
	usrService := user.NewUserService(usrStorage, &zlog)
	urlHandler := NewURLHandler(&urlService, &usrService)
	requestBody := "https://practicum.yandex.ru/"
	req := models.Request{
		URL: requestBody,
	}
	body, err := json.Marshal(req)
	if err != nil {
		return
	}

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
		req = httptest.NewRequest(http.MethodGet, "/{id}", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", key)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		w = httptest.NewRecorder()
		urlHandler.HandleGet(w, req)
		res = w.Result()
		assert.Equal(t, http.StatusTemporaryRedirect, res.StatusCode)
		assert.Equal(t, requestBody, res.Header.Get("Location"))
		err = res.Body.Close()
		require.NoError(t, err)
	})

	t.Run("handler test #2", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewReader(body))
		w := httptest.NewRecorder()
		urlHandler.HandleShorten(w, req)
		res := w.Result()
		assert.Equal(t, http.StatusCreated, res.StatusCode)
		assert.Equal(t, "application/json", res.Header.Get("Content-Type"))
		responseBody, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		response := models.Response{}
		err = json.Unmarshal(responseBody, &response)
		require.NoError(t, err)
		url := strings.Split(response.Result, "/")
		key := url[len(url)-1]
		err = res.Body.Close()
		require.NoError(t, err)
		req = httptest.NewRequest(http.MethodGet, "/{id}", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", key)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		w = httptest.NewRecorder()
		urlHandler.HandleGet(w, req)
		res = w.Result()
		assert.Equal(t, http.StatusTemporaryRedirect, res.StatusCode)
		assert.Equal(t, requestBody, res.Header.Get("Location"))
		err = res.Body.Close()
		require.NoError(t, err)
	})
}
