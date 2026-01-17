package postgres

import (
	"context"

	"github.com/bladimirbalbin/portafolio-api/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProjectRepo struct {
	db *pgxpool.Pool
}

func NewProjectRepo(db *pgxpool.Pool) *ProjectRepo {
	return &ProjectRepo{db: db}
}

func (r *ProjectRepo) List(ctx context.Context, featured *bool, tag *string) ([]domain.Project, error) {
	// Filtros opcionales
	query := `
SELECT id, slug, name, description, repo_url, demo_url, tags, featured, created_at, updated_at
FROM projects
WHERE ($1::boolean IS NULL OR featured = $1)
  AND ($2::text IS NULL OR $2 = ANY(tags))
ORDER BY featured DESC, updated_at DESC;
`
	rows, err := r.db.Query(ctx, query, featured, tag)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	projects := make([]domain.Project, 0)
	for rows.Next() {
		var p domain.Project
		if err := rows.Scan(
			&p.ID, &p.Slug, &p.Name, &p.Description, &p.RepoURL, &p.DemoURL,
			&p.Tags, &p.Featured, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return projects, nil
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
		return nil, err
	}
	return &p, nil
}
