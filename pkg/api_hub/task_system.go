package api_hub

import "github.com/VideoTranslationTools/base/db/task_system"

// GetOneTaskReq 获取一个任务的请求
type GetOneTaskReq struct {
	TaskType     task_system.TaskType `json:"task_type"`      // 请求什么类型的任务
	ApiKey       string               `json:"api_key"`        // 身份密钥
	WorkerApiKey string               `json:"worker_api_key"` // worker 身份密钥
}

// GetOneTaskResp 获取一个任务的响应
type GetOneTaskResp struct {
	Status              int                  `json:"status"`                 // 任务的状态 0 失败，1 成功
	Message             string               `json:"message"`                // 任务的状态信息
	TaskType            task_system.TaskType `json:"task_type"`              // 任务的类型
	SrcDataDownloadUrl  string               `json:"src_data_download_url"`  // 任务数据的下载地址
	FinishDataUploadUrl string               `json:"finish_data_upload_url"` // 任务数据的上传地址
	AudioSrcLanguage    string               `json:"audio_src_language"`     // 音频的源语言
	TranslatedLanguage  string               `json:"translated_language"`    // 期望的翻译后的语言
	Token               string               `json:"token"`                  // 任务的 token
	DataVersion         string               `json:"task_data_version"`      // 任务数据的版本
}

type GetOneTaskFinishDataUploadUrlReq struct {
	TaskType       task_system.TaskType `json:"task_type"`        // 任务的类型
	Token          string               `json:"token"`            // 任务的 token
	FinishDataSize int                  `json:"finish_data_size"` // 任务数据的大小
	ApiKey         string               `json:"api_key"`          // 身份密钥
	WorkerApiKey   string               `json:"worker_api_key"`   // worker 身份密钥
}

type SetOneTaskDoneReq struct {
	TaskType     task_system.TaskType `json:"task_type"`      // 任务的类型
	Token        string               `json:"token"`          // 任务的 token
	ApiKey       string               `json:"api_key"`        // 身份密钥
	WorkerApiKey string               `json:"worker_api_key"` // worker 身份密钥
}

// ----------------------------------------------

// AddMachineTranslationTaskPackageReq 添加一个任务的请求
type AddMachineTranslationTaskPackageReq struct {
	ImdbId  string `json:"imdb_id"`
	IsMovie bool   `json:"is_movie"` // 是电影还是连续剧
	Season  int    `json:"season"`   // 电影则是 -1
	Episode int    `json:"episode"`  // 连续剧则是 -1

	IsAudioOrSRT bool   `json:"is_audio_or_srt"` // 是音频还是字幕
	FileSha256   string `json:"file_sha256"`     // 文件的 SHA256
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
	Token         string `json:"token"`           // 针对这次任务的 token，需要使用来标记已经完成任务
}

type SetFirstPackageTaskDoneReq struct {
	TaskPackageId string `json:"task_package_id"` // 任务包的ID
	Token         string `json:"token"`           // 针对这次任务的 token，需要使用来标记已经完成任务
	ApiKey        string `json:"api_key"`         // 身份密钥
}

type GetTaskPackageInfoReq struct {
	TaskPackageId string `json:"task_package_id"` // 任务包的ID
	ApiKey        string `json:"api_key"`         // 身份密钥
}

type GetTaskPackageInfoResp struct {
	Status                   int                `json:"status"`                      // 任务的状态 0 失败，1 成功
	Message                  string             `json:"message"`                     // 任务的状态信息
	TaskPackageStatus        task_system.Status `json:"task_package_status"`         // 任务包的总状态
	AudioToSubtitleStatus    task_system.Status `json:"audio_to_subtitle_status"`    // 音频转字幕的状态，如果总任务包是上传的字幕，那么这个状态就是完成
	SplitTaskStatus          task_system.Status `json:"split_task_status"`           // 分割任务的状态
	TranslationTaskCount     int                `json:"translation_task_count"`      // 翻译任务的总数量
	TranslationTaskDoneCount int                `json:"translation_task_done_count"` // 翻译任务完成的数量
}

type GetTranslatedResultResp struct {
	Status            int    `json:"status"`              // 任务的状态 0 失败，1 成功
	Message           string `json:"message"`             // 任务的状态信息
	ResultDownloadUrl string `json:"result_download_url"` // 结果的下载地址
}

// ----------------------------------------------

// ConfirmSplitTaskInfoReq 确认分割任务的信息
type ConfirmSplitTaskInfoReq struct {
	Token        string `json:"token"`          // 针对这次任务的 token
	SrcDataSizes []int  `json:"src_data_sizes"` // 每个分割后的文件的大小
	ApiKey       string `json:"api_key"`        // 身份密钥
	WorkerApiKey string `json:"worker_api_key"` // worker 身份密钥
}

type GetOneSplitPartUploadURLReq struct {
	Token        string `json:"token"`          // 针对这次任务的 token
	PartIndex    int    `json:"part_index"`     // 分割后的文件的索引 任务的分片索引，从 0 开始
	PartCount    int    `json:"part_count"`     // 分割后的文件的总数
	ApiKey       string `json:"api_key"`        // 身份密钥
	WorkerApiKey string `json:"worker_api_key"` // worker 身份密钥
}

type GetOneSplitPartUploadURLResp struct {
	Status    int    `json:"status"`     // 任务的状态 0 失败，1 成功
	Message   string `json:"message"`    // 任务的状态信息
	UploadUrl string `json:"upload_url"` // 上传文件的URL
}

// ----------------------------------------------

type GetAllMergeTaskReq struct {
	Token        string `json:"token"`          // 针对这次任务的 token
	ApiKey       string `json:"api_key"`        // 身份密钥
	WorkerApiKey string `json:"worker_api_key"` // worker 身份密钥
}

type GetAllMergeTaskResp struct {
	Status                       int      `json:"status"`                           // 任务的状态 0 失败，1 成功
	Message                      string   `json:"message"`                          // 任务的状态信息
	AllTranslatedDataDownloadUrl []string `json:"all_translated_data_download_url"` // 所有翻译后的数据的下载地址
	OrgSrtDownloadUrl            string   `json:"org_srt_download_url"`             // 原始的字幕的下载地址
}

// ----------------------------------------------

type WorkerPingPongReq struct {
	TaskType     task_system.TaskType `json:"task_type"`      // 任务的类型
	ApiKey       string               `json:"api_key"`        // 身份密钥
	WorkerApiKey string               `json:"worker_api_key"` // worker 身份密钥
}

type WorkerPingPongResp struct {
	Status  int    `json:"status"`  // 任务的状态 0 失败，1 成功
	Message string `json:"message"` // 任务的状态信息
}

type WorkerCountReq struct {
	ApiKey string `json:"api_key"` // 身份密钥
}

type WorkerCountResp struct {
	Status       int      `json:"status"`        // 任务的状态 0 失败，1 成功
	Message      string   `json:"message"`       // 任务的状态信息
	WorkerTypes  []string `json:"worker_types"`  // worker 的类型
	WorkerCounts []int    `json:"worker_counts"` // worker 的数量
}
