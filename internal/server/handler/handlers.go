package handler

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/lookeme/short-url/internal/app/domain/shorten"
	"github.com/lookeme/short-url/internal/configuration"
	"github.com/lookeme/short-url/internal/models"
	"io"
	"net/http"
	"net/url"
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
	var request models.Request
	body, _ := io.ReadAll(req.Body)
	if err := json.Unmarshal(body, &request); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
	}
	_, err := url.Parse(request.Url)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
	}
	val, err := h.urlService.CreateAndSave(request.Url)
	if err != nil {
		fmt.Println("error during creating hash ", err.Error())
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	result := fmt.Sprintf("%s/%s", h.cfg.BaseURL, val)
	response := models.Response{
		Result: result,
	}

	b, err := json.Marshal(response)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
	}
	_, err = res.Write(b)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
	}
}
func (h *URLHandler) HandleGet(res http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
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
