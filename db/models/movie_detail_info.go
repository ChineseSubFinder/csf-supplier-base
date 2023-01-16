package models

import (
	"github.com/jinzhu/now"
	"gorm.io/gorm"
	"time"
)

type MovieDetailInfo struct {
	gorm.Model
	IMDBid      string `gorm:"column:imdb_id;type:varchar(20);not null"`
	ReleaseDate string `gorm:"column:release_date;type:varchar(20);not null"`
	Status      string `gorm:"column:status;type:varchar(20);not null"`
}

func NewMovieDetailInfo(IMDBid string, releaseDate string, status string) *MovieDetailInfo {
	return &MovieDetailInfo{IMDBid: IMDBid, ReleaseDate: releaseDate, Status: status}
}

// GetReleaseDate 获取电影的上映日期
func (m MovieDetailInfo) GetReleaseDate() time.Time {
	tm, err := now.Parse(m.ReleaseDate)
	if err != nil {
		return time.Time{}
	}
	return tm
}
