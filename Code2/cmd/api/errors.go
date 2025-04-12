package main

import (
	"backend/ehr/internal/logging"
	"net/http"
)

// error method to send internal server error message
func internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	logging.Logger.Errorw("internal server error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJSONError(w, http.StatusInternalServerError, "the server encountered a problem", 500)
}

// error method to send bad request server error message
func badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	logging.Logger.Warnf("bad request error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJSONError(w, http.StatusBadRequest, err.Error(), 400)
}

// error method to send not found request error message
func notFoundError(w http.ResponseWriter, r *http.Request, err error) {
	logging.Logger.Errorw("not found error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJSONError(w, http.StatusNotFound, "record not found", 404)
}

// error method to send conflict error message
func conflictResponse(w http.ResponseWriter, r *http.Request, err error) {
	logging.Logger.Errorw("conflict error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJSONError(w, http.StatusConflict, err.Error(), 409)
}
