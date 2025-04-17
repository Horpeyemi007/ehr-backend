package main

import (
	"backend/ehr/internal/dtos"
	"backend/ehr/internal/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func (app *application) adminLoginAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload dtos.AdminLoginRequest
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

		ctx := c.Request.Context()

		admin, err := app.store.Admins.GetByEmail(ctx, payload.Email)

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
		if err := admin.Password.Compare(payload.Password); err != nil {
			unauthorizedErrorResponse(c, err)
			return
		}

		// generate the token -< add claims
		claims := jwt.MapClaims{
			"sub": admin.ID,
			"exp": time.Now().Add(app.config.Auth.Token.Exp).Unix(),
			"iat": time.Now().Unix(),
			"nbf": time.Now().Unix(),
			"iss": app.config.Auth.Token.Iss,
			"aud": app.config.Auth.Token.Iss,
		}
		token, err := app.authenticator.GenerateToken(claims)
		if err != nil {
			internalServerError(c, err)
			return
		}

		response := &dtos.AdminDTO{
			ID:       admin.ID,
			FullName: admin.FullName,
			Email:    admin.Email,
			Slug:     admin.Slug.Value,
			Token:    token,
		}

		jsonResponse(c, http.StatusOK, response)
	}
}
