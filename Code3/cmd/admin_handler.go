package main

import (
	"backend/ehr/internal/dtos"
	"backend/ehr/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *application) createAdminHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload dtos.CreateAdminRequest
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

		// Check if the email exist on the db
		isFound, _ := app.store.Admins.Find(c, payload.Email)
		if isFound != nil {
			badRequestError(c, model.ErrAdminDuplicateEmail)
			return
		}

		admin := &model.Admin{
			FullName: payload.FullName,
			Email:    payload.Email,
		}

		// hash the admin password
		if err := admin.Password.Set(payload.Password); err != nil {
			internalServerError(c, err)
			return
		}

		// generate random string (slug)
		if err := admin.Slug.GenerateRandomString(3, true, false); err != nil {
			internalServerError(c, err)
		}

		ctx := c.Request.Context()
		// create the user
		if err := app.store.Admins.Create(ctx, admin); err != nil {
			switch err {
			case model.ErrAdminDuplicateEmail:
				conflictResponse(c, err)
				return
			default:
				internalServerError(c, err)
			}
			return
		}

		jsonResponse(c, http.StatusCreated, "Admin Created")
	}
}
