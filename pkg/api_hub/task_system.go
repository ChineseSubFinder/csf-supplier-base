package api_hub

import "github.com/ChineseSubFinder/csf-supplier-base/db/task_system"

// GetOneTaskReq 获取一个任务的请求
type GetOneTaskReq struct {
	TaskType task_system.TaskType `json:"task_type"` // 请求什么类型的任务
	ApiKey   string               `json:"api_key"`   // 身份密钥
}

// GetOneTaskResp 获取一个任务的响应
type GetOneTaskResp struct {
	Status          int                  `json:"status"`                 // 任务的状态 0 失败，1 成功
	Message         string               `json:"message"`                // 任务的状态信息
	TaskType        task_system.TaskType `json:"task_type"`              // 任务的类型
	DataDownloadUrl string               `json:"task_data_download_url"` // 任务数据的下载地址
	DataVersion     string               `json:"task_data_version"`      // 任务数据的版本
}

// ----------------------------------------------

// AddMachineTranslationTaskPackageReq 添加一个任务的请求
type AddMachineTranslationTaskPackageReq struct {
	ImdbId  string `json:"imdb_id"`
	IsMovie bool   `json:"is_movie"` // 是电影还是连续剧
	Season  int    `json:"season"`   // 电影则是 -1
	Episode int    `json:"episode"`  // 连续剧则是 -1

	IsAudioOrSRT bool   `json:"is_audio_or_srt"` // 是音频还是字幕
	SubSha256    string `json:"sub_sha256"`      // 文件的 SHA256
	FileName     string `json:"file_name"`       // 文件的名称
	FileSize     int    `json:"file_size"`       // 文件大小，单位：字节

	AudioSrcLanguage   string `json:"audio_src_language"`  // 音频的源语言
	TranslatedLanguage string `json:"translated_language"` // 期望的翻译后的语言

	ApiKey string `json:"api_key"` // 身份密钥
}

// AddMachineTranslationTaskPackageResp 添加一个任务的响应
type AddMachineTranslationTaskPackageResp struct {
	Status        int    `json:"status"`          // 任务的状态 0 失败，1 成功
	Message       string `json:"message"`         // 任务的状态信息
	TaskPackageId string `json:"task_package_id"` // 任务包的ID
	UploadURL     string `json:"upload_url"`      // 上传文件的URL
}
