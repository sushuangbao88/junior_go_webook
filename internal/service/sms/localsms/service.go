package localsms

import (
	"context"
	"fmt"
)

/*
本地模拟发送sms消息
*/
type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Send(ctx context.Context, tmpId string, args []string, numbers ...string) error {
	fmt.Printf("给手机号[%s]发送验证码：%s,成功！", numbers[0], args[0])

	return nil
}
