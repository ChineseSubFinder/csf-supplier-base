package api_hub

// GenerateUploadURLReq 生成临时上传的URL 请求
type GenerateUploadURLReq struct {
	ImdbId           string `json:"imdb_id"`
	SubSha256        string `json:"sub_sha256"`         // 文件的 SHA256
	Title            string `json:"title"`              // 过滤关键词后的标题，如：The.Walking.Dead.S09E01.720p.HDTV.x264-AVS[eztv].srt
	Language         int    `json:"language"`           // 语言，参考 MyLanguage
	Ext              string `json:"ext"`                // 文件扩展名，如：srt、ass
	IsMovie          bool   `json:"is_movie"`           // 是电影还是连续剧
	Season           int    `json:"season"`             // 电影则是 -1
	Episode          int    `json:"episode"`            // 连续剧则是 -1
	FileSize         int    `json:"file_size"`          // 文件大小，单位：字节
	Score            int    `json:"score"`              // 评分，参考 subtitle_mark.Score
	MarkType         int    `json:"mark_type"`          // 标记类型，参考 subtitle_mark.MarkType
	SaveRelativePath string `json:"save_relative_path"` // 保存的相对路径，包含文件名，会处理后一定是一个具体的字幕文件 /movie/2020/12/12/177838.srt
	ApiKey           string `json:"api_key"`
}

// GenerateUploadURLResp 生成临时上传的URL 响应
type GenerateUploadURLResp struct {
	UploadURL string `json:"upload_url"` // 上传的URL
	Token     string `json:"token"`      // 如果上传成功用于回报的 token
}

type UploadSubtitleDoneReq struct {
	ApiKey string `json:"api_key"`
	Token  string `json:"token"` // 如果上传成功用于回报的 token
}

type UploadSubtitleDoneResp struct {
	Success bool   `json:"success"` // 是否成功
	Token   string `json:"token"`   // 后续标记这一批都完成需要这个 token，电影无需
	Message string `json:"message"` // 如果失败的原因
}

type MarkUploadTVSubsDoneReq struct {
	ApiKey string   `json:"api_key"`
	Tokens []string `json:"tokens"` // 如果上传成功用于回报的 tokens
}

type MarkUploadTVSubsDoneResp struct {
	Success bool   `json:"success"` // 是否成功
	Message string `json:"message"` // 如果失败的原因
}

type GetGenerateUploadURLReq struct {
	IsMovie    bool   `json:"is_movie"` // 是电影还是连续剧
	VideoFPath string `json:"video_f_path"`
	SubFPath   string `json:"sub_f_path"`
	Season     int    `json:"season"`
	Episode    int    `json:"episode"`
}
