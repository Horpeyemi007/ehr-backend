package model

import (
	"backend/ehr/internal/utils"
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

// define the user struct model
type User struct {
	ID         int64
	FullName   string
	Email      string
	Phone      string
	DOB        string
	Gender     string
	Address    string
	Occupation string
	Role       []string
	Slug       utils.Slug
	CreatedAt  string
	IsActive   bool
}

// define the db instance user struct
type userStore struct {
	db *sql.DB
}

var (
	ErrUserDuplicateEmail = errors.New("a user with the given email already exists")
)

func (u *userStore) Create(ctx context.Context, user *User) error {
	query := `INSERT INTO users (fullname, email, phone, dob, gender, address, occupation, slug, role) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id, fullname, email, slug, role, created_at`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	err := u.db.QueryRowContext(
		ctx,
		query,
		user.FullName,
		user.Email,
		user.Phone,
		user.DOB,
		user.Gender,
		user.Address,
		user.Occupation,
		user.Slug.Value,
		pq.Array(user.Role),
	).Scan(&user.ID, &user.FullName, &user.Email, &user.Slug.Value, pq.Array(&user.Role), &user.CreatedAt)

	if err != nil {
		switch {
		case err.Error() == "pq: duplicate key value violates unique constraint \"users_email_key\"":
			return ErrUserDuplicateEmail
		default:
			return err
		}
	}
	return nil
}
func (s *userStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `SELECT	id, fullname, email, slug, role FROM users
	WHERE email = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	user := &User{}

	err := s.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.Slug.Value,
		pq.Array(&user.Role),
	)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return user, nil
}
