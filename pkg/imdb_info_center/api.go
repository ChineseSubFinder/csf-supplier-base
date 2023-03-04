package imdb_info_center

import (
	"github.com/ChineseSubFinder/csf-supplier-base/pkg/imdb_info_center/dao"
	models2 "github.com/ChineseSubFinder/csf-supplier-base/pkg/imdb_info_center/models"
	"github.com/jinzhu/now"
	"strings"
	"time"
)

// IsMovie 判断这个 ID 是电影还是连续剧。因为本地构建的数据库只有，电影和连续剧，不会有其他的类型
// return：这个 ID 是否在数据库中找到，true 电影，false 连续剧
func IsMovie(imdbId string) (bool, bool) {

	var tbs []models2.TitleBasic
	dao.Get().Where("tconst = ?", imdbId).Find(&tbs)

	if len(tbs) == 0 {
		return false, false
	}

	lowTypeName := strings.ToLower(tbs[0].TitleType)
	if strings.Contains(lowTypeName, "movie") == true ||
		strings.Contains(lowTypeName, "tvspecial") == true {
		return true, true
	} else {
		return true, false
	}
}

func EpsId2TvId(epsId string) (bool, string, int, int) {

	// 这里首先要考虑，这个 epsID 是否可能是 TV 的 ID
	// 如果是，那么就直接返回这个 ID，-1, -1
	// 如果不是，那么就要去查询这个 epsID 的父级 TV ID
	var tbs []models2.TitleEpisode
	dao.Get().Where("parent_imdb_id = ?", epsId).Find(&tbs)
	if len(tbs) == 0 {
		// 说明这个 epsID 不是 TV 的 ID，继续查找
		dao.Get().Where("tconst = ?", epsId).Find(&tbs)

		if len(tbs) == 0 {
			return false, "", -1, -1
		}

		return true, tbs[0].ParentTConst, tbs[0].SeasonNumber, tbs[0].EpisodeNumber
	} else {
		// 说明这个 epsID 是 TV 的 ID，直接返回
		return true, epsId, -1, -1
	}
}

// MovieIsRecent2Years 这个电影是否是近两年的电影
func MovieIsRecent2Years(imdbId string) bool {

	var tbs []models2.TitleBasic
	dao.Get().Where("tconst = ?", imdbId).Find(&tbs)
	if len(tbs) == 0 {
		return false
	}
	parseTime, err := now.Parse(tbs[0].StartYear)
	if err != nil {
		return false
	}
	nTime := time.Now()
	// 2022					2022-2=2020
	if parseTime.Year() >= nTime.Year()-2 {
		return true
	}

	return false
}

func GetYearInfo(imdbId string) (bool, YearInfo) {

	var tbs []models2.TitleBasic
	dao.Get().Where("tconst = ?", imdbId).Find(&tbs)
	if len(tbs) == 0 {
		return false, YearInfo{}
	}

	return true, YearInfo{tbs[0].StartYear, tbs[0].EndYear}
}

type YearInfo struct {
	Start string
	End   string
}

func (y YearInfo) StartYear() int {
	if y.Start == "" {
		return 0
	}
	parseTime, err := now.Parse(y.Start)
	if err != nil {
		return 0
	}
	return parseTime.Year()
}

func (y YearInfo) EndYear() int {
	if y.End == "" {
		return 0
	}
	parseTime, err := now.Parse(y.End)
	if err != nil {
		return 0
	}
	return parseTime.Year()
}