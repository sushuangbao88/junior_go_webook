package service

import (
	"context"
	"fmt"
	"math/rand"

	"example.com/basic-gin/webook/internal/repository"
	"example.com/basic-gin/webook/internal/service/sms"
)

type CodeService struct {
	repo *repository.CodeRepository
	sms  sms.Service
}

func NewCodeService(repo *repository.CodeRepository, smsSvc sms.Service) *CodeService {
	return &CodeService{
		repo: repo,
		sms:  smsSvc,
	}
}

func (s *CodeService) Set(ctx context.Context, biz string, phone string) error {
	//随机生成code
	code := s.generateCode()
	tmpId := "1877556"

	//设置验证码
	err := s.repo.Set(ctx, biz, phone, code)
	if err != nil {
		return err
	}

	//发送验证码
	return s.sms.Send(ctx, tmpId, []string{code}, phone)
}

func (s *CodeService) Verify(ctx context.Context, biz string, phone string, code string) (bool, error) {
	ok, err := s.repo.Verify(ctx, biz, phone, code)
	if err == repository.ErrCodeVerifyTooMany {
		return false, nil //在这里截住“验证频繁”的错误，不暴露出去，化成“校验失败”的错误
	}
	return ok, err
}

// 返回6位数字的字符串
func (s *CodeService) generateCode() string {
	code := rand.Intn(1000000)
	return fmt.Sprintf("%06d", code)
}
