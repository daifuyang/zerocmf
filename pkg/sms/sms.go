package sms

type Sms interface {
	SendSms(phoneNumber string, code string) error
}

const (
	AliSms     = "ali"
	NetEaseSms = "netease"
	TencentSms = "tencent"
)

func NewSms(accessKeyId string, accessKeySecret string, signName string, templateCode string, provider string) Sms {
	switch provider {
	case "ali":
		return &Dysms{
			AccessKeyId:     accessKeyId,
			AccessKeySecret: accessKeySecret,
			SignName:        signName,
			TemplateCode:    templateCode,
		}
	default:
		return nil
	}
}
