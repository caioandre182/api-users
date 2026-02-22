package api

import (
	"encoding/json"
	"net/http"

	"github.com/caioandre182/api-users/domain"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

func JsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func NewRouter(db map[string]domain.User) http.Handler {
	r := chi.NewRouter()

	r.Use(JsonMiddleware)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	r.Get("/user/{id}", getUserByID(db))
	r.Get("/users", getAllUsers(db))
	r.Post("/user", postUser(db))
	r.Put("/user/{id}", putUser(db))
	r.Delete("/user/{id}", deleteUser(db))

	return r
}

func getUserByID(db map[string]domain.User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			http.Error(w, "missing id", http.StatusBadRequest)
			return
		}

		user, ok := db[id]
		if !ok {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}

		if err := json.NewEncoder(w).Encode(user); err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	}
}

func getAllUsers(db map[string]domain.User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users := make([]domain.User, 0, len(db))

		for _, user := range db {
			users = (append(users, user))
		}

		if err := json.NewEncoder(w).Encode(users); err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	}
}

func postUser(db map[string]domain.User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user domain.User

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}

		user.ID = uuid.NewString()

		db[user.ID] = user
		w.WriteHeader(http.StatusCreated)

		if err := json.NewEncoder(w).Encode(user); err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	}
}

func putUser(db map[string]domain.User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		var newUser domain.User

		if id == "" {
			http.Error(w, "missing id", http.StatusBadRequest)
			return
		}

		oldUser, ok := db[id]
		if !ok {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}

		oldUser.FirstName = newUser.FirstName
		oldUser.LastName = newUser.LastName
		oldUser.Biography = newUser.Biography

		db[id] = oldUser

		if err := json.NewEncoder(w).Encode(oldUser); err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	}
}

func deleteUser(db map[string]domain.User) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		if id == "" {
			http.Error(w, "missing id", http.StatusBadRequest)
			return
		}

		if _, ok := db[id]; !ok {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}

		delete(db, id)

		w.WriteHeader(http.StatusNoContent)
	}
}
