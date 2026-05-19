package customer

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/GuilhermeRossiKirsten/CustomerRegistryAPI/internal/error_handler"
)

func respond(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data == nil {
		return
	}
	_ = json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respond(w, status, ErrorResponse{Error: message})
}

func handleServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, error_handler.ErrCustomerNotFound):
		respondError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, error_handler.ErrDuplicateDocument):
		respondError(w, http.StatusConflict, err.Error())
	case errors.Is(err, error_handler.ErrInvalidScoreRange),
		errors.Is(err, error_handler.ErrInvalidRiskLevel),
		errors.Is(err, error_handler.ErrInvalidStatus),
		errors.Is(err, error_handler.ErrMissingDocument),
		errors.Is(err, error_handler.ErrMissingName):
		respondError(w, http.StatusBadRequest, err.Error())
	default:
		respondError(w, http.StatusInternalServerError, "unexpected error")
	}
}
