package model

import (
	"backend/ehr/internal/utils"
	"context"
	"database/sql"
	"errors"
)

// define the patient struct model
type Admin struct {
	ID        int64
	FullName  string
	Email     string
	Password  password
	Slug      utils.Slug
	CreatedAt string
	IsActive  bool
}

// define a db instance patient struct
type adminStore struct {
	db *sql.DB
}

var (
	ErrAdminDuplicateEmail = errors.New("an admin with the given email already exists")
)

func (a *adminStore) Create(ctx context.Context, admin *Admin) error {
	query := `INSERT INTO admins (fullname, email, password, slug) 
	VALUES ($1, $2, $3, $4) RETURNING id, fullname, email, slug`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	err := a.db.QueryRowContext(
		ctx,
		query,
		admin.FullName,
		admin.Email,
		admin.Password.hash,
		admin.Slug.Value,
	).Scan(&admin.ID, &admin.FullName, &admin.Email, &admin.Slug.Value)

	if err != nil {
		switch {
		case err.Error() == "pq: duplicate key value violates unique constraint \"admins_email_key\"":
			return ErrAdminDuplicateEmail
		default:
			return err
		}
	}
	return nil
}

func (a *adminStore) GetByEmail(ctx context.Context, email string) (*Admin, error) {
	query := `SELECT id, fullname, email, password, slug FROM admins
	WHERE email = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	admin := &Admin{}
	err := a.db.QueryRowContext(ctx, query, email).Scan(
		&admin.ID,
		&admin.FullName,
		&admin.Email,
		&admin.Password.hash,
		&admin.Slug.Value,
	)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return admin, nil
}

func (a *adminStore) Find(ctx context.Context, email string) (*Admin, error) {
	query := `SELECT id, email FROM admins
	WHERE email = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	admin := &Admin{}
	err := a.db.QueryRowContext(ctx, query, email).Scan(
		&admin.ID,
		&admin.Email,
	)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return admin, nil
}
