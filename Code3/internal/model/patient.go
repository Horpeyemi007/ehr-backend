package model

import (
	"backend/ehr/internal/utils"
	"context"
	"database/sql"
	"errors"
)

// define the patient struct model
type Patient struct {
	ID                int64
	FullName          string
	DOB               string
	Gender            string
	Phone             string
	Email             string
	Slug              utils.Slug
	Address           string
	Occupation        string
	GuardianName      string
	GuardianTelephone string
	CreatedAt         string
}

// define a db instance patient struct
type patientStore struct {
	db *sql.DB
}

var (
	ErrPatientDuplicateEmail = errors.New("a patient with the given email already exists")
)

func (p *patientStore) Create(ctx context.Context, patient *Patient) error {
	query := `INSERT INTO patient (fullname, dob, gender, phone, email, slug, address, occupation, guardian_name, guardian_telephone)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id, fullname, email, slug, created_at`

	// create a timeout context on the query
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	err := p.db.QueryRowContext(
		ctx,
		query,
		patient.FullName,
		patient.DOB,
		patient.Gender,
		patient.Phone,
		patient.Email,
		patient.Slug.Value,
		patient.Address,
		patient.Occupation,
		patient.GuardianName,
		patient.GuardianTelephone,
	).Scan(&patient.ID, &patient.FullName, &patient.Email, &patient.Slug.Value, &patient.CreatedAt)

	if err != nil {
		switch {
		case err.Error() == "pq: duplicate key value violates unique constraint \"patient_email_key\"":
			return ErrPatientDuplicateEmail
		default:
			return err
		}
	}
	return nil
}
