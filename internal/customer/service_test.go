package customer

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/GuilhermeRossiKirsten/CustomerRegistryAPI/internal/error_handler"
	"github.com/stretchr/testify/assert"
)

// MockRepository é um mock do Repository para testes
type MockRepository struct {
	createFunc        func(ctx context.Context, c *Customer) error
	listFunc          func(ctx context.Context, limit, offset int) ([]Customer, error)
	getByIDFunc       func(ctx context.Context, id string) (*Customer, error)
	getByDocumentFunc func(ctx context.Context, doc string) (*Customer, error)
	updateStatusFunc  func(ctx context.Context, id, status string, updatedAt time.Time) error
}

func (m *MockRepository) Create(ctx context.Context, c *Customer) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, c)
	}
	return nil
}

func (m *MockRepository) List(ctx context.Context, limit, offset int) ([]Customer, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx, limit, offset)
	}
	return []Customer{}, nil
}

func (m *MockRepository) GetByID(ctx context.Context, id string) (*Customer, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, error_handler.ErrCustomerNotFound
}

func (m *MockRepository) GetByDocument(ctx context.Context, doc string) (*Customer, error) {
	if m.getByDocumentFunc != nil {
		return m.getByDocumentFunc(ctx, doc)
	}
	return nil, error_handler.ErrCustomerNotFound
}

func (m *MockRepository) UpdateStatus(ctx context.Context, id, status string, updatedAt time.Time) error {
	if m.updateStatusFunc != nil {
		return m.updateStatusFunc(ctx, id, status, updatedAt)
	}
	return nil
}

// ============== CREATE TESTS ==============

func TestCreateValidCustomer(t *testing.T) {
	repo := &MockRepository{
		createFunc: func(ctx context.Context, c *Customer) error {
			return nil
		},
	}
	svc := NewService(repo)

	input := CreateCustomerInput{
		Document:    "DOC-001",
		Name:        "John Doe",
		Score:       750,
		RiskLevel:   "LOW",
		IncomeRange: "5000-8000",
		Status:      "ACTIVE",
	}

	customer, err := svc.Create(context.Background(), input)

	assert.NoError(t, err)
	assert.NotNil(t, customer)
	assert.Equal(t, "DOC-001", customer.Document)
	assert.Equal(t, "John Doe", customer.Name)
	assert.Equal(t, 750, customer.Score)
	assert.Equal(t, "LOW", customer.RiskLevel)
	assert.Equal(t, "ACTIVE", customer.Status)
	assert.NotEmpty(t, customer.ID)
	assert.NotZero(t, customer.CreatedAt)
	assert.NotZero(t, customer.UpdatedAt)
}

func TestCreateMissingDocument(t *testing.T) {
	repo := &MockRepository{}
	svc := NewService(repo)

	input := CreateCustomerInput{
		Document:  "",
		Name:      "John Doe",
		Score:     750,
		RiskLevel: "LOW",
		Status:    "ACTIVE",
	}

	customer, err := svc.Create(context.Background(), input)

	assert.Error(t, err)
	assert.Nil(t, customer)
	assert.True(t, errors.Is(err, error_handler.ErrMissingDocument))
}

func TestCreateMissingName(t *testing.T) {
	repo := &MockRepository{}
	svc := NewService(repo)

	input := CreateCustomerInput{
		Document:  "DOC-001",
		Name:      "",
		Score:     750,
		RiskLevel: "LOW",
		Status:    "ACTIVE",
	}

	customer, err := svc.Create(context.Background(), input)

	assert.Error(t, err)
	assert.Nil(t, customer)
	assert.True(t, errors.Is(err, error_handler.ErrMissingName))
}

func TestCreateInvalidScoreNegative(t *testing.T) {
	repo := &MockRepository{}
	svc := NewService(repo)

	input := CreateCustomerInput{
		Document:  "DOC-001",
		Name:      "John Doe",
		Score:     -1,
		RiskLevel: "LOW",
		Status:    "ACTIVE",
	}

	customer, err := svc.Create(context.Background(), input)

	assert.Error(t, err)
	assert.Nil(t, customer)
	assert.True(t, errors.Is(err, error_handler.ErrInvalidScoreRange))
}

func TestCreateInvalidScoreTooHigh(t *testing.T) {
	repo := &MockRepository{}
	svc := NewService(repo)

	input := CreateCustomerInput{
		Document:  "DOC-001",
		Name:      "John Doe",
		Score:     1001,
		RiskLevel: "LOW",
		Status:    "ACTIVE",
	}

	customer, err := svc.Create(context.Background(), input)

	assert.Error(t, err)
	assert.Nil(t, customer)
	assert.True(t, errors.Is(err, error_handler.ErrInvalidScoreRange))
}

func TestCreateInvalidRiskLevel(t *testing.T) {
	repo := &MockRepository{}
	svc := NewService(repo)

	input := CreateCustomerInput{
		Document:  "DOC-001",
		Name:      "John Doe",
		Score:     750,
		RiskLevel: "INVALID",
		Status:    "ACTIVE",
	}

	customer, err := svc.Create(context.Background(), input)

	assert.Error(t, err)
	assert.Nil(t, customer)
	assert.True(t, errors.Is(err, error_handler.ErrInvalidRiskLevel))
}

func TestCreateInvalidStatus(t *testing.T) {
	repo := &MockRepository{}
	svc := NewService(repo)

	input := CreateCustomerInput{
		Document:  "DOC-001",
		Name:      "John Doe",
		Score:     750,
		RiskLevel: "LOW",
		Status:    "INVALID",
	}

	customer, err := svc.Create(context.Background(), input)

	assert.Error(t, err)
	assert.Nil(t, customer)
	assert.True(t, errors.Is(err, error_handler.ErrInvalidStatus))
}

func TestCreateRepositoryError(t *testing.T) {
	repo := &MockRepository{
		createFunc: func(ctx context.Context, c *Customer) error {
			return error_handler.ErrDuplicateDocument
		},
	}
	svc := NewService(repo)

	input := CreateCustomerInput{
		Document:  "DOC-001",
		Name:      "John Doe",
		Score:     750,
		RiskLevel: "LOW",
		Status:    "ACTIVE",
	}

	customer, err := svc.Create(context.Background(), input)

	assert.Error(t, err)
	assert.Nil(t, customer)
	assert.True(t, errors.Is(err, error_handler.ErrDuplicateDocument))
}

func TestCreateWithAllRiskLevels(t *testing.T) {
	testCases := []struct {
		riskLevel string
	}{
		{"LOW"},
		{"MEDIUM"},
		{"HIGH"},
	}

	for _, tc := range testCases {
		t.Run(tc.riskLevel, func(t *testing.T) {
			repo := &MockRepository{
				createFunc: func(ctx context.Context, c *Customer) error {
					return nil
				},
			}
			svc := NewService(repo)

			input := CreateCustomerInput{
				Document:  "DOC-001",
				Name:      "John Doe",
				Score:     750,
				RiskLevel: tc.riskLevel,
				Status:    "ACTIVE",
			}

			customer, err := svc.Create(context.Background(), input)

			assert.NoError(t, err)
			assert.NotNil(t, customer)
			assert.Equal(t, tc.riskLevel, customer.RiskLevel)
		})
	}
}

func TestCreateWithAllStatuses(t *testing.T) {
	testCases := []struct {
		status string
	}{
		{"ACTIVE"},
		{"INACTIVE"},
		{"UNDER_REVIEW"},
	}

	for _, tc := range testCases {
		t.Run(tc.status, func(t *testing.T) {
			repo := &MockRepository{
				createFunc: func(ctx context.Context, c *Customer) error {
					return nil
				},
			}
			svc := NewService(repo)

			input := CreateCustomerInput{
				Document:  "DOC-001",
				Name:      "John Doe",
				Score:     750,
				RiskLevel: "LOW",
				Status:    tc.status,
			}

			customer, err := svc.Create(context.Background(), input)

			assert.NoError(t, err)
			assert.NotNil(t, customer)
			assert.Equal(t, tc.status, customer.Status)
		})
	}
}

// ============== LIST TESTS ==============

func TestListWithDefaultLimit(t *testing.T) {
	expectedCustomers := []Customer{
		{ID: "1", Document: "DOC-001", Name: "John"},
		{ID: "2", Document: "DOC-002", Name: "Jane"},
	}

	repo := &MockRepository{
		listFunc: func(ctx context.Context, limit, offset int) ([]Customer, error) {
			assert.Equal(t, 20, limit)
			assert.Equal(t, 0, offset)
			return expectedCustomers, nil
		},
	}
	svc := NewService(repo)

	customers, err := svc.List(context.Background(), 0, 0)

	assert.NoError(t, err)
	assert.Equal(t, expectedCustomers, customers)
}

func TestListWithCustomLimit(t *testing.T) {
	expectedCustomers := []Customer{
		{ID: "1", Document: "DOC-001", Name: "John"},
	}

	repo := &MockRepository{
		listFunc: func(ctx context.Context, limit, offset int) ([]Customer, error) {
			assert.Equal(t, 10, limit)
			assert.Equal(t, 5, offset)
			return expectedCustomers, nil
		},
	}
	svc := NewService(repo)

	customers, err := svc.List(context.Background(), 10, 5)

	assert.NoError(t, err)
	assert.Equal(t, expectedCustomers, customers)
}

func TestListWithNegativeLimit(t *testing.T) {
	repo := &MockRepository{
		listFunc: func(ctx context.Context, limit, offset int) ([]Customer, error) {
			assert.Equal(t, 20, limit)
			return []Customer{}, nil
		},
	}
	svc := NewService(repo)

	customers, err := svc.List(context.Background(), -1, 0)

	assert.NoError(t, err)
	assert.NotNil(t, customers)
}

func TestListWithLimitGreaterThan100(t *testing.T) {
	repo := &MockRepository{
		listFunc: func(ctx context.Context, limit, offset int) ([]Customer, error) {
			assert.Equal(t, 20, limit)
			return []Customer{}, nil
		},
	}
	svc := NewService(repo)

	customers, err := svc.List(context.Background(), 150, 0)

	assert.NoError(t, err)
	assert.NotNil(t, customers)
}

func TestListWithNegativeOffset(t *testing.T) {
	repo := &MockRepository{
		listFunc: func(ctx context.Context, limit, offset int) ([]Customer, error) {
			assert.Equal(t, 20, limit)
			assert.Equal(t, 0, offset)
			return []Customer{}, nil
		},
	}
	svc := NewService(repo)

	customers, err := svc.List(context.Background(), 20, -5)

	assert.NoError(t, err)
	assert.NotNil(t, customers)
}

func TestListRepositoryError(t *testing.T) {
	repo := &MockRepository{
		listFunc: func(ctx context.Context, limit, offset int) ([]Customer, error) {
			return nil, errors.New("database error")
		},
	}
	svc := NewService(repo)

	customers, err := svc.List(context.Background(), 20, 0)

	assert.Error(t, err)
	assert.Nil(t, customers)
}

// ============== GET BY ID TESTS ==============

func TestGetByIDSuccess(t *testing.T) {
	expectedCustomer := &Customer{
		ID:       "test-id",
		Document: "DOC-001",
		Name:     "John Doe",
	}

	repo := &MockRepository{
		getByIDFunc: func(ctx context.Context, id string) (*Customer, error) {
			assert.Equal(t, "test-id", id)
			return expectedCustomer, nil
		},
	}
	svc := NewService(repo)

	customer, err := svc.GetByID(context.Background(), "test-id")

	assert.NoError(t, err)
	assert.Equal(t, expectedCustomer, customer)
}

func TestGetByIDNotFound(t *testing.T) {
	repo := &MockRepository{
		getByIDFunc: func(ctx context.Context, id string) (*Customer, error) {
			return nil, error_handler.ErrCustomerNotFound
		},
	}
	svc := NewService(repo)

	customer, err := svc.GetByID(context.Background(), "invalid-id")

	assert.Error(t, err)
	assert.Nil(t, customer)
	assert.True(t, errors.Is(err, error_handler.ErrCustomerNotFound))
}

// ============== GET BY DOCUMENT TESTS ==============

func TestGetByDocumentSuccess(t *testing.T) {
	expectedCustomer := &Customer{
		ID:       "test-id",
		Document: "DOC-001",
		Name:     "John Doe",
	}

	repo := &MockRepository{
		getByDocumentFunc: func(ctx context.Context, doc string) (*Customer, error) {
			assert.Equal(t, "DOC-001", doc)
			return expectedCustomer, nil
		},
	}
	svc := NewService(repo)

	customer, err := svc.GetByDocument(context.Background(), "DOC-001")

	assert.NoError(t, err)
	assert.Equal(t, expectedCustomer, customer)
}

func TestGetByDocumentNotFound(t *testing.T) {
	repo := &MockRepository{
		getByDocumentFunc: func(ctx context.Context, doc string) (*Customer, error) {
			return nil, error_handler.ErrCustomerNotFound
		},
	}
	svc := NewService(repo)

	customer, err := svc.GetByDocument(context.Background(), "INVALID")

	assert.Error(t, err)
	assert.Nil(t, customer)
	assert.True(t, errors.Is(err, error_handler.ErrCustomerNotFound))
}

// ============== UPDATE STATUS TESTS ==============

func TestUpdateStatusSuccess(t *testing.T) {
	repo := &MockRepository{
		updateStatusFunc: func(ctx context.Context, id, status string, updatedAt time.Time) error {
			assert.Equal(t, "test-id", id)
			assert.Equal(t, "ACTIVE", status)
			return nil
		},
	}
	svc := NewService(repo)

	err := svc.UpdateStatus(context.Background(), "test-id", "ACTIVE")

	assert.NoError(t, err)
}

func TestUpdateStatusWithAllStatuses(t *testing.T) {
	testCases := []struct {
		status string
	}{
		{"ACTIVE"},
		{"INACTIVE"},
		{"UNDER_REVIEW"},
	}

	for _, tc := range testCases {
		t.Run(tc.status, func(t *testing.T) {
			repo := &MockRepository{
				updateStatusFunc: func(ctx context.Context, id, status string, updatedAt time.Time) error {
					assert.Equal(t, tc.status, status)
					return nil
				},
			}
			svc := NewService(repo)

			err := svc.UpdateStatus(context.Background(), "test-id", tc.status)

			assert.NoError(t, err)
		})
	}
}

func TestUpdateStatusInvalidStatus(t *testing.T) {
	repo := &MockRepository{}
	svc := NewService(repo)

	err := svc.UpdateStatus(context.Background(), "test-id", "INVALID")

	assert.Error(t, err)
	assert.True(t, errors.Is(err, error_handler.ErrInvalidStatus))
}

func TestUpdateStatusRepositoryError(t *testing.T) {
	repo := &MockRepository{
		updateStatusFunc: func(ctx context.Context, id, status string, updatedAt time.Time) error {
			return error_handler.ErrCustomerNotFound
		},
	}
	svc := NewService(repo)

	err := svc.UpdateStatus(context.Background(), "invalid-id", "ACTIVE")

	assert.Error(t, err)
	assert.True(t, errors.Is(err, error_handler.ErrCustomerNotFound))
}

// ============== BOUNDARY TESTS ==============

func TestCreateWithBoundaryScores(t *testing.T) {
	testCases := []struct {
		score int
		valid bool
	}{
		{0, true},
		{1, true},
		{500, true},
		{999, true},
		{1000, true},
		{-1, false},
		{1001, false},
	}

	for _, tc := range testCases {
		t.Run("score_"+string(rune(tc.score)), func(t *testing.T) {
			repo := &MockRepository{
				createFunc: func(ctx context.Context, c *Customer) error {
					return nil
				},
			}
			svc := NewService(repo)

			input := CreateCustomerInput{
				Document:  "DOC-001",
				Name:      "John",
				Score:     tc.score,
				RiskLevel: "LOW",
				Status:    "ACTIVE",
			}

			customer, err := svc.Create(context.Background(), input)

			if tc.valid {
				assert.NoError(t, err)
				assert.NotNil(t, customer)
			} else {
				assert.Error(t, err)
				assert.Nil(t, customer)
			}
		})
	}
}
