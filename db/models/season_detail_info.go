package models

import (
	tmdb "github.com/cyruzin/golang-tmdb"
	"github.com/jinzhu/now"
	"gorm.io/gorm"
	"time"
)

type SeasonDetailInfo struct {
	gorm.Model
	IMDBid       string `gorm:"column:imdb_id;type:varchar(20);not null"`
	AirDate      string `gorm:"column:air_date;type:varchar(20);not null"`
	SeasonNumber int    `gorm:"column:season_number;type:int;not null"`
	EpisodeCount int    `gorm:"column:episode_count;type:int;not null"`
}

func NewSeasonDetailInfos(IMDBid string, tmdbTvDetails *tmdb.TVDetails) (seasonDetailInfos []SeasonDetailInfo) {

	seasonDetailInfos = make([]SeasonDetailInfo, 0)
	for _, season := range tmdbTvDetails.Seasons {
		seasonDetailInfos = append(seasonDetailInfos, SeasonDetailInfo{
			IMDBid:       IMDBid,
			AirDate:      season.AirDate,
			SeasonNumber: season.SeasonNumber,
			EpisodeCount: season.EpisodeCount,
		})
	}
	return
}

func (s SeasonDetailInfo) GetAirDate() time.Time {
	tm, err := now.Parse(s.AirDate)
	if err != nil {
		return time.Time{}
	}
	return tm
}
