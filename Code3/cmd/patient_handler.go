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
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Validate the payload
		if err := Validate.Struct(payload); err != nil {
			badRequestError(c, err)
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
		jsonResponse(c, http.StatusCreated, "Patient Created")
	}
}
