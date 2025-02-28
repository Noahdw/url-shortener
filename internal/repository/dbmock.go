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
	sToO map[string]CreateUrlMappingParams
	oToS map[string]CreateUrlMappingParams
}

func NewRepoMock() *RepoMock {
	return &RepoMock{
		sToO: make(map[string]CreateUrlMappingParams),
		oToS: make(map[string]CreateUrlMappingParams),
	}
}

func (m *RepoMock) CreateUrlMapping(ctx context.Context, arg CreateUrlMappingParams) (pgconn.CommandTag, error) {
	_, has := m.sToO[arg.ShortCode]
	if has {
		return pgconn.CommandTag{}, &pgconn.PgError{Code: "23505"} // 23505 = unique key violation
	}
	m.sToO[arg.ShortCode] = arg
	m.oToS[arg.OriginalUrl] = arg
	return pgconn.CommandTag{}, nil
}

func (m *RepoMock) GetOriginalUrlFromShortCode(ctx context.Context, shortCode string) (string, error) {
	arg, has := m.sToO[shortCode]
	if !has {
		return "", errors.New("")
	}
	return arg.OriginalUrl, nil
}

func (m *RepoMock) GetShortCodeFromOriginalUrl(ctx context.Context, originalUrl string) (string, error) {
	arg, has := m.oToS[originalUrl]
	if !has {
		return "", errors.New("")
	}
	return arg.ShortCode, nil
}
