package models

import (
	"fmt"
	tmdb "github.com/cyruzin/golang-tmdb"
	"github.com/jinzhu/now"
	"gorm.io/gorm"
	"time"
)

type EpisodeDetailInfo struct {
	gorm.Model
	IMDBid        string `gorm:"column:imdb_id;type:varchar(20);not null"`
	AirDate       string `gorm:"column:air_date;type:varchar(20);not null"`
	SeasonNumber  int    `gorm:"column:season_number;type:int;not null"`
	EpisodeNumber int    `gorm:"column:episode_number;type:int;not null"`
}

func NewEpisodeDetailInfos(IMDBid string, nowSeasonInfo *tmdb.TVSeasonDetails) []EpisodeDetailInfo {

	episodeDetailInfos := make([]EpisodeDetailInfo, 0)
	for _, episode := range nowSeasonInfo.Episodes {
		episodeDetailInfos = append(episodeDetailInfos, EpisodeDetailInfo{
			IMDBid:        IMDBid,
			AirDate:       episode.AirDate,
			SeasonNumber:  episode.SeasonNumber,
			EpisodeNumber: episode.EpisodeNumber,
		})
	}
	return episodeDetailInfos
}

func (e EpisodeDetailInfo) GetAirDate() time.Time {
	tm, err := now.Parse(e.AirDate)
	if err != nil {
		return time.Time{}
	}
	return tm
}

func (e EpisodeDetailInfo) GetSeasonEpisode() string {
	return "S" + fmt.Sprintf("%02d", e.SeasonNumber) + "E" + fmt.Sprintf("%02d", e.EpisodeNumber)
}
