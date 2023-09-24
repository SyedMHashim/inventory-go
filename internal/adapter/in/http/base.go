package http

import (
	"encoding/json"
	"log"
	"net/http"

	"example.com/apiserver/internal/application/domain"
)

type responsePayload map[string]any

type responseData map[string]any

func sendResponse(w http.ResponseWriter, statusCode int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(response)
}

func sendError(w http.ResponseWriter, statusCode int, errResp responsePayload) {
	sendResponse(w, statusCode, errResp)
}

func newSuccessfulResponse(message string, data interface{}) responsePayload {
	return responsePayload{
		"message": message,
		"data":    data,
	}
}

func newErrorResponse(message, details string) responsePayload {
	return responsePayload{
		"error":   message,
		"details": details,
	}
}

func errorResponseOf(err error) responsePayload {
	if domainErr, ok := err.(domain.Error); ok {
		return responsePayload{
			"error":   domainErr.Message(),
			"details": domainErr.Details(),
		}
	}

	log.Println("Uncategorized error happened, returning default error", err)
	defaultNonSensitiveErrorResponse := responsePayload{
		"error":   "Something wrong happened at our end",
		"details": "Please try again after a while. If the issue keeps happening, Please open an issue at https://github.com/SyedMHashim/inventory-go/issues",
	}
	return defaultNonSensitiveErrorResponse
}

func statusCodeOf(err error) int {
	domainErr, ok := err.(domain.Error)
	if !ok {
		return http.StatusInternalServerError
	}

	switch domainErr.Category() {
	case domain.InvalidInput:
		return http.StatusBadRequest
	case domain.Unauthorized:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}
