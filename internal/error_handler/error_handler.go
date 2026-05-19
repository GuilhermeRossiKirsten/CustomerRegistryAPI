package error_handler

import "errors"

var (
	ErrCustomerNotFound      = errors.New("customer not found")
	ErrCustomerAlreadyExists = errors.New("customer already exists")
	ErrInvalidStatus         = errors.New("invalid status")
	ErrInvalidRiskLevel      = errors.New("invalid risk level")
	ErrInvalidScore          = errors.New("invalid score")
	ErrInvalidIncomeRange    = errors.New("invalid income range")
	ErrDuplicateDocument     = errors.New("duplicate document")
	ErrInvalidScoreRange     = errors.New("score must be between 0 and 1000")
	ErrMissingDocument       = errors.New("document is required")
	ErrMissingName           = errors.New("name is required")
)
