package models

import "gorm.io/gorm"

type HotTV struct {
	gorm.Model
	HotMedia
}
