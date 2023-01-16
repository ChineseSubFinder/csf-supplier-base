package models

import "gorm.io/gorm"

type TopTv struct {
	gorm.Model
	HotMedia
}
