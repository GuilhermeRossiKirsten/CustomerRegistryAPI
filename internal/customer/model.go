package customer

import "time"

// Customer representa um cliente cadastrado.
type Customer struct {
	// ID único em UUID v4.
	// example: 1f4f3c2a-8f57-4d3d-9b3e-4a2c4e2b5a1d
	ID string `json:"id"`
	// Documento único do cliente.
	// example: FAKE-00001
	Document string `json:"document"`
	// Nome do cliente.
	// example: Cliente Simulado A
	Name string `json:"name"`
	// Score de crédito (0–1000).
	// example: 742
	Score int `json:"score"`
	// Nível de risco.
	// enum: LOW,MEDIUM,HIGH
	// example: LOW
	RiskLevel string `json:"risk_level"`
	// Faixa de renda.
	// example: 3000-5000
	IncomeRange string `json:"income_range"`
	// Status atual.
	// enum: ACTIVE,INACTIVE,UNDER_REVIEW
	// example: ACTIVE
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateCustomerInput é o payload de criação.
type CreateCustomerInput struct {
	// required: true
	// example: FAKE-00001
	Document string `json:"document"`
	// required: true
	// example: Cliente Simulado A
	Name string `json:"name"`
	// required: true
	// minimum: 0
	// maximum: 1000
	// example: 742
	Score int `json:"score"`
	// required: true
	// enum: LOW,MEDIUM,HIGH
	// example: LOW
	RiskLevel string `json:"risk_level"`
	// example: 3000-5000
	IncomeRange string `json:"income_range"`
	// required: true
	// enum: ACTIVE,INACTIVE,UNDER_REVIEW
	// example: ACTIVE
	Status string `json:"status"`
}

// UpdateStatusInput é o payload de atualização de status.
type UpdateStatusInput struct {
	// required: true
	// enum: ACTIVE,INACTIVE,UNDER_REVIEW
	// example: UNDER_REVIEW
	Status string `json:"status"`
}

// ErrorResponse é o envelope padrão de erro.
type ErrorResponse struct {
	// Mensagem descritiva do erro.
	// example: invalid JSON payload
	Error string `json:"error"`
}
