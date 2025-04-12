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
		FullName:          payload.FullName,
		DOB:               payload.DOB,
		Gender:            payload.Gender,
		Phone:             payload.Phone,
		Email:             payload.Email,
		Address:           payload.Address,
		Occupation:        payload.Occupation,
		GuardianName:      payload.GuardianName,
		GuardianTelephone: payload.GuardianTelephone,
	}

	// generate random string (slug)
	if err := patient.Slug.GenerateRandomString(5, true, false); err != nil {
		internalServerError(w, r, err)
	}

	ctx := r.Context()
	// create the patient
	if err := app.store.Patients.Create(ctx, patient); err != nil {
		switch err {
		case model.ErrPatientDuplicateEmail:
			badRequestError(w, r, err)
			return
		default:
			internalServerError(w, r, err)
		}
		return
	}

	if err := jsonResponse(w, http.StatusCreated, "Patient created"); err != nil {
		internalServerError(w, r, err)
		return
	}
}
