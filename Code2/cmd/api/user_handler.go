package main

import (
	"backend/ehr/internal/dtos"
	"backend/ehr/internal/model"
	"net/http"
)

func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload dtos.CreateUserRequest
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

	user := &model.User{
		FullName:   payload.FullName,
		Email:      payload.Email,
		Phone:      payload.Phone,
		DOB:        payload.DOB,
		Gender:     payload.Gender,
		Address:    payload.Address,
		Occupation: payload.Occupation,
		Role:       payload.Role,
	}

	// hash the user password
	if err := user.Password.Set(payload.Password); err != nil {
		internalServerError(w, r, err)
		return
	}
	// generate random string (slug)
	if err := user.Slug.GenerateRandomString(5, true, false); err != nil {
		internalServerError(w, r, err)
	}

	ctx := r.Context()
	// create the user
	if err := app.store.Users.Create(ctx, user); err != nil {
		switch err {
		case model.ErrUserDuplicateEmail:
			badRequestError(w, r, err)
			return
		default:
			internalServerError(w, r, err)
		}
		return
	}

	if err := jsonResponse(w, http.StatusCreated, "User created"); err != nil {
		internalServerError(w, r, err)
		return
	}
}
