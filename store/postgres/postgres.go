package postgres

import (
	"context"
	"database/sql"

	"github.com/caioandre182/api-users/domain"
	"github.com/caioandre182/api-users/store"
)

type Store struct {
	db *sql.DB
}

func New(db *sql.DB) *Store { return &Store{db: db} }

func (s *Store) Create(ctx context.Context, u domain.User) (domain.User, error) {
	const q = `
		INSERT INTO users (id, first_name, last_name, biography)
		VALUES ($1, $2, $3, $4)
	`
	_, err := s.db.ExecContext(ctx, q, u.ID, u.FirstName, u.LastName, u.Biography)
	if err != nil {
		return domain.User{}, err
	}
	return u, nil
}

func (s *Store) FindByID(ctx context.Context, id string) (domain.User, error) {
	const q = `
		SELECT id, first_name, last_name, biography
		FROM users
		WHERE id = $1
	`

	var u domain.User
	err := s.db.QueryRowContext(ctx, q, id).
		Scan(&u.ID, &u.FirstName, &u.LastName, &u.Biography)

	if err == sql.ErrNoRows {
		return domain.User{}, store.ErrNotFound
	}
	if err != nil {
		return domain.User{}, err
	}

	return u, nil
}

func (s *Store) FindAll(ctx context.Context) ([]domain.User, error) {
	const q = `
		SELECT id, first_name, last_name, biography
		FROM users
		ORDER BY first_name
	`

	rows, err := s.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]domain.User, 0)
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Biography); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (s *Store) Update(ctx context.Context, u domain.User) error {
	const q = `
		UPDATE users
		SET first_name = $2,
		    last_name = $3,
		    biography = $4
		WHERE id = $1
	`

	res, err := s.db.ExecContext(ctx, q, u.ID, u.FirstName, u.LastName, u.Biography)
	if err != nil {
		return err
	}

	aff, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if aff == 0 {
		return store.ErrNotFound
	}

	return nil
}

func (s *Store) Delete(ctx context.Context, id string) error {
	const q = `
		DELETE FROM users
		WHERE id = $1
	`

	res, err := s.db.ExecContext(ctx, q, id)
	if err != nil {
		return err
	}

	aff, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if aff == 0 {
		return store.ErrNotFound
	}

	return nil
}
