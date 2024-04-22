package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var (
	//go:embed lua/set_code.lua
	luaRedisSetCode string
	//go:embed lua/verify_code.lua
	luaRedisVerifyCode string
)

type CodeRedisCache struct {
	cmd redis.Cmdable
}

func NewCodeRedisCache(cmd redis.Cmdable) CodeCache {
	return &CodeRedisCache{
		cmd: cmd,
	}
}
func (crc *CodeRedisCache) Set(ctx context.Context, biz string, phone string, code string) error {
	res, err := crc.cmd.Eval(ctx, luaRedisSetCode, []string{crc.key(biz, phone)}, code).Int()
	if err != nil {
		return err
	}

	switch res {
	case -2:
		return errors.New("验证码存在，但是没有过期时间")
	case -1:
		return ErrCodeSendTooMany
	default:
		return nil
	}
}

func (crc *CodeRedisCache) Verify(ctx context.Context, biz string, phone string, code string) (bool, error) {
	res, err := crc.cmd.Eval(ctx, luaRedisVerifyCode, []string{crc.key(biz, phone)}, code).Int()
	if err != nil {
		return false, err
	}
	switch res {
	case -2:
		return false, nil
	case -1:
		return false, ErrCodeVerifyTooMany
	default:
		return true, nil
	}
}

func (crc *CodeRedisCache) key(biz, phone string) string {
	return fmt.Sprintf("%s:%s", biz, phone)
}
