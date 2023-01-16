package models

import "gorm.io/gorm"

type TopMovie struct {
	gorm.Model
	HotMedia
}
