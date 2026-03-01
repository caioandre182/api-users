package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/caioandre182/api-users/domain"
	"github.com/caioandre182/api-users/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Handler struct {
	store store.UserStore
}

func New(s store.UserStore) *Handler {
	return &Handler{store: s}
}

func JsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) Router() http.Handler {
	r := chi.NewRouter()

	r.Use(JsonMiddleware)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	r.Get("/user/{id}", h.getUserByID())
	r.Get("/users", h.getAllUsers())
	r.Post("/user", h.postUser())
	r.Put("/user/{id}", h.putUser())
	r.Delete("/user/{id}", h.deleteUser())

	return r
}

func (h *Handler) getUserByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			http.Error(w, "missing id", http.StatusBadRequest)
			return
		}

		user, err := h.store.FindByID(r.Context(), id)

		if err != nil {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}

		_ = json.NewEncoder(w).Encode(user)

	}
}

func (h *Handler) getAllUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := h.store.FindAll(r.Context())

		if err != nil {
			http.Error(w, "something wrong", http.StatusInternalServerError)
			return
		}

		_ = json.NewEncoder(w).Encode(users)
	}
}

func (h *Handler) postUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user domain.User

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			slog.Error("create user failed", "error", err)
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}

		userCreated, err := h.store.Create(r.Context(), user)

		if err != nil {
			http.Error(w, "something wrong", http.StatusInternalServerError)
			return
		}

		_ = json.NewEncoder(w).Encode(userCreated)
		w.WriteHeader(http.StatusCreated)

	}
}

func (h *Handler) putUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		if id == "" {
			http.Error(w, "missing id", http.StatusBadRequest)
			return
		}

		var user domain.User

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}

		_, err := h.store.FindByID(r.Context(), id)

		if err != nil {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}

		h.store.Update(r.Context(), user)

		w.WriteHeader(http.StatusNoContent)
	}
}

func (h *Handler) deleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		if id == "" {
			http.Error(w, "missing id", http.StatusBadRequest)
			return
		}

		_, err := h.store.FindByID(r.Context(), id)

		if err != nil {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}

		h.store.Delete(r.Context(), id)

		w.WriteHeader(http.StatusNoContent)
	}
}
