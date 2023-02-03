package api_hub

type SearchMovieSubtitleReq struct {
	ImdbId string `from:"imdb_id"`
	ApiKey string `from:"api_key"`
}

type SearchTVSubtitleByEpsReq struct {
	ImdbId  string `from:"imdb_id"`
	Season  int    `from:"season"`
	Episode int    `from:"episode"`
	ApiKey  string `from:"api_key"`
}

type SearchTVSubtitleBySeasonReq struct {
	ImdbId string `from:"imdb_id"`
	Season int    `from:"season"`
	ApiKey string `from:"api_key"`
}

type SearchTVSubtitleBySeasonPackageReq struct {
	ImdbId          string `from:"imdb_id"`
	SeasonPackageID string `from:"season_package_id"`
	ApiKey          string `from:"api_key"`
}

type GetSubtitleDownloadLinkReq struct {
	ImdbId          string `from:"imdb_id"`
	IsMovie         bool   `from:"is_movie"`
	Season          int    `json:"season"`
	Episode         int    `json:"episode"`
	SeasonPackageID string `json:"season_package_id"`
	SubSha256       string `from:"sub_sha256"` // 文件的 SHA256
	Language        int    `from:"language"`
	ApiKey          string `from:"api_key"`
}

type SearchTVSubtitleBySeasonResp struct {
	Status           int      `json:"status"`                      // 0 失败，1 成功
	Message          string   `json:"message"`                     // 返回的信息，包括成功和失败的原因
	SeasonPackageIDs []string `json:"season_package_ids,optional"` // 字幕包的 ID 列表
}

type GetSubtitleDownloadLinkResp struct {
	Status       int    `json:"status"`        // 0 失败，1 成功
	Message      string `json:"message"`       // 返回的信息，包括成功和失败的原因
	DownloadLink string `json:"download_link"` // 下载链接
}

type SearchSubtitleResp struct {
	Status   int        `json:"status"`            // 0 失败，1 成功
	Message  string     `json:"message"`           // 返回的信息，包括成功和失败的原因
	Subtitle []Subtitle `json:"subtitle,optional"` // 如果查询成功，返回的字幕信息
}

type Subtitle struct {
	SubSha256 string `json:"sub_sha256"` // 文件的 SHA256
	Title     string `json:"title"`      // 过滤关键词后的标题，如：The.Walking.Dead.S09E01.720p.HDTV.x264-AVS[eztv].srt
	Language  int    `json:"language"`   // 语言，参考 MyLanguage
	Ext       string `json:"ext"`        // 文件扩展名，如：srt、ass
	IsMovie   bool   `json:"is_movie"`   // 是电影还是连续剧
	Season    int    `json:"season"`     // 电影则是 -1
	Episode   int    `json:"episode"`    // 连续剧则是 -1
}
