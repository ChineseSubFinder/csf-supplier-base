package models

import "gorm.io/gorm"

type SubtitleTV struct {
	gorm.Model
	IsFullSeason     bool   `gorm:"column:is_full_season;type:tinyint(1);not null;default:0" json:"is_full_season"` // 是否是全季字幕，那么多个字幕会分开多个存储
	FullSeasonSha256 string `gorm:"column:full_season_sha256;type:char(64)"`                                        // 一整季压缩包文件的 SHA256
	Season           int    `gorm:"column:season;type:int;not null"`                                                // 如果无法识别就是 -1
	Episode          int    `gorm:"column:episode;type:int;not null"`                                               // 如果无法识别就是 -1
	CantParseName    bool   `gorm:"column:cant_parse_name;type:tinyint(1);not null;default:0"`                      // 是否无法识别字幕的名字提取 Season 和 Episode 信息
	SubtitleInfo     `gorm:"embedded"`
}

type OrderSubtitleTV []SubtitleTV

func (d OrderSubtitleTV) Len() int {
	return len(d)
}

func (d OrderSubtitleTV) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

func (d OrderSubtitleTV) Less(i, j int) bool {

	priorityI := d[i].Score*float32(d[i].Votes) + float32(d[i].DownloadTimes)
	priorityJ := d[j].Score*float32(d[j].Votes) + float32(d[j].DownloadTimes)

	return priorityI < priorityJ
}
