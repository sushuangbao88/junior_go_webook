package sms

import "context"

type Service interface {
	Send(ctx context.Context, tmpId string, args []string, numbers ...string) error
}
