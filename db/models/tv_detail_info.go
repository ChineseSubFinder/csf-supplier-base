package models

import (
	tmdb "github.com/cyruzin/golang-tmdb"
	"github.com/jinzhu/now"
	"gorm.io/gorm"
	"time"
)

// TvDetailInfo 连续的扩展信息
type TvDetailInfo struct {
	gorm.Model
	IMDBid           string `gorm:"column:imdb_id;unique;primary_key;type:varchar(20);not null"`
	FirstAirDate     string `gorm:"column:first_air_date;type:varchar(20);not null"`
	NumberOfSeasons  int    `gorm:"column:number_of_seasons;type:int;not null"`
	NumberOfEpisodes int    `gorm:"column:number_of_episodes;type:int;not null"`
	InProduction     bool   `gorm:"column:in_production;type:tinyint(1);not null;default:0"`
	Status           string `gorm:"column:status;type:varchar(20);not null"`
}

func NewTvDetailInfo(IMDBid string, tmdbTvDetails *tmdb.TVDetails) *TvDetailInfo {
	return &TvDetailInfo{
		IMDBid:           IMDBid,
		FirstAirDate:     tmdbTvDetails.FirstAirDate,
		NumberOfSeasons:  tmdbTvDetails.NumberOfSeasons,
		NumberOfEpisodes: tmdbTvDetails.NumberOfEpisodes,
		InProduction:     tmdbTvDetails.InProduction,
		Status:           tmdbTvDetails.Status,
	}
}

func (t TvDetailInfo) GetFirstAirDate() time.Time {
	tm, err := now.Parse(t.FirstAirDate)
	if err != nil {
		return time.Time{}
	}
	return tm
}
