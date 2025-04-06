package model

import (
	"errors"
	"time"
)

var (
	ErrNotFound          = errors.New("resource not found")
	QueryTimeOutDuration = time.Second * 5
	ErrConflict          = errors.New("resource already exists")
)

// type Storage struct {
// 	Patients interface {
// 		Create(context.Context, *Patient) error
// 	}
// }

// // Initialize the store
// func InitializeStore(db *sql.DB) Storage {
// 	return Storage{
// 		Patients: &PatientStore{db},
// 	}
// }
