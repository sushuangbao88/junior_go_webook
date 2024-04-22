package sms

import (
	"log"

	"example.com/basic-gin/webook/internal/service/sms/localsms"
	"example.com/basic-gin/webook/internal/service/sms/tencent"
)

func NewService(name string) Service {
	switch name {
	case "tencent":
		client, err := tencent.NewClient()
		if err != nil {
			log.Println("sms服务初始化失败!")
		}
		s := tencent.NewService(client, "app_id", "sign_name")
		return s
	case "local":
		return localsms.NewService()
	default:
		log.Println("没有对应的sms实现")
		return nil
	}
}
