package main

import (
	"github.com/gin-gonic/gin"
)

func setupRouter(app *application) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	v1 := r.Group("/api")
	{
		// -> (1) admins route
		admin := v1.Group("/admin")
		{
			admin.POST("/register", app.createAdminHandler())
			admin.POST("/login", app.adminLoginAuth())
		}

		// -> (2) users route
		user := v1.Group("/user")
		{
			user.POST("/register", app.createUserHandler())
		}

		// -> (3) patients route
		patient := v1.Group("/patient")
		{
			patient.POST("/register", app.createPatientHandler())
			patient.GET("/all", app.getAllPatient())
		}
	}
	return r
}
