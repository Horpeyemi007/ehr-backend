package model

import (
	"backend/ehr/internal/utils"
	"context"
	"database/sql"
	"errors"
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
	UserType   string
	IdType     string
	IdNumber   string
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
	query := `INSERT INTO users (fullname, email, phone, dob, gender, address, occupation, slug, user_type, id_type, id_number) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING slug`

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
		user.UserType,
		user.IdType,
		user.IdNumber,
	).Scan(&user.Slug.Value)

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

func (a *userStore) Find(ctx context.Context, email string) (*User, error) {
	query := `SELECT id, email FROM users
	WHERE email = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	user := &User{}
	err := a.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
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

func (s *userStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `SELECT	id, fullname, email, slug, user_type FROM users
	WHERE email = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	user := &User{}

	err := s.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.Slug.Value,
		&user.UserType,
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
