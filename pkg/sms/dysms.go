// 阿里短信
package sms

import (
	"fmt"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

type Dysms struct {
	// 必填，您的 AccessKey ID
	AccessKeyId string
	// 必填，您的 AccessKey Secret
	AccessKeySecret string
	SignName        string
	TemplateCode    string
}

// func (d *Dysms) NewDysms() *Dysms {
// 	accessKeyId := d.AccessKeyId
// 	accessKeySecret := d.AccessKeySecret
// 	signName := d.SignName
// 	templateCode := d.TemplateCode
// 	return &Dysms{
// 		AccessKeyId:     accessKeyId,
// 		AccessKeySecret: accessKeySecret,
// 		SignName:        signName,
// 		TemplateCode:    templateCode,
// 	}
// }

func (d *Dysms) SendSms(phoneNumber string, templateParam string) error {

	fmt.Println("Sending SMS", phoneNumber, templateParam)

	sendSmsRequest := &dysmsapi.SendSmsRequest{
		PhoneNumbers:  tea.String(phoneNumber),
		SignName:      tea.String(d.SignName),
		TemplateCode:  tea.String(d.TemplateCode),
		TemplateParam: tea.String(templateParam),
	}
	client, err := dysmsapi.NewClient(&openapi.Config{
		// 必填，您的 AccessKey ID
		AccessKeyId: tea.String(d.AccessKeyId),
		// 必填，您的 AccessKey Secret
		AccessKeySecret: tea.String(d.AccessKeySecret),
		Endpoint:        tea.String("dysmsapi.aliyuncs.com"),
	})

	if err != nil {
		return err
	}

	runtime := &util.RuntimeOptions{}

	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()

		// 复制代码运行请自行打印 API 的返回值
		_, err = client.SendSmsWithOptions(sendSmsRequest, runtime)
		if err != nil {
			return err
		}
		return nil
	}()

	if tryErr != nil {
		var error = &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			error = _t
		} else {
			error.Message = tea.String(tryErr.Error())
		}
		// 如有需要，请打印 error
		_, err = util.AssertAsString(error.Message)
		if err != nil {
			return err
		}
		return error
	}
	return nil
}
