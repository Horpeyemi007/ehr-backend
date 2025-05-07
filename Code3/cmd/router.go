package main

import (
	"backend/ehr/internal/utils"

	"github.com/gin-gonic/gin"
)

func setupRouter(app *application) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	v1 := r.Group("/api")
	{
		// -> (1) entity route
		entity := v1.Group("/entity")
		{
			entity.POST("/register", app.createEntityHandler())
			entity.POST("/login", app.entityLoginAuth())
			entity.GET("/:id", app.getEntityByIdHandler())
			entity.POST("/role/create", app.createRoleHandler())
			entity.POST("/permission/role", app.assignPermissionToRoleHandler())
			entity.POST("/role/user", app.assignRoleToUserHandler())
			entity.POST("/permission/user", app.assignPermissionToUserHandler())
			entity.DELETE("/permission/user/delete", app.removePermissionFromUser())
			entity.DELETE("/permission/role/delete", app.removePermissionFromRole())
			entity.GET("/role/all", app.getAllRoles())
			entity.GET("/permission/all", app.getAllPermissions())
		}

		// -> (2) users route
		user := v1.Group("/user")
		{
			user.POST("/admin/register", app.createAdminUserHandler())
			user.POST("/register", app.createUserHandler())
		}

		// -> (3) patients route
		patient := v1.Group("/patient")
		{
			patient.POST("/register", app.createPatientHandler())
			patient.GET("/all", app.permit(utils.TestAccess, app.getAllPatient()))
		}
	}
	return r
}
