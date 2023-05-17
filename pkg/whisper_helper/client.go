package whisper_helper

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
)

type WhisperClient struct {
	client    *resty.Client
	serverUrl string
	token     string
}

func NewWhisperClient(serverUrl string, token string) *WhisperClient {
	return &WhisperClient{
		serverUrl: serverUrl,
		client:    resty.New(),
		token:     token}
}

/*
支持的语言列表
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

// SendTask 发送任务
func (w WhisperClient) SendTask(taskID int, audioFPath string, language string) (*SendTaskReply, error) {

	task := SendTask{
		TaskID:     taskID,
		InputAudio: audioFPath,
		Language:   language,
	}
	// 发送请求
	resp, err := w.client.R().
		SetBody(task).
		SetAuthToken(w.token).
		Post(w.serverUrl + "/transcribe")

	if err != nil {
		return nil, err
	}
	var reply SendTaskReply
	// 从字符串转Struct
	err = json.Unmarshal(resp.Body(), &reply)
	if err != nil {
		return nil, err
	}

	return &reply, nil
}

// GetTaskStatus 获取任务状态
func (w WhisperClient) GetTaskStatus(taskID int) (*SendTaskReply, error) {

	resp, err := w.client.R().
		SetAuthToken(w.token).
		Get(w.serverUrl + "/transcribe" + "/" + fmt.Sprintf("%d", taskID))
	if err != nil {
		return nil, err
	}

	var reply SendTaskReply
	// 从字符串转Struct
	err = json.Unmarshal(resp.Body(), &reply)
	if err != nil {
		return nil, err
	}

	return &reply, nil
}
