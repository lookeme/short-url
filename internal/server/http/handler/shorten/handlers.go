package shorten

import (
	"github.com/lookeme/short-url/internal/app/domain/shorten"
	"github.com/lookeme/short-url/internal/configuration"
	"io"
	"net/http"
)

type UrlHandler struct {
	urlService *shorten.UrlService
	cfg        *configuration.NetworkCfg
}

func NewUrlHandler(urlService *shorten.UrlService, cfg *configuration.NetworkCfg) *UrlHandler {
	return &UrlHandler{
		urlService: urlService,
		cfg:        cfg,
	}
}

func (h *UrlHandler) Index(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleGet(w, r)
	case http.MethodPost:
		h.handlePOST(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *UrlHandler) handlePOST(res http.ResponseWriter, req *http.Request) {
	b, err := io.ReadAll(req.Body)
	val, err := h.urlService.CreateAndSave(string(b))
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
	}
	res.Header().Set("content-type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	var host string
	if h.cfg == nil {
		host = "localhost:8080"
	}
	url := "http://" + host + "/" + val
	res.Write([]byte(url))
}
func (h *UrlHandler) handleGet(res http.ResponseWriter, req *http.Request) {
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
