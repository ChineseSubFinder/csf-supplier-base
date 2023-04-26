package models

type UserShareSubtitleTV struct {
	UserShareSubtitleInfo
	IsFullSeason     bool   `gorm:"column:is_full_season;type:tinyint(1);index;not null;default:0" json:"is_full_season"` // 是否是全季字幕，那么多个字幕会分开多个存储
	FullSeasonSha256 string `gorm:"column:full_season_sha256;type:char(64);index"`                                        // 一整季的 SHA256，因为上传的时候是分散的文件，这里计算会比较特别，具体看详细的实现
	Season           int    `gorm:"column:season;type:int;index;not null"`                                                // 如果无法识别就是 -1
	Episode          int    `gorm:"column:episode;type:int;index;not null"`                                               // 如果无法识别就是 -1
}
