package models

import "gorm.io/gorm"

/*
	首先通过这里来记录，到底有那些需要爬取的连续剧
*/
type Tv struct {
	gorm.Model
	MixMediaInfo `gorm:"embedded"`
}

func NewTv(mixMediaInfo MixMediaInfo) *Tv {
	return &Tv{MixMediaInfo: mixMediaInfo}
}
