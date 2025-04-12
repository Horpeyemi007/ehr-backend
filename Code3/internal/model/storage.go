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
)

type Password struct {
	text string
	hash []byte
}

type Storage struct {
	Patients interface {
		Create(context.Context, *Patient) error
	}
	Users interface {
		Create(context.Context, *User) error
	}
}

// InitializeStore initializes the storage layer with the given database connection
// and returns a Storage instance.
func InitializeStore(db *sql.DB) Storage {
	return Storage{
		Patients: &patientStore{db},
		Users:    &userStore{db},
	}
}

// hashing the password
func (p *Password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	p.text = text
	p.hash = hash

	return nil
}
