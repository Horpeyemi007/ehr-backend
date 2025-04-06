package model

import (
	"context"
	"database/sql"
)

// define the patient struct model
type Patient struct {
	ID        int64  `json:"id"`
	FullName  string `json:"fullname"`
	DOB       string `json:"dob"`
	Sex       string `json:"sex"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Street    string `json:"street"`
	City      string `json:"city"`
	State     string `json:"state"`
	Zipcode   string `json:"zipcode"`
	CreatedAt string `json:"created_at"`
}

// define the patient repository
type PatientRepository interface {
	Create(context.Context, *Patient) error
}

// define a db instance struct
type PatientStore struct {
	db *sql.DB
}

func NewPatientRepository(db *sql.DB) PatientRepository {
	return &PatientStore{db: db}
}

func (s *PatientStore) Create(ctx context.Context, patient *Patient) error {
	query := `INSERT INTO patient (fullname, dob, sex, phone, email, street, city, state, zipcode)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id, fullname, email, created_at`

	// create a timeout context on the query
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		patient.FullName,
		patient.DOB,
		patient.Sex,
		patient.Phone,
		patient.Email,
		patient.Street,
		patient.City,
		patient.State,
		patient.Zipcode,
	).Scan(&patient.ID, &patient.FullName, &patient.Email, &patient.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}
