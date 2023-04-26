package models

type UserShareSubtitleTV struct {
	UserShareSubtitleInfo
	Season  int `gorm:"column:season;type:int;index;not null"`  // 如果无法识别就是 -1
	Episode int `gorm:"column:episode;type:int;index;not null"` // 如果无法识别就是 -1
}
