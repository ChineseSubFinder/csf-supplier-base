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
