package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/faizainur/hands-on-golang/rest-server/models"
	"github.com/faizainur/hands-on-golang/rest-server/services"
	"github.com/gorilla/mux"
)

type Password struct {
	Username_email string `json:"username_email,omitempty"`
	Password       string `json:"password,omitempty"`
}

type PasswordRoutes struct {
	service *services.PasswordService
}

func NewPasswordRoutes(s *services.PasswordService) *PasswordRoutes {
	return &PasswordRoutes{service: s}
}

func (p *PasswordRoutes) AddPassword(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	dbResponse, err := p.service.AddPassword(vars["name"], vars["username"], vars["email"], vars["password"])
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "%s", "Bad Request")
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(dbResponse)
}

func (p *PasswordRoutes) GetPassword(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(400)
		return
	}

	if id < 1 {
		w.WriteHeader(405)
		return
	}

	response, err := p.service.GetPassword(uint(id))
	if err != nil {
		// log.Fatal(err)
		fmt.Fprintf(w, "%s", err.Error())
		w.WriteHeader(400)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (p *PasswordRoutes) ListPasswords(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "%s", "Hell")
	limit, err := strconv.Atoi(mux.Vars(r)["limit"])
	if err != nil {
		fmt.Fprintf(w, "%s", err.Error())
		w.WriteHeader(400)
		return
	}
	offset, err := strconv.Atoi(mux.Vars(r)["offset"])
	if err != nil {
		fmt.Fprintf(w, "%s", err.Error())
		w.WriteHeader(400)
		return
	}
	// fmt.Fprintf(w, "%d %d", limit, offset)
	response, count, err := p.service.ListPasswords(uint32(limit), uint32(offset))
	if err != nil {
		fmt.Fprintf(w, "%s", err.Error())
		w.WriteHeader(400)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Count  int64             `json:"count"`
		Limit  uint              `json:"limit"`
		Offset uint              `json:"offset"`
		Data   []models.Password `json:"data"`
	}{
		Count:  *count,
		Limit:  uint(limit),
		Offset: uint(offset),
		Data:   response,
	})
}
