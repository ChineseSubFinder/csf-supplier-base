package models

type MixMediaInfo struct {
	IMDBid      string `gorm:"column:imdb_id;unique;primary_key;type:varchar(20);not null"`
	TMDBid      string `gorm:"column:tmdb_id;type:varchar(20);not null" json:"tmdb_id"`
	TVDBid      string `gorm:"column:tvdb_id;type:varchar(20);not null" json:"tvdb_id"`
	DouBanId    string `gorm:"column:douban_id;type:varchar(20);not null" json:"douban_id"`
	NameCn      string `gorm:"column:name_cn;type:varchar(255);not null" json:"name_cn"`          // 中文名称 或者 译名
	NameEn      string `gorm:"column:name_en;type:varchar(255);not null" json:"name_en"`          // 英文名称 或者 译名
	NameOrg     string `gorm:"column:name_org;type:varchar(255);not null" json:"name_org"`        // 原始名称
	ReleaseTime string `gorm:"column:release_time;type:varchar(20);not null" json:"release_time"` // 播出的时间，如果是连续剧是第一次播出的时间 2019-01-01  2022-01-01

	IMDbScore           IMDBScore           `gorm:"embedded;embeddedPrefix:im_"`
	TMDBScore           TMDBScore           `gorm:"embedded;embeddedPrefix:tm_"`
	DouBanScore         DouBanScore         `gorm:"embedded;embeddedPrefix:db_"`
	MetaCriticsScore    MetaCriticScore     `gorm:"embedded;embeddedPrefix:mc_"`
	RottenTomatoesScore RottenTomatoesScore `gorm:"embedded;embeddedPrefix:rt_"`

	Count            int   `gorm:"column:count;type:int;not null"`                 // 爬取的次数
	LastCrawlingTime int64 `gorm:"column:last_crawling_time;type:bigint;not null"` // 上次爬取的时间
}

type IMDBScore struct {
	ScoreInfo
}

func NewIMDBScore(scoreInfo ScoreInfo) *IMDBScore {
	return &IMDBScore{ScoreInfo: scoreInfo}
}

type TMDBScore struct {
	ScoreInfo
}

func NewTMDBScore(scoreInfo ScoreInfo) *TMDBScore {
	return &TMDBScore{ScoreInfo: scoreInfo}
}

type DouBanScore struct {
	ScoreInfo
}

func NewDouBanScore(scoreInfo ScoreInfo) *DouBanScore {
	return &DouBanScore{ScoreInfo: scoreInfo}
}

// MetaCriticScore 是一个专门收集对于电影、电视节目、音乐专辑、游戏的评论的网站
type MetaCriticScore struct {
	ScoreInfo
}

func NewMetaCriticScore(scoreInfo ScoreInfo) *MetaCriticScore {
	return &MetaCriticScore{ScoreInfo: scoreInfo}
}

// RottenTomatoesScore 烂番茄，越高，近期关注值越大，值得爬取
type RottenTomatoesScore struct {
	ScoreInfo
}

func NewRottenTomatoesScore(scoreInfo ScoreInfo) *RottenTomatoesScore {
	return &RottenTomatoesScore{ScoreInfo: scoreInfo}
}

type ScoreInfo struct {
	Url      string  `gorm:"column:url;type:varchar(255);not null" json:"url"` // 评分的 URL
	Score    float32 `gorm:"column:score" json:"score"`                        // 评分
	Comments int     `gorm:"column:comments" json:"comments"`                  // 评论数量
}

func NewScoreInfo(url string, score float32, comments int) *ScoreInfo {
	return &ScoreInfo{Url: url, Score: score, Comments: comments}
}

const (
	IMDBMaxScore           = 10
	TMDBMaxScore           = 10
	DoubanMaxScore         = 10
	MetaCriticMaxScore     = 100
	RottenTomatoesMaxScore = 100
)
