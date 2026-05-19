package customer

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /customers", h.create)
	mux.HandleFunc("GET /customers", h.list)
	mux.HandleFunc("GET /customers/{id}", h.getByID)
	mux.HandleFunc("GET /customers/document/{document}", h.getByDocument)
	mux.HandleFunc("PATCH /customers/{id}/status", h.updateStatus)
}

// @Summary Cria um novo cliente
// @Description Cria um novo cliente no sistema
// @Tags customers
// @Accept json
// @Produce json
// @Param input body CreateCustomerInput true "Customer data"
// @Success 201 {object} Customer
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /customers [post]
func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var in CreateCustomerInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		respondError(w, http.StatusBadRequest, "invalid JSON payload")
		return
	}
	c, err := h.svc.Create(r.Context(), in)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	respond(w, http.StatusCreated, c)
}

// @Summary Lista clientes
// @Description Lista clientes com paginação
// @Tags customers
// @Accept json
// @Produce json
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {array} Customer
// @Router /customers [get]
func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	list, err := h.svc.List(r.Context(), limit, offset)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "unexpected error")
		return
	}
	if list == nil {
		list = []Customer{}
	}
	respond(w, http.StatusOK, list)
}

// @Summary Busca cliente por ID
// @Description Busca um cliente pelo ID
// @Tags customers
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Success 200 {object} Customer
// @Failure 404 {object} ErrorResponse
// @Router /customers/{id} [get]
func (h *Handler) getByID(w http.ResponseWriter, r *http.Request) {
	c, err := h.svc.GetByID(r.Context(), r.PathValue("id"))
	if err != nil {
		handleServiceError(w, err)
		return
	}
	respond(w, http.StatusOK, c)
}

// @Summary Busca cliente por documento
// @Description Busca um cliente pelo documento
// @Tags customers
// @Accept json
// @Produce json
// @Param document path string true "Customer document"
// @Success 200 {object} Customer
// @Failure 404 {object} ErrorResponse
// @Router /customers/document/{document} [get]
func (h *Handler) getByDocument(w http.ResponseWriter, r *http.Request) {
	c, err := h.svc.GetByDocument(r.Context(), r.PathValue("document"))
	if err != nil {
		handleServiceError(w, err)
		return
	}
	respond(w, http.StatusOK, c)
}

// @Summary Atualiza status do cliente
// @Description Atualiza apenas o status de um cliente
// @Tags customers
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Param input body UpdateStatusInput true "New status"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /customers/{id}/status [patch]
func (h *Handler) updateStatus(w http.ResponseWriter, r *http.Request) {
	var in UpdateStatusInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		respondError(w, http.StatusBadRequest, "invalid JSON payload")
		return
	}
	err := h.svc.UpdateStatus(r.Context(), r.PathValue("id"), strings.ToUpper(in.Status))
	if err != nil {
		handleServiceError(w, err)
		return
	}
	respond(w, http.StatusNoContent, nil)
}
