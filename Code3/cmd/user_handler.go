package main

import (
	"backend/ehr/internal/dtos"
	"backend/ehr/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *application) createUserHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload dtos.CreateUserRequest
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
		isFound, _ := app.store.Users.Find(c, payload.Email)
		if isFound != nil {
			badRequestError(c, model.ErrUserDuplicateEmail)
			return
		}

		user := &model.User{
			FullName:   payload.FullName,
			Email:      payload.Email,
			Phone:      payload.Phone,
			DOB:        payload.DOB,
			Gender:     payload.Gender,
			Address:    payload.Address,
			Occupation: payload.Occupation,
			UserType:   payload.UserType,
			IdType:     payload.IdType,
			IdNumber:   payload.IdNumber,
		}

		// generate random string (slug)
		if err := user.Slug.GenerateRandomString(3, true, false); err != nil {
			internalServerError(c, err)
		}

		ctx := c.Request.Context()
		// create the user
		if err := app.store.Users.Create(ctx, user); err != nil {
			switch err {
			case model.ErrUserDuplicateEmail:
				conflictResponse(c, err)
				return
			default:
				internalServerError(c, err)
			}
			return
		}

		jsonResponse(c, http.StatusCreated, "User Created")
	}
}
