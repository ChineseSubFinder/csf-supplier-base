package whisper_helper

import "github.com/go-resty/resty/v2"

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
func (w WhisperClient) SendTask(taskID int, audioFPath string, language string) error {

	return nil
}

func (w WhisperClient) GetTaskStatus(taskID int) {

}
