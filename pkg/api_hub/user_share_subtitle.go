package api_hub

// GenerateUploadURLReq 生成临时上传的URL 请求
type GenerateUploadURLReq struct {
	ImdbId    string `json:"imdb_id"`
	SubSha256 string `json:"sub_sha256"` // 文件的 SHA256
	Title     string `json:"title"`      // 过滤关键词后的标题，如：The.Walking.Dead.S09E01.720p.HDTV.x264-AVS[eztv].srt
	Language  int    `json:"language"`   // 语言，参考 MyLanguage
	Ext       string `json:"ext"`        // 文件扩展名，如：srt、ass
	IsMovie   bool   `json:"is_movie"`   // 是电影还是连续剧
	Season    int    `json:"season"`     // 电影则是 -1
	Episode   int    `json:"episode"`    // 连续剧则是 -1
	FileSize  int    `json:"file_size"`  // 文件大小，单位：字节
	Score     int    `json:"score"`      // 评分，参考 subtitle_mark.Score
	MarkType  int    `json:"mark_type"`  // 标记类型，参考 subtitle_mark.MarkType
	ApiKey    string `json:"api_key"`
}

// GenerateUploadURLResp 生成临时上传的URL 响应
type GenerateUploadURLResp struct {
	UploadURL string `json:"upload_url"` // 上传的URL
	Token     string `json:"api_key"`    // 如果上传成功用于回报的 token
}

type UploadSubtitleDoneReq struct {
	Token string `json:"api_key"` // 如果上传成功用于回报的 token
}

type UploadSubtitleDoneResp struct {
	Success bool   `json:"success"` // 是否成功
	Message string `json:"message"` // 如果失败的原因
}
