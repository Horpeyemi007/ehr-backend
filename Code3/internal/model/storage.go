package model

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrNotFound          = errors.New("resource not found")
	QueryTimeOutDuration = time.Second * 5
	ErrConflict          = errors.New("resource already exists")
	ErrPasswordMismatch  = errors.New("password does not match")
)

type password struct {
	text string
	hash []byte
}

type Storage struct {
	Patients interface {
		Create(context.Context, *Patient) error
		Find(context.Context, string) (*Patient, error)
		GetAll(ctx context.Context) ([]Patient, error)
	}
	Users interface {
		Create(context.Context, *User) error
		Find(context.Context, string) (*User, error)
	}
	Admins interface {
		Create(context.Context, *Admin) error
		GetByEmail(context.Context, string) (*Admin, error)
		Find(context.Context, string) (*Admin, error)
	}
}

// InitializeStore initializes the storage layer with the given database connection
// and returns a Storage instance.
func InitializeStore(db *sql.DB) Storage {
	return Storage{
		Patients: &patientStore{db},
		Users:    &userStore{db},
		Admins:   &adminStore{db},
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
