package handler

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgerrcode"
	"github.com/lookeme/short-url/internal/app/domain/user"
	"github.com/lookeme/short-url/internal/security"
	"github.com/lookeme/short-url/internal/utils"
	"io"
	"net/http"
	"net/url"

	"github.com/lookeme/short-url/internal/app/domain/shorten"
	"github.com/lookeme/short-url/internal/models"
)

type URLHandler struct {
	urlService *shorten.URLService
	usrService *user.UsrService
}

func NewURLHandler(urlService *shorten.URLService, usrService *user.UsrService) *URLHandler {
	return &URLHandler{
		urlService: urlService,
		usrService: usrService,
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
	val, err := h.urlService.CreateAndSave(request.URL, 1)
	res.Header().Set("Content-Type", "application/json")
	if err != nil {
		h.urlService.Log.Log.Error(err.Error())
		code := utils.ErrorCode(err)
		if code == pgerrcode.UniqueViolation {
			res.WriteHeader(http.StatusConflict)
			data, ok := h.urlService.FindByURL(request.URL)
			if !ok {
				http.Error(res, err.Error(), http.StatusBadRequest)
			} else {
				val = data.ShortURL
			}
		}
	} else {
		res.WriteHeader(http.StatusCreated)
	}
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
	defer req.Body.Close()
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
	}
	_, err = url.Parse(string(b))
	if err != nil {

		http.Error(res, err.Error(), http.StatusBadRequest)
	}
	urlToSave := string(b)
	token := req.Header.Get("Authorization")
	token, err = utils.GetToken(token)
	if err != nil {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
	userID := security.GetUserID(token)
	if userID == 0 {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
	val, err := h.urlService.CreateAndSave(urlToSave, userID)
	res.Header().Set("content-type", "text/plain")
	if err != nil {
		h.urlService.Log.Log.Error(err.Error())
		code := utils.ErrorCode(err)
		if code == pgerrcode.UniqueViolation {
			res.WriteHeader(http.StatusConflict)
			data, ok := h.urlService.FindByURL(urlToSave)
			if !ok {
				http.Error(res, err.Error(), http.StatusBadRequest)
			} else {
				val = data.ShortURL
			}
		}
	} else {
		res.WriteHeader(http.StatusCreated)
	}
	if err != nil {
		h.urlService.Log.Log.Error(err.Error())
	}
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
	if val.DeletedFlag {
		res.WriteHeader(http.StatusGone)
	} else {
		res.Header().Set("Location", val.OriginalURL)
		res.WriteHeader(http.StatusTemporaryRedirect)
	}
}

func (h *URLHandler) HandleUserURLs(res http.ResponseWriter, r *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	token := r.Header.Get("Authorization")
	token, err := utils.GetToken(token)
	if err != nil {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
	userID := security.GetUserID(token)
	if userID == 0 {
		http.Error(res, "userID is not presented in token", http.StatusUnauthorized)
	}
	urls, err := h.urlService.FindAllByUserID(userID)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
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

func (h *URLHandler) HandleShortenBatch(res http.ResponseWriter, req *http.Request) {
	var request []models.BatchRequest
	body, err := io.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
	}
	if err := json.Unmarshal(body, &request); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
	}
	val, err := h.urlService.CreateAndSaveBatch(request)
	if err != nil {
		h.urlService.Log.Log.Error(err.Error())
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	b, err := json.Marshal(val)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
	}
	_, err = res.Write(b)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
	}
}

func (h *URLHandler) HandleDeleteURLs(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	token := req.Header.Get("Authorization")
	token, err := utils.GetToken(token)
	if err != nil {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
	userID := security.GetUserID(token)
	if userID == 0 {
		http.Error(res, "userID is not presented in token", http.StatusUnauthorized)
	}
	var request []string
	body, _ := io.ReadAll(req.Body)
	if err := json.Unmarshal(body, &request); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
	}
	if len(request) != 0 {
		go h.urlService.DeleteByShortURLAndUserID(request, userID)
	}
	res.WriteHeader(http.StatusAccepted)
}
