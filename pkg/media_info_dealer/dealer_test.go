package media_info_dealer

import (
	"github.com/ChineseSubFinder/csf-supplier-base/pkg/settings"
	"github.com/ChineseSubFinder/csf-supplier-base/pkg/tmdb_api"
	"github.com/WQGroup/logger"
	"github.com/allanpk716/rod_helper"
	"testing"
)

func init() {
	rod_helper.InitFakeUA(true, settings.Get().CacheRootDirPath, settings.Get().TMDBConfig.TMDBHttpProxy)
}

func TestDealers_IMDBEpsId2TVId(t *testing.T) {

	var err error
	tmdbApi, err := tmdb_api.NewTmdbHelper(logger.GetLogger(), settings.Get().TMDBConfig.ApiKey)
	if err != nil {
		logger.Panicln("NewTmdbHelper err: ", err)
	}
	dealer := NewDealers(tmdbApi)

	//found, mixInfo := dealer.GetMoreInfoById(common.Movie, common.DouBan, "27605669")
	//if found == false {
	//	t.Error("GetMoreInfoById err")
	//}
	//println(mixInfo.NameCn)

	// EPS ID
	//epsID := "tt17663758"
	// main ID
	//mainID := "tt7197768"
	// 电影，黑亚当
	//mainID := "tt6443346"
	// 连续剧：无为大师
	//mainID := "tt4635276"
	// 连续剧：危险关系
	//mainID := "tt14792896"
	// 连续剧：大楼里只有谋杀 正确的 IMDB ID 是
	// http://www.imdb.com/title/tt12851524/
	// http://www.imdb.com/title/tt15425160/
	// http://www.imdb.com/title/tt21261218/
	mainID := "tt12851524"
	mainTVId, season, eps, err := dealer.IMDBEpsId2TVId(mainID)
	if err != nil {
		t.Fatal(err)
	}

	println(mainTVId)
	println(mainID)
	println(season)
	println(eps)
}

func TestDealers_GetMediaInfo(t *testing.T) {

	var err error
	tmdbApi, err := tmdb_api.NewTmdbHelper(logger.GetLogger(), settings.Get().TMDBConfig.ApiKey)
	if err != nil {
		logger.Panicln("NewTmdbHelper err: ", err)
	}
	dealer := NewDealers(tmdbApi)

	// FBI
	imdbID := "tt20195158"
	// 模范出租车
	imdbID = "tt13759970"
	mediaInfo, isMovie, err := dealer.GetMediaInfo(imdbID, tmdb_api.Imdb, false)
	if err != nil {
		t.Fatal(err)
	}
	println(mediaInfo.TitleCn)
	println(mediaInfo.TitleEn)
	println("isMovie: ", isMovie)
}

func TestDealers_GetTVDetailInfo(t *testing.T) {

	var err error
	tmdbApi, err := tmdb_api.NewTmdbHelper(logger.GetLogger(), settings.Get().TMDBConfig.ApiKey)
	if err != nil {
		logger.Panicln("NewTmdbHelper err: ", err)
	}
	// 模范出租车
	imdbID := "tt13759970"
	dealer := NewDealers(tmdbApi)
	tvDetailInfo, err := dealer.GetTVDetailInfo(imdbID, false)
	if err != nil {
		t.Fatal(err)
	}
	tvDetailInfo, err = dealer.InsertOrUpdateTVDetailInfo(imdbID)
	if err != nil {
		t.Fatal(err)
	}
	println(tvDetailInfo.IMDBid)
}

func TestDealers_GetMovieDetailInfo(t *testing.T) {

	var err error
	tmdbApi, err := tmdb_api.NewTmdbHelper(logger.GetLogger(), settings.Get().TMDBConfig.ApiKey)
	if err != nil {
		logger.Panicln("NewTmdbHelper err: ", err)
	}
	// 电影，黑亚当
	mainID := "tt6443346"
	dealer := NewDealers(tmdbApi)
	movieDetailInfo, err := dealer.GetMovieDetailInfo(mainID, false)
	if err != nil {
		t.Fatal(err)
	}
	movieDetailInfo, err = dealer.InsertOrUpdateMovieDetailInfo(mainID)
	if err != nil {
		t.Fatal(err)
	}
	println(movieDetailInfo.IMDBid)
}
