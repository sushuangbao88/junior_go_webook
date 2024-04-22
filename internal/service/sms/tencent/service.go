package tencent

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/ecodeclub/ekit"
	"github.com/ecodeclub/ekit/slice"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tencentSms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

type Service struct {
	client   *tencentSms.Client
	appId    *string
	signName *string
}

func NewClient() (*tencentSms.Client, error) {
	//从配置中读取密钥对
	secretId := os.Getenv("TENCENTCLOUD+SECRET_ID")
	secretKey := os.Getenv("TENCENTCLOUD+SECRET_KEY")

	//根据密钥对，实例认证对象
	credential := common.NewCredential(secretId, secretKey)
	cpf := profile.NewClientProfile()

	return tencentSms.NewClient(credential, "ap-nanjing", cpf)
}

func NewService(client *tencentSms.Client, appId string, signName string) *Service {
	return &Service{
		client:   client,
		appId:    &appId,
		signName: &signName,
	}
}

func (s *Service) Send(ctx context.Context, tmpId string, args []string, numbers ...string) error {
	request := tencentSms.NewSendSmsRequest()
	request.SetContext(ctx)
	request.SmsSdkAppId = s.appId
	request.SignName = s.signName
	request.TemplateId = ekit.ToPtr[string](tmpId)
	request.TemplateParamSet = s.toPtrSlice(args)
	request.PhoneNumberSet = s.toPtrSlice(numbers)
	respose, err := s.client.SendSms(request)
	if err != nil {
		return err
	}
	b, _ := json.Marshal(respose.Response)
	fmt.Printf("%s", b) //尝试打印返回结果

	return nil
}

func (s *Service) toPtrSlice(data []string) []*string {
	return slice.Map[string, *string](data, func(idx int, src string) *string {
		return &src
	})
}
