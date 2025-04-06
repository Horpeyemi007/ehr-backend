package main

import (
	"backend/ehr/internal/model"
	"backend/ehr/internal/service"
	"net/http"
)

type patientHandler struct {
	patient service.PatientService
}

func NewPatientHandler(ps service.PatientService) *patientHandler {
	return &patientHandler{patient: ps}
}

// create the patient
func (h *patientHandler) CreatePatientHandler(w http.ResponseWriter, r *http.Request) {
	type createPatientPayload struct {
		FullName string `json:"fullname"`
		DOB      string `json:"dob"`
		Sex      string `json:"sex"`
		Phone    string `json:"phone"`
		Email    string `json:"email"`
		Street   string `json:"street"`
		City     string `json:"city"`
		State    string `json:"state"`
		Zipcode  string `json:"zipcode"`
	}

	var payload createPatientPayload
	if err := readJSON(w, r, &payload); err != nil {
		internalServerError(w, r, err)
		return
	}

	patient := &model.Patient{
		FullName: payload.FullName,
		DOB:      payload.DOB,
		Sex:      payload.Sex,
		Phone:    payload.Phone,
		Email:    payload.Email,
		Street:   payload.Street,
		City:     payload.City,
		State:    payload.State,
		Zipcode:  payload.Zipcode,
	}

	ctx := r.Context()
	err := h.patient.CreatePatient(ctx, patient)
	if err != nil {
		badRequestError(w, r, err)
		return
	}

	if err := jsonResponse(w, http.StatusCreated, patient); err != nil {
		internalServerError(w, r, err)
		return
	}
}
