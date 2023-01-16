package models

import "gorm.io/gorm"

/*
	首先通过这里来记录，到底有那些需要爬取的电影
*/
type Movie struct {
	gorm.Model
	MixMediaInfo
}

func NewMovie(mixMediaInfo MixMediaInfo) *Movie {
	return &Movie{MixMediaInfo: mixMediaInfo}
}
