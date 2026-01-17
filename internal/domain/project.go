package domain

import "time"

type Project struct {
	ID          int64     `json:"id"`
	Slug        string    `json:"slug"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	RepoURL     *string   `json:"repo_url,omitempty"`
	DemoURL     *string   `json:"demo_url,omitempty"`
	Tags        []string  `json:"tags"`
	Featured    bool      `json:"featured"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
