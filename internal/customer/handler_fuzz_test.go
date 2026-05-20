package customer

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/GuilhermeRossiKirsten/CustomerRegistryAPI/internal/error_handler"
)

// MockService é um mock do ServiceInterface para testes
type MockService struct {
	createFunc        func(ctx context.Context, input CreateCustomerInput) (*Customer, error)
	listFunc          func(ctx context.Context, limit, offset int) ([]Customer, error)
	getByIDFunc       func(ctx context.Context, id string) (*Customer, error)
	getByDocumentFunc func(ctx context.Context, document string) (*Customer, error)
	updateStatusFunc  func(ctx context.Context, id, status string) error
}

func (m *MockService) Create(ctx context.Context, input CreateCustomerInput) (*Customer, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, input)
	}
	return nil, error_handler.ErrCustomerNotFound
}

func (m *MockService) List(ctx context.Context, limit, offset int) ([]Customer, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx, limit, offset)
	}
	return []Customer{}, nil
}

func (m *MockService) GetByID(ctx context.Context, id string) (*Customer, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, error_handler.ErrCustomerNotFound
}

func (m *MockService) GetByDocument(ctx context.Context, document string) (*Customer, error) {
	if m.getByDocumentFunc != nil {
		return m.getByDocumentFunc(ctx, document)
	}
	return nil, error_handler.ErrCustomerNotFound
}

func (m *MockService) UpdateStatus(ctx context.Context, id, status string) error {
	if m.updateStatusFunc != nil {
		return m.updateStatusFunc(ctx, id, status)
	}
	return nil
}

func FuzzCreateHandler(f *testing.F) {

	f.Add([]byte(`{"document":"DOC-001","name":"Test","score":750,"risk_level":"LOW","income_range":"5000-8000","status":"ACTIVE"}`))
	f.Add([]byte(`{"document":""}`))
	f.Add([]byte(`{}`))
	f.Add([]byte(`invalid json`))
	f.Add([]byte(`{"document":"X","name":"","score":-1,"risk_level":"","income_range":"","status":""}`))
	f.Add([]byte(`{"document":"DOC","name":"Test","score":9999,"risk_level":"INVALID","income_range":"123","status":"INVALID"}`))

	f.Fuzz(func(t *testing.T, jsonData []byte) {
		mockSvc := &MockService{
			createFunc: func(ctx context.Context, input CreateCustomerInput) (*Customer, error) {
				return &Customer{
					ID:          "test-id",
					Document:    input.Document,
					Name:        input.Name,
					Score:       input.Score,
					RiskLevel:   input.RiskLevel,
					IncomeRange: input.IncomeRange,
					Status:      input.Status,
				}, nil
			},
		}

		handler := NewHandler(mockSvc)

		req := httptest.NewRequest("POST", "/customers", bytes.NewReader(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.create(w, req)

		var response interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil && w.Body.Len() > 0 {
			t.Errorf("Response não é JSON válido: %v", err)
		}

		if w.Code < 100 || w.Code >= 600 {
			t.Errorf("Status code inválido: %d", w.Code)
		}
	})
}

func FuzzListHandlerQueryParams(f *testing.F) {
	f.Add("", "")
	f.Add("limit", "20")
	f.Add("limit", "-1")
	f.Add("limit", "999999")
	f.Add("limit", "abc")
	f.Add("offset", "10")
	f.Add("offset", "-100")
	f.Add("limit", "5")
	f.Add("offset", "1000000")

	f.Fuzz(func(t *testing.T, key, value string) {
		mockSvc := &MockService{
			listFunc: func(ctx context.Context, limit, offset int) ([]Customer, error) {
				return []Customer{}, nil
			},
		}

		handler := NewHandler(mockSvc)

		u := &url.URL{Path: "/customers"}
		if key != "" && value != "" {
			q := url.Values{}
			q.Set(key, value)
			u.RawQuery = q.Encode()
		}

		req := httptest.NewRequest("GET", u.String(), nil)
		w := httptest.NewRecorder()

		handler.list(w, req)

		var response interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil && w.Body.Len() > 0 {
			t.Errorf("Response não é JSON válido: %v", err)
		}

		if w.Code != http.StatusOK {
			t.Errorf("Status code esperado 200, obteve %d", w.Code)
		}
	})
}

func FuzzGetByIDHandler(f *testing.F) {
	f.Add("")
	f.Add("123")
	f.Add("test-id")
	f.Add("aaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	f.Add("abc123def456")
	f.Add("path-traversal")
	f.Add("uuid-style-id")

	f.Fuzz(func(t *testing.T, id string) {
		mockSvc := &MockService{
			getByIDFunc: func(ctx context.Context, idParam string) (*Customer, error) {
				if idParam == "test-id" {
					return &Customer{ID: "test-id", Document: "DOC-001"}, nil
				}
				return nil, error_handler.ErrCustomerNotFound
			},
		}

		handler := NewHandler(mockSvc)

		req := httptest.NewRequest("GET", "/customers/_", nil)
		req.SetPathValue("id", id)
		w := httptest.NewRecorder()

		handler.getByID(w, req)

		var response interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil && w.Body.Len() > 0 {
			t.Errorf("Response não é JSON válido: %v", err)
		}

		if (w.Code != http.StatusOK && w.Code != http.StatusNotFound) || w.Code < 100 || w.Code >= 600 {
			t.Errorf("Status code inválido: %d", w.Code)
		}
	})
}

func FuzzGetByDocumentHandler(f *testing.F) {
	f.Add("")
	f.Add("DOC-001")
	f.Add("123456789")
	f.Add("DOC/001")
	f.Add("aaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	f.Add("DOC-XSS-TEST")
	f.Add("DOC-SQL-TEST")

	f.Fuzz(func(t *testing.T, document string) {
		mockSvc := &MockService{
			getByDocumentFunc: func(ctx context.Context, doc string) (*Customer, error) {
				if doc == "DOC-001" {
					return &Customer{Document: "DOC-001", Name: "Test"}, nil
				}
				return nil, error_handler.ErrCustomerNotFound
			},
		}

		handler := NewHandler(mockSvc)

		req := httptest.NewRequest("GET", "/customers/document/_", nil)
		req.SetPathValue("document", document)
		w := httptest.NewRecorder()

		handler.getByDocument(w, req)

		var response interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil && w.Body.Len() > 0 {
			t.Errorf("Response não é JSON válido: %v", err)
		}

		if (w.Code != http.StatusOK && w.Code != http.StatusNotFound) || w.Code < 100 || w.Code >= 600 {
			t.Errorf("Status code inválido: %d", w.Code)
		}
	})
}

func FuzzUpdateStatusHandler(f *testing.F) {
	f.Add([]byte(`{"status":"ACTIVE"}`))
	f.Add([]byte(`{"status":"INACTIVE"}`))
	f.Add([]byte(`{"status":"UNDER_REVIEW"}`))
	f.Add([]byte(`{"status":""}`))
	f.Add([]byte(`{"status":"invalid"}`))
	f.Add([]byte(`{}`))
	f.Add([]byte(`invalid`))
	f.Add([]byte(`{"status":123}`))
	f.Add([]byte(`{"status":null}`))

	f.Fuzz(func(t *testing.T, jsonData []byte) {
		mockSvc := &MockService{
			updateStatusFunc: func(ctx context.Context, id, status string) error {
				return nil
			},
		}

		handler := NewHandler(mockSvc)

		req := httptest.NewRequest("PATCH", "/customers/test-id/status", bytes.NewReader(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.SetPathValue("id", "test-id")
		w := httptest.NewRecorder()

		handler.updateStatus(w, req)

		if w.Code != http.StatusNoContent {
			var response interface{}
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil && w.Body.Len() > 0 {
				t.Errorf("Response não é JSON válido: %v", err)
			}
		}

		if w.Code < 100 || w.Code >= 600 {
			t.Errorf("Status code inválido: %d", w.Code)
		}
	})
}

func FuzzCreatePayloadStructure(f *testing.F) {
	f.Add([]byte(`{"document":"DOC001","name":"John","score":500,"risk_level":"MEDIUM","income_range":"2000-5000","status":"ACTIVE"}`))
	f.Add([]byte(`{"document":"DOC002","name":"Jane Doe","score":900,"risk_level":"LOW","income_range":"10000+","status":"ACTIVE"}`))
	f.Add([]byte(`{"document":"","name":"NoDoc","score":0,"risk_level":"HIGH","income_range":"","status":"INACTIVE"}`))
	f.Add([]byte(`{"name":"NoDocument","score":100}`))
	f.Add([]byte(`{"document":"DOC003"}`))

	f.Fuzz(func(t *testing.T, jsonData []byte) {
		mockSvc := &MockService{
			createFunc: func(ctx context.Context, input CreateCustomerInput) (*Customer, error) {
				return &Customer{
					ID:       "fuzz-id",
					Document: input.Document,
					Name:     input.Name,
				}, nil
			},
		}

		handler := NewHandler(mockSvc)

		req := httptest.NewRequest("POST", "/customers", bytes.NewReader(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.create(w, req)

		if w == nil {
			t.Fatal("Handler retornou nil ResponseWriter")
		}

		contentType := w.Header().Get("Content-Type")
		if w.Code == http.StatusCreated || w.Code == http.StatusBadRequest {
			if contentType != "application/json" {
				t.Errorf("Content-Type esperado application/json, obteve %s", contentType)
			}
		}
	})
}
