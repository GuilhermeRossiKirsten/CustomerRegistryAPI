package customer

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/GuilhermeRossiKirsten/CustomerRegistryAPI/internal/error_handler"
	"github.com/jackc/pgx/v5/pgconn"
)

type Repository interface {
	Create(ctx context.Context, customer *Customer) error
	List(ctx context.Context, limit, offset int) ([]Customer, error)
	GetByID(ctx context.Context, id string) (*Customer, error)
	GetByDocument(ctx context.Context, document string) (*Customer, error)
	UpdateStatus(ctx context.Context, id, status string, updatedAt time.Time) error
}

type DatabaseRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *DatabaseRepository {
	return &DatabaseRepository{db: db}
}

func (repo *DatabaseRepository) Create(ctx context.Context, c *Customer) error {
	const query = `
		INSERT INTO customers
			(id, document, name, score, risk_level, income_range, status, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`
	_, err := repo.db.ExecContext(ctx, query,
		c.ID, c.Document, c.Name, c.Score, c.RiskLevel,
		c.IncomeRange, c.Status, c.CreatedAt, c.UpdatedAt,
	)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return error_handler.ErrDuplicateDocument
	}
	return err
}

func (repo *DatabaseRepository) List(ctx context.Context, limit, offset int) ([]Customer, error) {
	const query = `
		SELECT id, document, name, score, risk_level, income_range, status, created_at, updated_at
		FROM customers
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`
	rows, err := repo.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Customer
	for rows.Next() {
		var c Customer
		if err := rows.Scan(&c.ID, &c.Document, &c.Name, &c.Score,
			&c.RiskLevel, &c.IncomeRange, &c.Status, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func (repo *DatabaseRepository) GetByID(ctx context.Context, id string) (*Customer, error) {
	return repo.getOne(ctx, "id = $1", id)
}

func (repo *DatabaseRepository) GetByDocument(ctx context.Context, doc string) (*Customer, error) {
	return repo.getOne(ctx, "document = $1", doc)
}

func (repo *DatabaseRepository) getOne(ctx context.Context, where, arg string) (*Customer, error) {
	q := `SELECT id, document, name, score, risk_level, income_range, status, created_at, updated_at
	      FROM customers WHERE ` + where
	var c Customer
	err := repo.db.QueryRowContext(ctx, q, arg).Scan(
		&c.ID, &c.Document, &c.Name, &c.Score,
		&c.RiskLevel, &c.IncomeRange, &c.Status, &c.CreatedAt, &c.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, error_handler.ErrCustomerNotFound
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (repo *DatabaseRepository) UpdateStatus(ctx context.Context, id, status string, updatedAt time.Time) error {
	const query = `UPDATE customers SET status = $1, updated_at = $2 WHERE id = $3`
	res, err := repo.db.ExecContext(ctx, query, status, updatedAt, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return error_handler.ErrCustomerNotFound
	}
	return nil
}
