package handler

import (
	"fmt"
	"github.com/lookeme/short-url/internal/app/domain/shorten"
	"github.com/lookeme/short-url/internal/configuration"
	"io"
	"net/http"
)

type URLHandler struct {
	urlService *shorten.URLService
	cfg        *configuration.NetworkCfg
}

func NewURLHandler(urlService *shorten.URLService, cfg *configuration.Config) *URLHandler {
	return &URLHandler{
		urlService: urlService,
		cfg:        cfg.Network,
	}
}

func (h *URLHandler) HandlePOST(res http.ResponseWriter, req *http.Request) {
	b, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
	}
	val, err := h.urlService.CreateAndSave(string(b))
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
	}
	res.Header().Set("content-type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	url := fmt.Sprintf("%s/%s", h.cfg.BaseAddress, val)
	_, err = res.Write([]byte(url))
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
	}
}
func (h *URLHandler) HandleGet(res http.ResponseWriter, req *http.Request) {
	id := req.RequestURI[1:]
	if id == "" {
		http.Error(res, "ID is not provided in path", http.StatusBadRequest)
		return
	}
	val, ok := h.urlService.FindByKey(id)
	if !ok {
		http.Error(res, "Value is not found", http.StatusBadRequest)
		return
	}
	res.Header().Set("Location", val)
	res.WriteHeader(http.StatusTemporaryRedirect)
}
