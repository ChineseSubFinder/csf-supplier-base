package models

import "gorm.io/gorm"

type FinishMovie struct {
	gorm.Model
	HotMedia
}
