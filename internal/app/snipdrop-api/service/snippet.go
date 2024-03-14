package service

import (
	"context"
	"snipdrop-rest-api/internal/app/snipdrop-api/model"
	"snipdrop-rest-api/internal/app/snipdrop-api/repository"

	"go.uber.org/zap"
)

type SnippetService struct {
	Repo   *repository.SnippetRepository
	Logger *zap.Logger
}

func (s *SnippetService) ListSnippets(ctx context.Context, params repository.SnippetQueryParams) ([]model.Snippet, error) {
	return s.Repo.ListSnippets(ctx, params)
}

func (s *SnippetService) CreateSnippet(ctx context.Context, snippet model.Snippet) error {
	return s.Repo.NewSnippet(ctx, snippet)
}

func (s *SnippetService) GetSnippet(ctx context.Context, id string) (*model.Snippet, error) {
	return s.Repo.GetSnippet(ctx, id)
}

func (s *SnippetService) DeleteSnippet(ctx context.Context, id string) error {
	return s.Repo.DeleteSnippet(ctx, id)
}
