package repository

import (
	"context"
	"database/sql"
	"fmt"
	"snipdrop-rest-api/internal/app/snipdrop-api/model"
	"strings"
)

type SnippetRepository struct {
	DB *sql.DB
}

type SnippetQueryParams struct {
	Page     int    // Page number
	PageSize int    // Number of snippets per page
	Sort     string // Sort field and direction, e.g., "created_at desc"
	Search   string // Search term
}

// ListSnippets retrieves snippets based on pagination, sorting, and searching parameters.
func (r *SnippetRepository) ListSnippets(ctx context.Context, params SnippetQueryParams) ([]model.Snippet, error) {
	// Set default values for pagination if not provided
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 10
	}

	offset := (params.Page - 1) * params.PageSize

	// Start building the query
	query := "SELECT id, created_at, value, name, author FROM snippets"
	whereClauses := []string{}
	args := []interface{}{}

	// Search filter
	if params.Search != "" {
		searchQuery := "%" + params.Search + "%"
		whereClauses = append(whereClauses, fmt.Sprintf("(name ILIKE $%d OR value ILIKE $%d)", len(args)+1, len(args)+1))
		args = append(args, searchQuery)
	}

	// Add WHERE clauses if any
	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Sorting
	if params.Sort != "" {
		query += fmt.Sprintf(" ORDER BY %s", params.Sort)
	} else {
		// Default sorting
		query += " ORDER BY created_at DESC"
	}

	// Pagination
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
	args = append(args, params.PageSize, offset)

	// Execute the query
	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	snippets := []model.Snippet{}
	for rows.Next() {
		var snippet model.Snippet
		if err := rows.Scan(&snippet.ID, &snippet.CreatedAt, &snippet.Value, &snippet.Name, &snippet.Author); err != nil {
			return nil, err
		}
		snippets = append(snippets, snippet)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}

// NewSnippet creates a new snippet in the database.
func (r *SnippetRepository) NewSnippet(ctx context.Context, s model.Snippet) error {
	_, err := r.DB.ExecContext(ctx, "INSERT INTO snippets (id, created_at, value, name, author) VALUES ($1, $2, $3, $4, $5)",
		s.ID, s.CreatedAt, s.Value, s.Name, s.Author)
	return err
}

// GetSnippet retrieves a snippet by its ID.
func (r *SnippetRepository) GetSnippet(ctx context.Context, id string) (*model.Snippet, error) {
	snippet := model.Snippet{}
	err := r.DB.QueryRowContext(ctx, "SELECT id, created_at, value, name, author FROM snippets WHERE id = $1", id).Scan(&snippet.ID, &snippet.CreatedAt, &snippet.Value, &snippet.Name, &snippet.Author)
	if err != nil {
		return nil, err
	}
	return &snippet, nil
}

// DeleteSnippet removes a snippet from the database.
func (r *SnippetRepository) DeleteSnippet(ctx context.Context, id string) error {
	_, err := r.DB.ExecContext(ctx, "DELETE FROM snippets WHERE id = $1", id)
	return err
}
