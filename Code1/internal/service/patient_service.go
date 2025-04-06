package service

import (
	"backend/ehr/internal/model"
	"context"
)

type PatientService interface {
	CreatePatient(context.Context, *model.Patient) error
}

type PatientServiceImp struct {
	model model.PatientRepository
}

func NewPatientService(m model.PatientRepository) PatientService {
	return &PatientServiceImp{model: m}
}

func (s *PatientServiceImp) CreatePatient(ctx context.Context, patient *model.Patient) error {
	err := s.model.Create(ctx, patient)
	if err != nil {
		return err
	}
	return nil
}
