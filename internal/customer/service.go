package customer

import (
	"context"
	"time"

	"github.com/GuilhermeRossiKirsten/CustomerRegistryAPI/internal/error_handler"
	"github.com/google/uuid"
)

var (
	validRiskLevels = map[string]bool{"LOW": true, "MEDIUM": true, "HIGH": true}
	validStatuses   = map[string]bool{"ACTIVE": true, "INACTIVE": true, "UNDER_REVIEW": true}
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, input CreateCustomerInput) (*Customer, error) {
	if input.Document == "" {
		return nil, error_handler.ErrMissingDocument
	}
	if input.Name == "" {
		return nil, error_handler.ErrMissingName
	}
	if input.Score < 0 || input.Score > 1000 {
		return nil, error_handler.ErrInvalidScoreRange
	}
	if !validRiskLevels[input.RiskLevel] {
		return nil, error_handler.ErrInvalidRiskLevel
	}
	if !validStatuses[input.Status] {
		return nil, error_handler.ErrInvalidStatus
	}

	now := time.Now().UTC()
	c := &Customer{
		ID:          uuid.NewString(),
		Document:    input.Document,
		Name:        input.Name,
		Score:       input.Score,
		RiskLevel:   input.RiskLevel,
		IncomeRange: input.IncomeRange,
		Status:      input.Status,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := s.repo.Create(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) List(ctx context.Context, limit, offset int) ([]Customer, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.List(ctx, limit, offset)
}

func (s *Service) GetByID(ctx context.Context, id string) (*Customer, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) GetByDocument(ctx context.Context, doc string) (*Customer, error) {
	return s.repo.GetByDocument(ctx, doc)
}

func (s *Service) UpdateStatus(ctx context.Context, id, status string) error {
	if !validStatuses[status] {
		return error_handler.ErrInvalidStatus
	}
	return s.repo.UpdateStatus(ctx, id, status, time.Now().UTC())
}
