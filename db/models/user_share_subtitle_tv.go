package models

type UserShareSubtitleTV struct {
	UserShareSubtitleInfo
	FullSeasonSha256 string `gorm:"column:full_season_sha256;type:char(64);index"` // 一整季的 SHA256，因为上传的时候是分散的文件，这里计算会比较特别，具体看详细的实现
	Season           int    `gorm:"column:season;type:int;index;not null"`         // 如果无法识别就是 -1
	Episode          int    `gorm:"column:episode;type:int;index;not null"`        // 如果无法识别就是 -1
}
