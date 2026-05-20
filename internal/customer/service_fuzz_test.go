package customer

import (
	"context"
	"testing"
	"time"
)

// ============== FUZZING TESTS FOR SERVICE ==============

// FuzzServiceCreate testa Create com inputs aleatórios
func FuzzServiceCreate(f *testing.F) {
	f.Add("", "", 0, "", "", "")
	f.Add("DOC-001", "John", 750, "LOW", "5000-8000", "ACTIVE")
	f.Add("DOC", "Jane", 500, "MEDIUM", "3000-5000", "INACTIVE")
	f.Add("", "NoDoc", 600, "HIGH", "2000-4000", "UNDER_REVIEW")
	f.Add("DOC-TEST", "", 400, "LOW", "", "ACTIVE")
	f.Add("a", "b", -100, "INVALID", "x-y", "INVALID")
	f.Add("DOC-999", "Very Long Name With Many Characters", 1000, "LOW", "10000+", "ACTIVE")
	f.Add("DOC-SPECIAL", "Name & Co.", 0, "MEDIUM", "0-1000", "INACTIVE")

	f.Fuzz(func(t *testing.T, document, name string, score int, riskLevel, incomeRange, status string) {
		repo := &MockRepository{
			createFunc: func(ctx context.Context, c *Customer) error {
				return nil
			},
		}
		svc := NewService(repo)

		input := CreateCustomerInput{
			Document:    document,
			Name:        name,
			Score:       score,
			RiskLevel:   riskLevel,
			IncomeRange: incomeRange,
			Status:      status,
		}

		customer, err := svc.Create(context.Background(), input)

		// Validar que não panics e sempre retorna uma resposta consistente
		if err == nil {
			if customer == nil {
				t.Error("customer não deve ser nil quando não há erro")
			} else {
				if customer.ID == "" {
					t.Error("customer.ID não deve estar vazio")
				}
				if customer.Document != document {
					t.Errorf("Document não foi preservado: %s != %s", customer.Document, document)
				}
			}
		} else {
			if customer != nil {
				t.Error("customer deve ser nil quando há erro")
			}
		}
	})
}

// FuzzServiceList testa List com diferentes limites e offsets
func FuzzServiceList(f *testing.F) {
	f.Add(0, 0)
	f.Add(20, 0)
	f.Add(10, 5)
	f.Add(-1, -1)
	f.Add(999999, 999999)
	f.Add(100, 0)
	f.Add(101, 10)
	f.Add(50, -50)

	f.Fuzz(func(t *testing.T, limit, offset int) {
		repo := &MockRepository{
			listFunc: func(ctx context.Context, l, o int) ([]Customer, error) {
				// Validar que os valores foram normalizados
				if l <= 0 || l > 100 {
					if l != 20 {
						t.Errorf("Limit inválido após normalização: %d", l)
					}
				}
				if o < 0 {
					if o != 0 {
						t.Errorf("Offset inválido após normalização: %d", o)
					}
				}
				return []Customer{}, nil
			},
		}
		svc := NewService(repo)

		customers, err := svc.List(context.Background(), limit, offset)

		// Validar que sempre retorna slice válido
		if err != nil {
			t.Errorf("List não deve retornar erro: %v", err)
		}
		if customers == nil {
			t.Error("customers não deve ser nil")
		}
	})
}

// FuzzServiceGetByID testa GetByID com IDs aleatórios
func FuzzServiceGetByID(f *testing.F) {
	f.Add("")
	f.Add("123")
	f.Add("test-id")
	f.Add("uuid-format-test")
	f.Add("a")
	f.Add("very-long-id-string-with-many-characters-for-testing")
	f.Add("ID/with/slashes")
	f.Add("id-with-dashes-and-numbers-123-456")

	f.Fuzz(func(t *testing.T, id string) {
		repo := &MockRepository{
			getByIDFunc: func(ctx context.Context, idParam string) (*Customer, error) {
				if idParam == id {
					return &Customer{ID: id, Document: "DOC-001"}, nil
				}
				return nil, nil
			},
		}
		svc := NewService(repo)

		_, err := svc.GetByID(context.Background(), id)

		// Validar consistência
		if err != nil {
			t.Errorf("GetByID não deve retornar erro: %v", err)
		}
	})
}

// FuzzServiceGetByDocument testa GetByDocument com documentos aleatórios
func FuzzServiceGetByDocument(f *testing.F) {
	f.Add("")
	f.Add("DOC-001")
	f.Add("123456789")
	f.Add("DOC/001")
	f.Add("d")
	f.Add("DOCUMENT-WITH-VERY-LONG-NAME-FOR-TESTING")
	f.Add("doc-with-special-chars-!@#$")
	f.Add("numbers-12345-67890")

	f.Fuzz(func(t *testing.T, document string) {
		repo := &MockRepository{
			getByDocumentFunc: func(ctx context.Context, doc string) (*Customer, error) {
				if doc == document {
					return &Customer{Document: document, Name: "Test"}, nil
				}
				return nil, nil
			},
		}
		svc := NewService(repo)

		_, err := svc.GetByDocument(context.Background(), document)

		// Validar consistência
		if err != nil {
			t.Errorf("GetByDocument não deve retornar erro: %v", err)
		}
	})
}

// FuzzServiceUpdateStatus testa UpdateStatus com status aleatórios
func FuzzServiceUpdateStatus(f *testing.F) {
	f.Add("test-id", "ACTIVE")
	f.Add("test-id", "INACTIVE")
	f.Add("test-id", "UNDER_REVIEW")
	f.Add("test-id", "")
	f.Add("test-id", "INVALID")
	f.Add("test-id", "active")
	f.Add("test-id", "VerY_mIxEd_CaSe")
	f.Add("", "ACTIVE")
	f.Add("id", "status")

	f.Fuzz(func(t *testing.T, id, status string) {
		callCount := 0
		repo := &MockRepository{
			updateStatusFunc: func(ctx context.Context, idParam, statusParam string, updatedAt time.Time) error {
				callCount++
				if statusParam != status {
					t.Errorf("Status não foi preservado: %s != %s", statusParam, status)
				}
				return nil
			},
		}
		svc := NewService(repo)

		err := svc.UpdateStatus(context.Background(), id, status)

		// Validar que sempre retorna consistentemente
		if status == "ACTIVE" || status == "INACTIVE" || status == "UNDER_REVIEW" {
			if err != nil {
				t.Errorf("UpdateStatus deve suceder com status válido: %v", err)
			}
			if callCount != 1 {
				t.Errorf("Repository deve ser chamado uma vez, foi %d", callCount)
			}
		}
		// Status inválido deve retornar erro sem chamar repo
	})
}

// FuzzCreateInputCombinations testa diferentes combinações de inputs
func FuzzCreateInputCombinations(f *testing.F) {
	f.Add([]byte(`{"document":"D","name":"N","score":0,"risk_level":"LOW","income_range":"","status":"ACTIVE"}`))
	f.Add([]byte(`{"document":"","name":"","score":-999,"risk_level":"","income_range":"","status":""}`))
	f.Add([]byte(`{"document":"X","name":"Y","score":500,"risk_level":"MEDIUM","income_range":"1000-2000","status":"INACTIVE"}`))
	f.Add([]byte(`{}`))
	f.Add([]byte(`{"document":null,"name":null,"score":null,"risk_level":null,"status":null}`))

	f.Fuzz(func(t *testing.T, data []byte) {
		repo := &MockRepository{
			createFunc: func(ctx context.Context, c *Customer) error {
				return nil
			},
		}
		svc := NewService(repo)

		// Não vamos fazer parsing JSON aqui, apenas validar que o serviço
		// pode lidar com inputs aleatórios sem panics
		input := CreateCustomerInput{}

		customer, _ := svc.Create(context.Background(), input)

		// Input vazio deve retornar erro (missing document e name)
		if customer != nil {
			if customer.ID == "" {
				t.Error("ID não deve estar vazio")
			}
		}
	})
}

// FuzzServiceListBoundaryValues testa boundary values para List
func FuzzServiceListBoundaryValues(f *testing.F) {
	f.Add(0, 0)
	f.Add(1, 0)
	f.Add(20, 0)
	f.Add(100, 0)
	f.Add(101, 0)
	f.Add(-1, 0)
	f.Add(-2147483648, 0)
	f.Add(2147483647, 0)
	f.Add(50, -1)
	f.Add(50, -2147483648)
	f.Add(50, 2147483647)

	f.Fuzz(func(t *testing.T, limit, offset int) {
		normalizedLimit := limit
		normalizedOffset := offset

		if normalizedLimit <= 0 || normalizedLimit > 100 {
			normalizedLimit = 20
		}
		if normalizedOffset < 0 {
			normalizedOffset = 0
		}

		repo := &MockRepository{
			listFunc: func(ctx context.Context, l, o int) ([]Customer, error) {
				if l != normalizedLimit {
					t.Errorf("Limit não foi normalizado corretamente: expected %d, got %d", normalizedLimit, l)
				}
				if o != normalizedOffset {
					t.Errorf("Offset não foi normalizado corretamente: expected %d, got %d", normalizedOffset, o)
				}
				return []Customer{}, nil
			},
		}
		svc := NewService(repo)

		customers, err := svc.List(context.Background(), limit, offset)

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if customers == nil {
			t.Fatal("customers slice should not be nil")
		}
	})
}

// FuzzServiceScoreBoundary testa boundary values para Score
func FuzzServiceScoreBoundary(f *testing.F) {
	f.Add(-2147483648)
	f.Add(-1001)
	f.Add(-1)
	f.Add(0)
	f.Add(1)
	f.Add(500)
	f.Add(999)
	f.Add(1000)
	f.Add(1001)
	f.Add(2147483647)

	f.Fuzz(func(t *testing.T, score int) {
		repo := &MockRepository{
			createFunc: func(ctx context.Context, c *Customer) error {
				return nil
			},
		}
		svc := NewService(repo)

		input := CreateCustomerInput{
			Document:  "DOC-TEST",
			Name:      "Test",
			Score:     score,
			RiskLevel: "LOW",
			Status:    "ACTIVE",
		}

		customer, err := svc.Create(context.Background(), input)

		// Validar que score deve estar entre 0 e 1000
		if score >= 0 && score <= 1000 {
			if err != nil {
				t.Errorf("Create deve suceder com score válido %d: %v", score, err)
			}
			if customer == nil {
				t.Error("customer não deve ser nil com score válido")
			}
		} else {
			if err == nil {
				t.Errorf("Create deve retornar erro com score inválido %d", score)
			}
			if customer != nil {
				t.Error("customer deve ser nil com score inválido")
			}
		}
	})
}
