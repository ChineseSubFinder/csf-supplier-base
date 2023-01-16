package media_info_dealer

import (
	"fmt"
	"github.com/ChineseSubFinder/csf-supplier/internal/models"
	"github.com/WQGroup/logger"
	"time"
)

type MovieDetailInfo struct {
	models.MovieDetailInfo
}

type TvDetailInfo struct {
	models.TvDetailInfo
	SeasonDetailInfos  []models.SeasonDetailInfo
	EpisodeDetailInfos []models.EpisodeDetailInfo
}

// IsAir 是否正在更新、连载中（注意，等待下一季的不算）
func (t TvDetailInfo) IsAir() bool {

	// 整理出所有的季数，且每一季对应有那些集
	// 然后找到最新一集的播出时间，如果已经超过7天（相对于当前时间）了，那么就认为已经最新一季完结了
	// 如果最新一集的播出时间距离现在不到7天，那么就认为正在更新
	epsMap := make(map[string]models.EpisodeDetailInfo)
	// SxEx 的形式
	for i, episodeDetailInfo := range t.EpisodeDetailInfos {
		epsMap[episodeDetailInfo.GetSeasonEpisode()] = t.EpisodeDetailInfos[i]
	}
	// 找出最大的 Season Number
	var maxSeasonNumber int
	for _, seasonDetailInfo := range t.SeasonDetailInfos {
		if seasonDetailInfo.SeasonNumber > maxSeasonNumber {
			maxSeasonNumber = seasonDetailInfo.SeasonNumber
		}
	}
	// 获取最大 Season Number 下的最后一集的集数
	var maxSeasonNumberMaxEpisodeNumber int
	for _, episodeDetailInfo := range t.EpisodeDetailInfos {
		if episodeDetailInfo.SeasonNumber == maxSeasonNumber &&
			episodeDetailInfo.EpisodeNumber > maxSeasonNumberMaxEpisodeNumber {
			maxSeasonNumberMaxEpisodeNumber = episodeDetailInfo.EpisodeNumber
		}
	}
	// 获取最大 Season Number 下的最后一集的播出时间
	var lastEpisodeDetailInfo models.EpisodeDetailInfo
	seasonEpsString := fmt.Sprintf("S%02dE%02d", maxSeasonNumber, maxSeasonNumberMaxEpisodeNumber)
	lastEpisodeDetailInfo, found := epsMap[seasonEpsString]
	if found == false {
		logger.Panicln(seasonEpsString, "not found")
	}
	airDate := lastEpisodeDetailInfo.GetAirDate()
	// 如果最后一集的播出时间距离现在不到7天，那么就认为正在更新
	if airDate.AddDate(0, 0, 7).After(time.Now()) {
		return true
	}

	return false
}
