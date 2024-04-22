package repository

import (
	"context"

	"example.com/basic-gin/webook/internal/repository/cache"
)

var (
	ErrCodeSendTooMany   = cache.ErrCodeSendTooMany
	ErrCodeVerifyTooMany = cache.ErrCodeVerifyTooMany
)

type CodeRepository struct {
	cache cache.CodeCache
}

func NewCodeRepository(cc cache.CodeCache) *CodeRepository {
	return &CodeRepository{
		cache: cc,
	}
}
func (c *CodeRepository) Set(ctx context.Context, biz string, phone string, code string) error {
	return c.cache.Set(ctx, biz, phone, code)
}

func (c *CodeRepository) Verify(ctx context.Context, biz string, phone string, code string) (bool, error) {
	return c.cache.Verify(ctx, biz, phone, code)
}
