package model

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrNotFound          = errors.New("no record found")
	QueryTimeOutDuration = time.Second * 5
	ErrConflict          = errors.New("resource already exists")
	ErrRoleNameConflict  = errors.New("role name already exists")
	ErrPasswordMismatch  = errors.New("password does not match")
	ErrNoChangesMade     = errors.New("no changes made")
)

type Storage struct {
	Patients interface {
		Create(context.Context, *Patient) error
		FindByEmail(context.Context, string) (*Patient, error)
		GetAll(context.Context) ([]Patient, error)
	}
	Users interface {
		CreateAdmin(context.Context, *User) error
		Create(context.Context, *User) error
		FindById(context.Context, string) (*User, error)
		VerifyPermission(context.Context, int64, []string) (bool, error)
	}
	Entity interface {
		Create(context.Context, *Entity) error
		FindByEmail(context.Context, string) (*Entity, error)
		FindById(context.Context, string) (*Entity, error)
		CreateRole(context.Context, *Role) error
		AssignPermissionToRole(context.Context, string, []string) error
		AssignRoleToUser(context.Context, string, string) error
		AssignPermissionToUser(context.Context, string, []string) error
		RemovePermissionFromUser(context.Context, string, string) error
		RemovePermissionFromRole(context.Context, string, string) error
		GetAllRoles(context.Context) ([]Role, error)
		GetAllPermissions(context.Context) ([]Permission, error)
	}
}

// InitializeStore initializes the storage layer with the given database connection
// and returns a Storage instance.
func InitializeStore(db *sql.DB) Storage {
	return Storage{
		Patients: &patientStore{db},
		Users:    &userStore{db},
		Entity:   &entityStore{db},
	}
}

// hashing the password
func (p *password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	p.text = text
	p.hash = hash

	return nil
}

func (p *password) Compare(text string) error {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(text))
	if err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return ErrPasswordMismatch
		default:
			return err
		}
	}
	return nil
}

// wrapper function using sql transactions
func withTransaction(ctx context.Context, db *sql.DB, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		// rollback on an error
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
