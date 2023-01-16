package models

import "gorm.io/gorm"

type ZiMuKuTV struct {
	gorm.Model
	Season     int `gorm:"column:season;type:int(11);not null" json:"season"`   // 季数
	Episode    int `gorm:"column:episode;type:int(11);not null" json:"episode"` // 集数
	ZiMuKuInfo `gorm:"embedded"`
}
