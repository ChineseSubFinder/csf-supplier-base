package captcha_client

import (
	"encoding/base64"
	"encoding/json"
	"github.com/ChineseSubFinder/csf-supplier-base/pkg/settings"
	"github.com/allanpk716/rod_helper"
	"github.com/go-resty/resty/v2"
	"io/ioutil"
)

type CaptchaClient struct {
	baseUrl string
	client  *resty.Client
}

func NewCaptchaClient(baseUrl string) *CaptchaClient {

	opt := rod_helper.NewHttpClientOptions(settings.Get().TimeConfig.GetOnePageTimeOut())
	client, err := rod_helper.NewHttpClient(opt)
	if err != nil {
		panic(err)
	}
	return &CaptchaClient{
		baseUrl: baseUrl,
		client:  client}
}

func (c CaptchaClient) GetFromImgFile(imgFPath string) (string, error) {

	srcByte, err := ioutil.ReadFile(imgFPath)
	if err != nil {
		return "", err
	}

	res := base64.StdEncoding.EncodeToString(srcByte)

	return c.GetFromImgBase64(res)
}

func (c CaptchaClient) GetFromImgBase64(imgBase64 string) (string, error) {

	var requestCaptchaInfo RequestCaptchaInfo
	requestCaptchaInfo.Image = imgBase64
	var replyCaptchaInfo ReplyCaptchaInfo
	resp, err := c.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(requestCaptchaInfo).
		Post(c.baseUrl + "/captcha/v1")
	if err != nil {
		return "", err
	}

	// 从字符串转Struct
	err = json.Unmarshal(resp.Body(), &replyCaptchaInfo)
	if err != nil {
		return "", err
	}

	if replyCaptchaInfo.Success == false {
		return "", err
	}

	return replyCaptchaInfo.Message, nil
}

type RequestCaptchaInfo struct {
	Image string `json:"image"`
}

type ReplyCaptchaInfo struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Success bool   `json:"success"`
	Uid     string `json:"uid"`
}
