package model

import (
	"backend/ehr/internal/utils"
	"context"
	"database/sql"
	"errors"
)

// define the patient struct model
type Patient struct {
	ID                 int64
	FullName           string
	Email              string
	DOB                string
	Gender             string
	Phone              string
	Slug               utils.Slug
	Address            string
	Occupation         string
	EmergencyName      string
	EmergencyTelephone string
	CreatedAt          string
}

// define a db instance patient struct
type patientStore struct {
	db *sql.DB
}

var (
	ErrPatientDuplicateEmail = errors.New("a patient with the given email already exists")
)

func (p *patientStore) Create(ctx context.Context, patient *Patient) error {
	query := `INSERT INTO patients (fullname, dob, gender, phone, email, slug, address, occupation, emergency_name, emergency_phone)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING slug`

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
		patient.EmergencyName,
		patient.EmergencyTelephone,
	).Scan(&patient.Slug.Value)

	if err != nil {
		switch {
		case err.Error() == "pq: duplicate key value violates unique constraint \"patients_email_key\"":
			return ErrPatientDuplicateEmail
		default:
			return err
		}
	}
	return nil
}

func (a *patientStore) Find(ctx context.Context, email string) (*Patient, error) {
	query := `SELECT id, email FROM patients
	WHERE email = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	patient := &Patient{}
	err := a.db.QueryRowContext(ctx, query, email).Scan(
		&patient.ID,
		&patient.Email,
	)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return patient, nil
}

func (p *patientStore) GetAll(ctx context.Context) ([]Patient, error) {
	query := ` 
		SELECT slug, fullname, email, phone, dob, gender, address, occupation, emergency_name, emergency_phone 
		FROM patients 
		ORDER BY id DESC
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := p.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	patients := []Patient{}
	for rows.Next() {
		var p Patient
		err := rows.Scan(
			&p.Slug.Value,
			&p.FullName,
			&p.Email,
			&p.Phone,
			&p.DOB,
			&p.Gender,
			&p.Address,
			&p.Occupation,
			&p.EmergencyName,
			&p.EmergencyTelephone,
		)
		if err != nil {
			return nil, err
		}
		patients = append(patients, p)
	}
	return patients, nil
}
