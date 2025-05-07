package main

import (
	"backend/ehr/internal/dtos"
	"backend/ehr/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *application) createPatientHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload dtos.CreatePatientRequest
		// read the payload
		if err := c.ShouldBindJSON(&payload); err != nil {
			badRequestError(c, err)
			return
		}
		// Validate the payload
		if err := Validate.Struct(payload); err != nil {
			badRequestError(c, err)
			return
		}

		// Check if the email exist on the db
		isFound, _ := app.store.Patients.FindByEmail(c, payload.Email)
		if isFound != nil {
			badRequestError(c, model.ErrPatientDuplicateEmail)
			return
		}

		patient := &model.Patient{
			FullName:           payload.FullName,
			DOB:                payload.DOB,
			Gender:             payload.Gender,
			Phone:              payload.Phone,
			Email:              payload.Email,
			Address:            payload.Address,
			Occupation:         payload.Occupation,
			EmergencyName:      payload.EmergencyName,
			EmergencyTelephone: payload.EmergencyTelephone,
		}
		// generate random string (slug)
		if err := patient.Slug.GenerateSlug(3, true, false); err != nil {
			internalServerError(c, err)
		}

		ctx := c.Request.Context()
		// create the patient
		if err := app.store.Patients.Create(ctx, patient); err != nil {
			switch err {
			case model.ErrPatientDuplicateEmail:
				conflictResponse(c, err)
				return
			default:
				internalServerError(c, err)
			}
			return
		}

		response := &dtos.PatientDTO{Slug: patient.Slug.Value}
		jsonResponse(c, http.StatusCreated, response)
	}
}

func (app *application) getAllPatient() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx := c.Request.Context()

		patient, err := app.store.Patients.GetAll(ctx)

		if err != nil {
			internalServerError(c, err)
			return
		}

		var response []dtos.GetAllPatientResponse
		for _, p := range patient {
			dto := dtos.GetAllPatientResponse{
				Slug:               p.Slug.Value,
				FullName:           p.FullName,
				Email:              p.Email,
				Gender:             p.Gender,
				Phone:              p.Phone,
				DOB:                p.DOB,
				Address:            p.Address,
				Occupation:         p.Occupation,
				EmergencyName:      p.EmergencyName,
				EmergencyTelephone: p.EmergencyTelephone,
			}
			response = append(response, dto)
		}

		jsonResponse(c, http.StatusOK, response)
	}
}
