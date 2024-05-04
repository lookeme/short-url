package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"

	"github.com/lookeme/short-url/internal/app/domain/shorten"
	"github.com/lookeme/short-url/internal/models"
)

type URLHandler struct {
	urlService *shorten.URLService
}

func NewURLHandler(urlService *shorten.URLService) *URLHandler {
	return &URLHandler{
		urlService: urlService,
	}
}

func (h *URLHandler) HandleShorten(res http.ResponseWriter, req *http.Request) {
	var request models.Request
	body, _ := io.ReadAll(req.Body)
	if err := json.Unmarshal(body, &request); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
	}
	_, err := url.Parse(request.URL)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
	}
	val, err := h.urlService.CreateAndSave(request.URL)
	if err != nil {
		fmt.Println("error during creating hash ", err.Error())
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	b, err := json.Marshal(models.Response{
		Result: val,
	})
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
	}
	_, err = res.Write(b)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
	}
}

func (h *URLHandler) HandlePOST(res http.ResponseWriter, req *http.Request) {
	b, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
	}
	_, err = url.Parse(string(b))
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
	}
	val, err := h.urlService.CreateAndSave(string(b))
	if err != nil {
		fmt.Println("error during creating hashL ", err.Error())
	}
	res.Header().Set("content-type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	_, err = res.Write([]byte(val))
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
	}
}

func (h *URLHandler) HandlePing(res http.ResponseWriter, _ *http.Request) {
	ctx := context.Background()
	err := h.urlService.Ping(ctx)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
	}
	res.WriteHeader(http.StatusOK)
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

func (h *URLHandler) HandleUserURLs(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	urls, err := h.urlService.FindAll()
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("URLS", urls)

	fmt.Println("params", req.URL.String())

	b1, _ := io.ReadAll(req.Body)
	fmt.Println("body", string(b1))
	if urls == nil {
		res.WriteHeader(http.StatusNoContent)
		return
	}
	b, err := json.Marshal(urls)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	res.WriteHeader(http.StatusOK)
	_, err = res.Write(b)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
}
