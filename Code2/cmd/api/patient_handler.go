package main

import (
	"backend/ehr/internal/dtos"
	"backend/ehr/internal/model"
	"net/http"
)

func (app *application) createPatientHandler(w http.ResponseWriter, r *http.Request) {

	var payload dtos.CreatePatientRequest
	// read the payload
	if err := readJSON(w, r, &payload); err != nil {
		badRequestError(w, r, err)
		return
	}

	// Validate the payload
	if err := Validate.Struct(payload); err != nil {
		badRequestError(w, r, err)
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
	if err := app.store.Patients.Create(ctx, patient); err != nil {
		switch err {
		case model.ErrDuplicateEmail:
			badRequestError(w, r, err)
			return
		default:
			internalServerError(w, r, err)
		}
		return
	}
	// construct the response data
	response := dtos.PatientDTO{ID: patient.ID, FullName: patient.FullName, Email: patient.Email, CreatedAt: patient.CreatedAt}
	if err := jsonResponse(w, http.StatusCreated, response); err != nil {
		internalServerError(w, r, err)
		return
	}
}
