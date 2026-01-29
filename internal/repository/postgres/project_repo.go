package postgres

import (
	"context"
	"errors"

	"github.com/bladimirbalbin/portafolio-api/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("project not found")

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
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (r *ProjectRepo) Create(ctx context.Context, p domain.Project) (*domain.Project, error) {
	query := `
INSERT INTO projects (slug, name, description, repo_url, demo_url, tags, featured)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, slug, name, description, repo_url, demo_url, tags, featured, created_at, updated_at;
`
	var out domain.Project
	err := r.db.QueryRow(ctx, query,
		p.Slug, p.Name, p.Description, p.RepoURL, p.DemoURL, p.Tags, p.Featured,
	).Scan(
		&out.ID, &out.Slug, &out.Name, &out.Description, &out.RepoURL, &out.DemoURL,
		&out.Tags, &out.Featured, &out.CreatedAt, &out.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *ProjectRepo) UpdateBySlug(ctx context.Context, slug string, p domain.Project) (*domain.Project, error) {
	query := `
UPDATE projects
SET name = $2,
    description = $3,
    repo_url = $4,
    demo_url = $5,
    tags = $6,
    featured = $7
WHERE slug = $1
RETURNING id, slug, name, description, repo_url, demo_url, tags, featured, created_at, updated_at;
`
	var out domain.Project
	err := r.db.QueryRow(ctx, query,
		slug, p.Name, p.Description, p.RepoURL, p.DemoURL, p.Tags, p.Featured,
	).Scan(
		&out.ID, &out.Slug, &out.Name, &out.Description, &out.RepoURL, &out.DemoURL,
		&out.Tags, &out.Featured, &out.CreatedAt, &out.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &out, nil
}

func (r *ProjectRepo) DeleteBySlug(ctx context.Context, slug string) error {
	cmd, err := r.db.Exec(ctx, `DELETE FROM projects WHERE slug = $1;`, slug)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
