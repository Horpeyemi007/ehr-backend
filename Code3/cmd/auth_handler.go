package main

import (
	"backend/ehr/internal/dtos"
	"backend/ehr/internal/model"
	"backend/ehr/internal/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func (app *application) entityLoginAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload dtos.EntityLoginRequest
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

		ctx := c.Request.Context()

		entity, err := app.store.Entity.FindByEmail(ctx, payload.Email)

		if err != nil {
			switch err {
			case model.ErrNotFound:
				unauthorizedErrorResponse(c, err)
			default:
				internalServerError(c, err)
			}
			return
		}

		// check for password validity
		if err := entity.Password.Compare(payload.Password); err != nil {
			unauthorizedErrorResponse(c, err)
			return
		}
		// generate the csrf token
		csrfToken := utils.GenerateCSRFToken()
		// generate the token -< add claims
		claims := jwt.MapClaims{
			"sub":      entity.ID,
			"exp":      time.Now().Add(app.config.Auth.Token.Exp).Unix(),
			"iat":      time.Now().Unix(),
			"iss":      app.config.Auth.Token.Iss,
			"aud":      app.config.Auth.Token.Iss,
			"userType": entity.UserType,
			"csrf":     csrfToken,
		}
		token, err := app.authenticator.GenerateToken(claims)
		if err != nil {
			internalServerError(c, err)
			return
		}
		response := &dtos.EntityDTO{
			ID:       entity.ID,
			FullName: entity.FullName,
			Email:    entity.Email,
			Slug:     entity.EntityInformation.Slug.Value,
			Token:    token,
		}
		app.setCSRFToken(c, csrfToken)
		jsonResponse(c, http.StatusOK, response)
	}
}
