package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New(validator.WithPrivateFieldValidation())
}

func (app *application) setCSRFToken(c *gin.Context, csrfValue string) {
	c.Header("X-CSRF-Token", csrfValue)
}

func jsonResponse(c *gin.Context, status int, data any) {
	c.JSON(status, gin.H{
		"success": true,
		"data":    data,
	})
}
