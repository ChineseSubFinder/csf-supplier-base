package models

import "gorm.io/gorm"

type ZiMuKuMovie struct {
	gorm.Model
	ZiMuKuInfo `gorm:"embedded"`
}
