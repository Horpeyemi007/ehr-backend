package main

import (
	"backend/ehr/internal/dtos"
	"backend/ehr/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

const roleType string = "user managed"

func (app *application) createEntityHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload dtos.CreateEntityRequest
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

		entity := &model.Entity{
			FullName: payload.FullName,
			Email:    payload.Email,
		}

		// hash the entity password
		if err := entity.Password.Set(payload.Password); err != nil {
			internalServerError(c, err)
			return
		}

		// generate random string (slug)
		if err := entity.EntityInformation.Slug.GenerateSlug(3, true, false); err != nil {
			internalServerError(c, err)
		}

		ctx := c.Request.Context()
		// create the user
		if err := app.store.Entity.Create(ctx, entity); err != nil {
			switch err {
			case model.ErrEntityDuplicateEmail:
				conflictResponse(c, err)
				return
			default:
				internalServerError(c, err)
				return
			}
		}

		jsonResponse(c, http.StatusCreated, "Entity Created")
	}
}

func (app *application) getEntityByIdHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		entityId := c.Param("id")
		if entityId == "" {
			badRequestError(c, nil)
			return
		}

		ctx := c.Request.Context()
		entity, err := app.store.Entity.FindById(ctx, entityId)
		if err != nil {
			switch err {
			case model.ErrNotFound:
				notFoundError(c, err)
				return
			default:
				internalServerError(c, err)
				return
			}
		}

		response := struct {
			EntityID int64  `json:"entityId"`
			Email    string `json:"email"`
			Slug     string `json:"slug"`
		}{
			EntityID: entity.ID,
			Email:    entity.Email,
			Slug:     entity.EntityInformation.Slug.Value,
		}
		jsonResponse(c, http.StatusOK, response)
	}
}

func (app *application) createRoleHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload dtos.CreateRoleRequest
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

		role := &model.Role{
			Name:        payload.Name,
			Description: payload.Description,
			Type:        roleType,
		}

		ctx := c.Request.Context()

		if err := app.store.Entity.CreateRole(ctx, role); err != nil {
			switch err {
			case model.ErrRoleNameConflict:
				conflictResponse(c, err)
				return
			default:
				internalServerError(c, err)
				return
			}
		}

		jsonResponse(c, http.StatusOK, "role created successfully")
	}
}

func (app *application) assignPermissionToRoleHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload dtos.RolePermissionAssign
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
		if err := app.store.Entity.AssignPermissionToRole(ctx, payload.RoleName, payload.Permissions); err != nil {
			internalServerError(c, err)
			return
		}

		jsonResponse(c, http.StatusOK, "Successfully assigned permission to role")
	}
}

func (app *application) assignRoleToUserHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload dtos.RoleUserAssign
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
		if err := app.store.Entity.AssignRoleToUser(ctx, payload.Slug, payload.RoleName); err != nil {
			internalServerError(c, err)
			return
		}

		jsonResponse(c, http.StatusOK, "Successfully assigned user to role")
	}
}

func (app *application) assignPermissionToUserHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload dtos.PermissionUserAssign
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
		if err := app.store.Entity.AssignPermissionToUser(ctx, payload.Slug, payload.Permissions); err != nil {
			internalServerError(c, err)
			return
		}

		jsonResponse(c, http.StatusOK, "Successfully assigned user to permission")
	}
}

func (app *application) removePermissionFromUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload dtos.DeleteUserPermission
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
		if err := app.store.Entity.RemovePermissionFromUser(ctx, payload.Slug, payload.Permission); err != nil {
			internalServerError(c, err)
			return
		}

		jsonResponse(c, http.StatusOK, "Successfully removed user permission")
	}
}

func (app *application) removePermissionFromRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload dtos.DeleteRolePermission
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
		if err := app.store.Entity.RemovePermissionFromRole(ctx, payload.RoleName, payload.Permission); err != nil {
			internalServerError(c, err)
			return
		}

		jsonResponse(c, http.StatusOK, "Successfully removed role permission")
	}
}

func (app *application) getAllRoles() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		roles, err := app.store.Entity.GetAllRoles(ctx)
		if err != nil {
			internalServerError(c, err)
			return
		}

		var response []dtos.RoleDTO
		for _, role := range roles {
			response = append(response, dtos.RoleDTO{
				ID:          role.ID,
				Name:        role.Name,
				Description: role.Description,
				Type:        role.Type,
			})
		}

		jsonResponse(c, http.StatusOK, response)
	}
}

func (app *application) getAllPermissions() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		permissions, err := app.store.Entity.GetAllPermissions(ctx)
		if err != nil {
			internalServerError(c, err)
			return
		}

		var response []dtos.PermissionDTO
		for _, permission := range permissions {
			response = append(response, dtos.PermissionDTO{
				ID:          permission.ID,
				Name:        permission.Name,
				Description: permission.Description,
				Type:        permission.Type,
			})
		}

		jsonResponse(c, http.StatusOK, response)
	}
}
