package task_system

import (
	"encoding/json"
	"github.com/ChineseSubFinder/csf-supplier-base/pkg"
	"github.com/ChineseSubFinder/csf-supplier-base/pkg/api_hub"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

type TaskSystemClient struct {
	client          *resty.Client
	apiKey          string // subtitle best api key
	comunicateToken string // 所有包通信时候写在 Authorization 的 Token
}

func NewTaskSystemClient(apiKey string) *TaskSystemClient {
	return &TaskSystemClient{
		client:          resty.New(),
		apiKey:          apiKey,
		comunicateToken: "5akwmGAbuFWqgaZf9QwT",
	}
}

/*
Audio 支持的语言列表
		DEFAULT_ALIGN_MODELS_TORCH = {
		"en": "WAV2VEC2_ASR_BASE_960H",
		"fr": "VOXPOPULI_ASR_BASE_10K_FR",
		"de": "VOXPOPULI_ASR_BASE_10K_DE",
		"es": "VOXPOPULI_ASR_BASE_10K_ES",
		"it": "VOXPOPULI_ASR_BASE_10K_IT",
	}

	DEFAULT_ALIGN_MODELS_HF = {
		"ja": "jonatasgrosman/wav2vec2-large-xlsr-53-japanese",
		"zh": "jonatasgrosman/wav2vec2-large-xlsr-53-chinese-zh-cn",
		"nl": "jonatasgrosman/wav2vec2-large-xlsr-53-dutch",
		"uk": "Yehor/wav2vec2-xls-r-300m-uk-with-small-lm",
		"pt": "jonatasgrosman/wav2vec2-large-xlsr-53-portuguese",
		"ar": "jonatasgrosman/wav2vec2-large-xlsr-53-arabic",
		"ru": "jonatasgrosman/wav2vec2-large-xlsr-53-russian",
		"pl": "jonatasgrosman/wav2vec2-large-xlsr-53-polish",
		"hu": "jonatasgrosman/wav2vec2-large-xlsr-53-hungarian",
		"fi": "jonatasgrosman/wav2vec2-large-xlsr-53-finnish",
		"fa": "jonatasgrosman/wav2vec2-large-xlsr-53-persian",
		"el": "jonatasgrosman/wav2vec2-large-xlsr-53-greek",
		"tr": "mpoyraz/wav2vec2-xls-r-300m-cv7-turkish",
		"he": "imvladikon/wav2vec2-xls-r-300m-hebrew",
}
*/

/*
	SRT 翻译支持的语言列表
	CN， 默认这个先，就是翻译到中文
*/

// AddMachineTranslationTask 添加机翻任务
func (c *TaskSystemClient) AddMachineTranslationTask(ImdbId string, IsMovie bool, Season, Episode int,
	IsAudioOrSRT bool, filFullPath string, AudioSrcLanguage, TranslatedLanguage string) (*api_hub.AddMachineTranslationTaskPackageResp, error) {

	if pkg.IsFile(filFullPath) == false {
		return nil, errors.New("file not exist")
	}
	// 获取文件的名称，文件的 SHA256 以及 文件的大小
	fileName, fileSha256, fileSize, err := pkg.GetFileSha256AndSize(filFullPath)
	if err != nil {
		return nil, err
	}

	addTaskReq := api_hub.AddMachineTranslationTaskPackageReq{
		ImdbId:             ImdbId,
		IsMovie:            IsMovie,
		Season:             Season,
		Episode:            Episode,
		IsAudioOrSRT:       IsAudioOrSRT,
		FileSha256:         fileSha256,
		FileName:           fileName,
		FileSize:           fileSize,
		AudioSrcLanguage:   AudioSrcLanguage,
		TranslatedLanguage: TranslatedLanguage,
		ApiKey:             c.apiKey,
	}
	// 发送请求
	resp, err := c.client.R().
		SetBody(addTaskReq).
		SetAuthToken(c.comunicateToken).
		Post(serverUrlBase + "/add-machine-translation-task")

	if err != nil {
		return nil, err
	}
	var reply api_hub.AddMachineTranslationTaskPackageResp
	// 从字符串转Struct
	err = json.Unmarshal(resp.Body(), &reply)
	if err != nil {
		return nil, err
	}

	return &reply, nil
}

// SetFirstPackageTaskDone 设置第一个包的任务完成，简单来说就是上传的 音频或者字幕 任务完成，需要反馈给服务器
func (c *TaskSystemClient) SetFirstPackageTaskDone(TaskPackageId, token string) (*ReplyCommon, error) {

	addTaskReq := api_hub.SetFirstPackageTaskDoneReq{
		TaskPackageId: TaskPackageId,
		Token:         token,
		ApiKey:        c.apiKey,
	}
	// 发送请求
	resp, err := c.client.R().
		SetBody(addTaskReq).
		SetAuthToken(c.comunicateToken).
		Post(serverUrlBase + "/set-first-package-task-done")

	if err != nil {
		return nil, err
	}
	var reply ReplyCommon
	// 从字符串转Struct
	err = json.Unmarshal(resp.Body(), &reply)
	if err != nil {
		return nil, err
	}

	return &reply, nil
}

// GetTaskPackageStatus 获取任务包的状态
func (c *TaskSystemClient) GetTaskPackageStatus(taskPackageID string) (*api_hub.GetTaskPackageInfoResp, error) {

	getTaskPackageInfoReq := api_hub.GetTaskPackageInfoReq{
		TaskPackageId: taskPackageID,
		ApiKey:        c.apiKey,
	}
	// 发送请求
	resp, err := c.client.R().
		SetBody(getTaskPackageInfoReq).
		SetAuthToken(c.comunicateToken).
		Post(serverUrlBase + "/get-task-package-status")

	if err != nil {
		return nil, err
	}
	var reply api_hub.GetTaskPackageInfoResp
	// 从字符串转Struct
	err = json.Unmarshal(resp.Body(), &reply)
	if err != nil {
		return nil, err
	}

	return &reply, nil
}

func (c *TaskSystemClient) GetTranslatedResult(taskPackageID string) (*api_hub.GetTranslatedResultResp, error) {

	getTaskPackageInfoReq := api_hub.GetTaskPackageInfoReq{
		TaskPackageId: taskPackageID,
		ApiKey:        c.apiKey,
	}
	// 发送请求
	resp, err := c.client.R().
		SetBody(getTaskPackageInfoReq).
		SetAuthToken(c.comunicateToken).
		Post(serverUrlBase + "/get-translated-result")

	if err != nil {
		return nil, err
	}
	var reply api_hub.GetTranslatedResultResp
	// 从字符串转Struct
	err = json.Unmarshal(resp.Body(), &reply)
	if err != nil {
		return nil, err
	}

	return &reply, nil
}

// CancelTaskPackage 取消任务包，提交和返回的数据都跟 GetTaskPackageStatus 一样
func (c *TaskSystemClient) CancelTaskPackage(taskPackageID string) (*api_hub.GetTaskPackageInfoResp, error) {

	getTaskPackageInfoReq := api_hub.GetTaskPackageInfoReq{
		TaskPackageId: taskPackageID,
		ApiKey:        c.apiKey,
	}
	// 发送请求
	resp, err := c.client.R().
		SetBody(getTaskPackageInfoReq).
		SetAuthToken(c.comunicateToken).
		Post(serverUrlBase + "/cancel-task-package")

	if err != nil {
		return nil, err
	}
	var reply api_hub.GetTaskPackageInfoResp
	// 从字符串转Struct
	err = json.Unmarshal(resp.Body(), &reply)
	if err != nil {
		return nil, err
	}

	return &reply, nil
}

const (
	//serverUrlBase = "https://api.subtitle.best/task-system/v1"
	serverUrlBase = "http://127.0.0.1:8888/task-system/v1"
)

type ReplyCommon struct {
	Status  int    `json:"status"` // 任务的状态 0 失败，1 成功
	Message string `json:"message,omitempty"`
}
