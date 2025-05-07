package main

import (
	"backend/ehr/internal/model"
	"errors"

	"github.com/gin-gonic/gin"
)

func (app *application) permit(permissionList []string, next gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := 1 // @FIXME: get the user id from the context

		// check if the user has the required permission
		_, err := app.store.Users.VerifyPermission(c.Request.Context(), int64(userId), permissionList)
		if err != nil {
			switch {
			case errors.Is(err, model.ErrNotFound):
				forbiddenErrorResponse(c, err)
				return
			default:
				internalServerError(c, err)
				return
			}
		}

		next(c)
	}
}
