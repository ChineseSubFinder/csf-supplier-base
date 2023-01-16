package models

import "gorm.io/gorm"

type HotMovie struct {
	gorm.Model
	HotMedia
}
