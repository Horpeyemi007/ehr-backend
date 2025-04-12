package main

import (
	"github.com/gin-gonic/gin"
)

func setupRouter(app *application) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	v1 := r.Group("/api")
	{
		// users route
		user := v1.Group("/user")
		{
			user.POST("/register", app.createUserHandler())
		}
		// patients route
		patient := v1.Group("/patient")
		{
			patient.POST("/register", app.createPatientHandler())
		}
	}
	return r
}
