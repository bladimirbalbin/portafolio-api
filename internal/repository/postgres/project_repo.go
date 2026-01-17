package postgres

import (
	"context"
	"errors"

	"github.com/bladimirbalbin/portafolio-api/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProjectRepo struct {
	db *pgxpool.Pool
}

func NewProjectRepo(db *pgxpool.Pool) *ProjectRepo {
	return &ProjectRepo{db: db}
}

type ProjectListOpts struct {
	Featured *bool
	Tag      *string
	Limit    int
	Offset   int
	Sort     string // "recent" | "featured"
}

func (r *ProjectRepo) List(ctx context.Context, opts ProjectListOpts) ([]domain.Project, error) {
	limit := opts.Limit
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	offset := opts.Offset
	if offset < 0 {
		offset = 0
	}

	// Whitelist para evitar SQL injection en ORDER BY
	orderBy := "featured DESC, updated_at DESC"
	switch opts.Sort {
	case "", "featured":
		orderBy = "featured DESC, updated_at DESC"
	case "recent":
		orderBy = "updated_at DESC"
	default:
		return nil, errors.New("invalid sort (use 'recent' or 'featured')")
	}

	query := `
SELECT id, slug, name, description, repo_url, demo_url, tags, featured, created_at, updated_at
FROM projects
WHERE ($1::boolean IS NULL OR featured = $1)
  AND ($2::text IS NULL OR $2 = ANY(tags))
ORDER BY ` + orderBy + `
LIMIT $3 OFFSET $4;
`
	rows, err := r.db.Query(ctx, query, opts.Featured, opts.Tag, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]domain.Project, 0)
	for rows.Next() {
		var p domain.Project
		if err := rows.Scan(
			&p.ID, &p.Slug, &p.Name, &p.Description, &p.RepoURL, &p.DemoURL,
			&p.Tags, &p.Featured, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *ProjectRepo) GetBySlug(ctx context.Context, slug string) (*domain.Project, error) {
	query := `
SELECT id, slug, name, description, repo_url, demo_url, tags, featured, created_at, updated_at
FROM projects
WHERE slug = $1
LIMIT 1;
`
	var p domain.Project
	err := r.db.QueryRow(ctx, query, slug).Scan(
		&p.ID, &p.Slug, &p.Name, &p.Description, &p.RepoURL, &p.DemoURL,
		&p.Tags, &p.Featured, &p.CreatedAt, &p.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}
