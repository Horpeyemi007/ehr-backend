package main

import (
	"backend/ehr/internal/dtos"
	"backend/ehr/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *application) createAdminUserHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload dtos.CreateAdminUserRequest
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
		if err := user.Slug.GenerateSlug(3, true, false); err != nil {
			internalServerError(c, err)
		}

		ctx := c.Request.Context()
		// create the admin user
		if err := app.store.Users.CreateAdmin(ctx, user); err != nil {
			switch err {
			case model.ErrUserDuplicateEmail:
				conflictResponse(c, err)
				return
			default:
				internalServerError(c, err)
			}
			return
		}

		jsonResponse(c, http.StatusCreated, "Admin User Created")
	}
}

func (app *application) createUserHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload dtos.CreateUserRequest
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
		if err := user.Slug.GenerateSlug(3, true, false); err != nil {
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
