package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/voutasaurus/cbtest/database"
)

type Handler struct {
	mux *http.ServeMux
	db  *database.DB
	log *log.Logger
}

type Config struct {
	Log *log.Logger
	DB  *database.Config
}

func NewHandler(c *Config) (*Handler, error) {
	db, err := database.NewDB(c.DB)
	if err != nil {
		return nil, err
	}
	logger := c.Log
	if logger == nil {
		logger = log.New(ioutil.Discard, "", 0)
	}
	h := Handler{
		mux: http.NewServeMux(),
		db:  db,
		log: logger,
	}
	h.mux.HandleFunc("/new", h.newHandler)
	h.mux.HandleFunc("/get", h.getHandler)
	h.mux.HandleFunc("/search", h.searchHandler)
	h.mux.HandleFunc("/update", h.updateHandler)
	h.mux.HandleFunc("/delete", h.deleteHandler)
	h.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		h.log.Println("hit")
		fmt.Fprintln(w, "hello")
	})
	return &h, nil
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func (h *Handler) newHandler(w http.ResponseWriter, r *http.Request) {
	var in database.Record
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		jsonError(w, r, fmt.Sprintf("bad json payload: %v", err), 400)
		return
	}
	out, err := h.db.New(r.Context(), &in)
	if err != nil {
		jsonError(w, r, fmt.Sprintf("bad database: %v", err), 500)
		return
	}
	json.NewEncoder(w).Encode(out)
}

func (h *Handler) getHandler(w http.ResponseWriter, r *http.Request) {
	var in database.Record
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		jsonError(w, r, fmt.Sprintf("bad json payload: %v", err), 400)
		return
	}
	out, err := h.db.Get(r.Context(), &in)
	if err != nil {
		jsonError(w, r, fmt.Sprintf("bad database: %v", err), 500)
		return
	}
	json.NewEncoder(w).Encode(out)
}

func (h *Handler) searchHandler(w http.ResponseWriter, r *http.Request) {
	var in database.Record
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		jsonError(w, r, fmt.Sprintf("bad json payload: %v", err), 400)
		return
	}
	out, err := h.db.Search(r.Context(), &in)
	if err != nil {
		jsonError(w, r, fmt.Sprintf("bad database: %v", err), 500)
		return
	}
	json.NewEncoder(w).Encode(out)
}

func (h *Handler) updateHandler(w http.ResponseWriter, r *http.Request) {
	var in database.Record
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		jsonError(w, r, fmt.Sprintf("bad json payload: %v", err), 400)
		return
	}
	out, err := h.db.Update(r.Context(), &in)
	if err != nil {
		jsonError(w, r, fmt.Sprintf("bad database: %v", err), 500)
		return
	}
	json.NewEncoder(w).Encode(out)
}

func (h *Handler) deleteHandler(w http.ResponseWriter, r *http.Request) {
	var in database.Record
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		jsonError(w, r, fmt.Sprintf("bad json payload: %v", err), 400)
		return
	}
	out, err := h.db.Delete(r.Context(), &in)
	if err != nil {
		jsonError(w, r, err.Error(), 500)
		return
	}
	json.NewEncoder(w).Encode(out)
}

func jsonError(w http.ResponseWriter, r *http.Request, msg string, code int) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(struct {
		Err string `json:"err"`
	}{
		Err: msg,
	})
}
