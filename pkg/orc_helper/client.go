package orc_helper

import (
	"encoding/json"
	"github.com/ChineseSubFinder/csf-supplier-base/pkg/settings"
	"github.com/allanpk716/rod_helper"
	"github.com/go-resty/resty/v2"
)

type ORCHelper struct {
	orcBaseUrl string
	client     *resty.Client
}

func NewORCHelper(orcBaseUrl string) *ORCHelper {

	opt := rod_helper.NewHttpClientOptions(settings.Get().TimeConfig.GetOnePageTimeOut())
	client, err := rod_helper.NewHttpClient(opt)
	if err != nil {
		panic(err)
	}
	return &ORCHelper{
		orcBaseUrl: orcBaseUrl,
		client:     client}
}

func (o ORCHelper) GetStatus() bool {

	var status ReplyStatusInfo
	_, err := o.client.R().SetResult(&status).Get(o.orcBaseUrl)
	if err != nil {
		return false
	}

	return true
}

func (o ORCHelper) File(imgFPath string) (string, error) {

	var imgFileInfo ReplyImgFileInfo
	resp, err := o.client.R().
		SetFile("file", imgFPath).
		Post(o.orcBaseUrl + "/file")
	if err != nil {
		return "", err
	}
	// 从字符串转Struct
	err = json.Unmarshal(resp.Body(), &imgFileInfo)
	if err != nil {
		return "", err
	}

	return imgFileInfo.Result, nil
}

type ReplyStatusInfo struct {
	Message   string `json:"message"`
	Tesseract struct {
		Languages []string `json:"languages"`
		Version   string   `json:"version"`
	} `json:"tesseract"`
	Version string `json:"version"`
}

type ReplyImgFileInfo struct {
	Result  string `json:"result"`
	Version string `json:"version"`
}
