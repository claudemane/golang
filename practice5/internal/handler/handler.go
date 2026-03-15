package handler

import (
	"encoding/json"
	"net/http"
	"practice5/internal/models"
	"practice5/internal/repository"
	"strconv"
)

type Handler struct {
	repo *repository.Repository
}

func New(repo *repository.Repository) *Handler {
	return &Handler{repo: repo}
}

// GET /users?page=1&page_size=5&order_by=name&name=john&gender=male&id=3&email=...&birth_date=1990-01-01
func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	q := r.URL.Query()

	page, _ := strconv.Atoi(q.Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(q.Get("page_size"))
	if pageSize < 1 {
		pageSize = 10
	}

	filter := models.UserFilter{
		Page:     page,
		PageSize: pageSize,
		OrderBy:  q.Get("order_by"),
	}

	if v := q.Get("id"); v != "" {
		id, err := strconv.Atoi(v)
		if err == nil {
			filter.ID = &id
		}
	}
	if v := q.Get("name"); v != "" {
		filter.Name = &v
	}
	if v := q.Get("email"); v != "" {
		filter.Email = &v
	}
	if v := q.Get("gender"); v != "" {
		filter.Gender = &v
	}
	if v := q.Get("birth_date"); v != "" {
		filter.BirthDate = &v
	}

	result, err := h.repo.GetPaginatedUsers(filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	json.NewEncoder(w).Encode(result)
}

// GET /users/common-friends?user1=1&user2=2
func (h *Handler) GetCommonFriends(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	q := r.URL.Query()

	user1, err1 := strconv.Atoi(q.Get("user1"))
	user2, err2 := strconv.Atoi(q.Get("user2"))

	if err1 != nil || err2 != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "user1 and user2 query params are required"})
		return
	}

	friends, err := h.repo.GetCommonFriends(user1, user2)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	json.NewEncoder(w).Encode(friends)
}
