package model

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound          = errors.New("resource not found")
	QueryTimeOutDuration = time.Second * 5
	ErrConflict          = errors.New("resource already exists")
)

type Storage struct {
	Patients interface {
		Create(context.Context, *Patient) error
	}
}

// InitializeStore initializes the storage layer with the given database connection
// and returns a Storage instance.
func InitializeStore(db *sql.DB) Storage {
	return Storage{
		Patients: &PatientStore{db},
	}
}
