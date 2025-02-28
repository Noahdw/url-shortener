package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

type Querier interface {
	CreateUrlMapping(ctx context.Context, arg CreateUrlMappingParams) (pgconn.CommandTag, error)
	GetOriginalUrlFromShortCode(ctx context.Context, shortCode string) (string, error)
	GetShortCodeFromOriginalUrl(ctx context.Context, originalUrl string) (string, error)
}

type RepoMock struct {
	db map[string]CreateUrlMappingParams
}

func NewRepoMock() *RepoMock {
	return &RepoMock{
		db: make(map[string]CreateUrlMappingParams),
	}
}

func (m *RepoMock) CreateUrlMapping(ctx context.Context, arg CreateUrlMappingParams) (pgconn.CommandTag, error) {
	m.db[arg.ShortCode] = arg
	return pgconn.CommandTag{}, nil
}

func (m *RepoMock) GetOriginalUrlFromShortCode(ctx context.Context, shortCode string) (string, error) {
	arg, has := m.db[shortCode]
	if !has {
		return "", errors.New("")
	}
	return arg.OriginalUrl, nil
}

func (m *RepoMock) GetShortCodeFromOriginalUrl(ctx context.Context, originalUrl string) (string, error) {
	arg, has := m.db[originalUrl]
	if !has {
		return "", errors.New("")
	}
	return arg.ShortCode, nil
}
