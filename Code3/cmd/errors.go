package main

import (
	"backend/ehr/internal/logging"
	"net/http"

	"github.com/gin-gonic/gin"
)

// error method to send internal server error message
func internalServerError(c *gin.Context, err error) {
	logging.Logger.Errorw("internal server error", "method", c.Request.Method, "path", c.Request.URL.Path, "error", err.Error())
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": "the server encountered a problem",
		"code":  500,
	})
}

// error method to send bad request server error message
func badRequestError(c *gin.Context, err error) {
	logging.Logger.Warnf("bad request error", "method", c.Request.Method, "path", c.Request.URL.Path, "error", err.Error())
	c.JSON(http.StatusBadRequest, gin.H{
		"error": err.Error(),
		"code":  400,
	})
}

// error method to send not found request error message
func notFoundError(c *gin.Context, err error) {
	logging.Logger.Errorw("not found error", "method", c.Request.Method, "path", c.Request.URL.Path, "error", err.Error())
	c.JSON(http.StatusNotFound, gin.H{
		"error": "record not found",
		"code":  404,
	})
}

// error method to send conflict error message
func conflictResponse(c *gin.Context, err error) {
	logging.Logger.Errorw("conflict error", "method", c.Request.Method, "path", c.Request.URL.Path, "error", err.Error())
	c.JSON(http.StatusConflict, gin.H{
		"error": err.Error(),
		"code":  409,
	})
}

// error method to send unauthorized error message
func unauthorizedErrorResponse(c *gin.Context, err error) {
	logging.Logger.Errorw("unauthorized error", "method", c.Request.Method, "path", c.Request.URL.Path, "error", err.Error())
	c.JSON(http.StatusUnauthorized, gin.H{
		"error": "unauthorized",
		"code":  401,
	})
}

// error method to send forbidden error message
func forbiddenErrorResponse(c *gin.Context, err error) {
	logging.Logger.Errorw("forbidden error", "method", c.Request.Method, "path", c.Request.URL.Path, "error", err.Error())
	c.JSON(http.StatusForbidden, gin.H{
		"error": "forbidden",
		"code":  403,
	})
}
