package media_info_dealer

import (
	"fmt"
	"github.com/ChineseSubFinder/csf-supplier/internal/dao"
	"github.com/ChineseSubFinder/csf-supplier/internal/models"
	"github.com/ChineseSubFinder/csf-supplier/pkg/common"
	"github.com/ChineseSubFinder/csf-supplier/pkg/imdb_info_center"
	"github.com/ChineseSubFinder/csf-supplier/pkg/settings"
	"github.com/ChineseSubFinder/csf-supplier/pkg/tmdb_api"
	"github.com/WQGroup/logger"
	"github.com/allanpk716/rod_helper"
	tmdb "github.com/cyruzin/golang-tmdb"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"strconv"
)

type Dealers struct {
	tmdbHelper  *tmdb_api.TMDBApi
	restyClient *resty.Client
}

func NewDealers(tmdbHelper *tmdb_api.TMDBApi) *Dealers {

	// 获取 http client 实例
	httpProxyUrl := ""
	if settings.Get().TMDBConfig.TMDBHttpProxyEnable == true && settings.Get().TMDBConfig.TMDBHttpProxy != "" {
		httpProxyUrl = settings.Get().TMDBConfig.TMDBHttpProxy
	}
	opt := rod_helper.NewHttpClientOptions(settings.Get().TimeConfig.GetOnePageTimeOut())
	opt.SetHttpProxy(httpProxyUrl)
	restyClient, err := rod_helper.NewHttpClient(opt)
	if err != nil {
		err = fmt.Errorf("error initializing resty client: %s", err)
		logger.Panicln(err)
	}

	return &Dealers{
		restyClient: restyClient,
		tmdbHelper:  tmdbHelper,
	}
}

func (d *Dealers) GetTmdbHelper() *tmdb_api.TMDBApi {
	return d.tmdbHelper
}

func (d *Dealers) GetRestyClient() *resty.Client {
	return d.restyClient
}

// ConvertId 目前仅仅支持 TMDB ID 转 IMDB ID，iD：TMDB ID，source：tmdb
func (d *Dealers) ConvertId(iD string, idSource tmdb_api.IdType, isMovieOrSeries bool) (convertIdResult *tmdb_api.ConvertIdResult, err error) {

	return d.tmdbHelper.ConvertId(iD, idSource, isMovieOrSeries)
}

// GetMediaInfo 通过用户自己的 tmdb api 查询媒体信息 "source"=imdb|tmdb  isMovie（预期的类型），返回的时候会带有查询出来的类型
func (d *Dealers) GetMediaInfo(id string, idSource tmdb_api.IdType, isMovie bool) (*MediaInfo, bool, error) {

	imdbId := ""
	var tmdbID int64
	if idSource == tmdb_api.Imdb {
		imdbId = id
	} else if idSource == tmdb_api.Tmdb {

	} else {
		return nil, false, errors.New("idSource is not support")
	}
	// 先查询英文信息，然后再查询中文信息
	findByIDEn, err := d.tmdbHelper.GetInfo(id, idSource, isMovie, true)
	if err != nil {
		return nil, false, fmt.Errorf("error while getting info from TMDB: %v", err)
	}
	findByIDCn, err := d.tmdbHelper.GetInfo(id, idSource, isMovie, false)
	if err != nil {
		return nil, false, fmt.Errorf("error while getting info from TMDB: %v", err)
	}

	OriginalTitle := ""
	OriginalLanguage := ""
	TitleEn := ""
	TitleCn := ""
	Year := ""
	var Vote float32
	var VoteCount int64
	outIsMovie := false
	// 电影
	/*
		理论上吧，填入这里查询的信息应该是电影，那么就在电影的部分应该可以得到信息
		实际上呢，因为爬取网站的原因，可能的情况是，写是电影，其实是连续剧，这个时候就会出现问题
		但是，这里其实已经把电影和连续剧的信息都爬取了，因为传入的是一个绝对的 ID 信息，这里如果都错，就没救了
		那么就认为，如果期望是电影，但是电影没得， 连续剧的信息有，那么就妥了，改为读取是连续剧的信息
		下面的连续剧判断逻辑也是如果，以此类推
	*/
	if len(findByIDEn.MovieResults) < 1 && len(findByIDEn.TvResults) < 1 {
		// 都没有找到
		// 再看看是否查询到的是一集的信息
		if len(findByIDEn.TvEpisodeResults) < 1 {
			return nil, false, common.MediaInfoDealerTMDBNotFoundGetMediaInfo
		}
		// 查询到的是一集的信息
		// 从这里获取 Show_id ,然后再通过 Conver ID 查询出剧集的总 IMDB ID
		showID := findByIDEn.TvEpisodeResults[0].ShowID
		var idResult *tmdb_api.ConvertIdResult
		idResult, err = d.tmdbHelper.ConvertId(fmt.Sprintf("%d", showID), tmdb_api.Tmdb, false)
		if err != nil {
			return nil, false, err
		}
		imdbId = idResult.ImdbID
		// 先查询英文信息，然后再查询中文信息
		findByIDEn, err = d.tmdbHelper.GetInfo(imdbId, tmdb_api.Imdb, isMovie, true)
		if err != nil {
			return nil, false, fmt.Errorf("error while getting info from TMDB: %v", err)
		}
		findByIDCn, err = d.tmdbHelper.GetInfo(imdbId, tmdb_api.Imdb, isMovie, false)
		if err != nil {
			return nil, false, fmt.Errorf("error while getting info from TMDB: %v", err)
		}
	}

	if len(findByIDEn.MovieResults) > 0 && len(findByIDEn.TvResults) > 0 {
		/*
			之所有这里会认为有可能有问题
			是因为如果是使用 IMDB ID 查询，那么就可能出现这种情况
			如果是使用 TMDB ID 查询，那么就不会出现这种情况
			那么核心点就是确认这个 IMDB ID 是电影还是连续剧，现在这里查找自己私有的数据库，看看是否有记录
		*/
		var found bool
		found, isMovie = imdb_info_center.IsMovie(imdbId)
		if found == false {
			// 说明这个 IMDB ID 在本地数据库找不到，那么就无法判断后续的逻辑
			return nil, false, common.MediaInfoDealerTMDBFoundMovieAndSeriesInfoDontKnowWhichOneIs
		}

		if isMovie == true {
			// 找到的是电影
			tmdbID = findByIDEn.MovieResults[0].ID
			OriginalTitle = findByIDEn.MovieResults[0].OriginalTitle
			OriginalLanguage = findByIDEn.MovieResults[0].OriginalLanguage
			TitleEn = findByIDEn.MovieResults[0].Title
			TitleCn = findByIDCn.MovieResults[0].Title
			Year = findByIDEn.MovieResults[0].ReleaseDate
			Vote = findByIDEn.MovieResults[0].VoteAverage
			VoteCount = findByIDEn.MovieResults[0].VoteCount
			outIsMovie = true
		} else {
			// 找到的是连续剧
			tmdbID = findByIDEn.TvResults[0].ID
			OriginalTitle = findByIDEn.TvResults[0].OriginalName
			OriginalLanguage = findByIDEn.TvResults[0].OriginalLanguage
			TitleEn = findByIDEn.TvResults[0].Name
			TitleCn = findByIDCn.TvResults[0].Name
			Year = findByIDEn.TvResults[0].FirstAirDate
			Vote = findByIDEn.TvResults[0].VoteAverage
			VoteCount = findByIDEn.TvResults[0].VoteCount
			outIsMovie = false
		}
	} else if len(findByIDEn.MovieResults) > 0 {
		// 找到的是电影
		tmdbID = findByIDEn.MovieResults[0].ID
		OriginalTitle = findByIDEn.MovieResults[0].OriginalTitle
		OriginalLanguage = findByIDEn.MovieResults[0].OriginalLanguage
		TitleEn = findByIDEn.MovieResults[0].Title
		TitleCn = findByIDCn.MovieResults[0].Title
		Year = findByIDEn.MovieResults[0].ReleaseDate
		Vote = findByIDEn.MovieResults[0].VoteAverage
		VoteCount = findByIDEn.MovieResults[0].VoteCount
		outIsMovie = true
	} else if len(findByIDEn.TvResults) > 0 {
		// 找到的是连续剧
		tmdbID = findByIDEn.TvResults[0].ID
		OriginalTitle = findByIDEn.TvResults[0].OriginalName
		OriginalLanguage = findByIDEn.TvResults[0].OriginalLanguage
		TitleEn = findByIDEn.TvResults[0].Name
		TitleCn = findByIDCn.TvResults[0].Name
		Year = findByIDEn.TvResults[0].FirstAirDate
		Vote = findByIDEn.TvResults[0].VoteAverage
		VoteCount = findByIDEn.TvResults[0].VoteCount
		outIsMovie = false
	}

	mediaInfo := &MediaInfo{
		ImdbId:           imdbId,
		TmdbId:           fmt.Sprintf("%d", tmdbID),
		OriginalTitle:    OriginalTitle,
		OriginalLanguage: OriginalLanguage,
		TitleEn:          TitleEn,
		TitleCn:          TitleCn,
		Year:             Year,
		Vote:             Vote,
		VoteCount:        int(VoteCount),
	}
	convId, err := d.tmdbHelper.ConvertId(fmt.Sprintf("%d", tmdbID), tmdb_api.Tmdb, isMovie)
	if err == nil {
		mediaInfo.TVdbId = convId.TvdbID
		if idSource == tmdb_api.Tmdb {
			// 因为是使用 TMDB ID 查询的，所以 IMDB ID 需要额外查询出来填充
			mediaInfo.ImdbId = convId.ImdbID
		}
	}
	return mediaInfo, outIsMovie, nil
}

// IMDBEpsId2TVId 从一集的 IMDB ID 查询这部连续剧的 ID，当然也会考虑，如果传入的就是连续剧的总 ID 的情况
// 返回： IMDB ID，season，eps
func (d *Dealers) IMDBEpsId2TVId(id string) (string, int, int, error) {

	// 优先查询私有 IMDB 数据库
	// 首先，要确认这个 IMDB ID 是连续剧的，还是电影的
	found, tvID, season, eps := imdb_info_center.EpsId2TvId(id)
	if found == true {
		// 如果找到了，那么就提前返回
		return tvID, season, eps, nil
	}
	// 先查询英文信息，然后再查询中文信息
	findByIDEn, err := d.tmdbHelper.GetInfo(id, tmdb_api.Imdb, false, true)
	if err != nil {
		return "", -1, -1, err
	}

	if len(findByIDEn.TvResults) > 0 {
		// 说明是总 ID
		return id, -1, -1, nil
	}
	// 这个查询的是 IMDB ID 的情况，且不是这一部电影或者连续剧的总 ID，而是一集的 ID
	// 那么 tv_episode_results 应该不为空
	if len(findByIDEn.TvEpisodeResults) < 1 {
		return "", -1, -1, common.MediaInfoDealerTMDBEpsId2TVIdNotFound
	}

	// 从这里获取 Show_id ,然后再通过 Conver ID 查询出剧集的总 IMDB ID
	showID := findByIDEn.TvEpisodeResults[0].ShowID
	var idResult *tmdb_api.ConvertIdResult
	idResult, err = d.tmdbHelper.ConvertId(fmt.Sprintf("%d", showID), tmdb_api.Tmdb, false)
	if err != nil {
		return "", -1, -1, err
	}
	// imdbId
	if idResult.ImdbID != "" {
		return idResult.ImdbID, findByIDEn.TvEpisodeResults[0].SeasonNumber, findByIDEn.TvEpisodeResults[0].EpisodeNumber, nil
	} else {
		return "", -1, -1, common.MediaInfoDealerIMDBEpsId2TVIdCanNotGetImdbId
	}
}

// GetMoreInfoById 从自己的数据库中查询数据
func (d *Dealers) GetMoreInfoById(mediaType common.MediaType, inIdType common.MediaIdType, id string) (bool, models.MixMediaInfo) {

	if mediaType == common.Movie {
		var movieInfos []models.Movie
		dao.Get().Where(inIdType.QueryString(), id).Find(&movieInfos)

		if len(movieInfos) > 0 {
			return true, movieInfos[0].MixMediaInfo
		}
		return false, models.MixMediaInfo{}

	} else if mediaType == common.TV {
		var tvInfos []models.Tv
		dao.Get().Where(inIdType.QueryString(), id).Find(&tvInfos)

		if len(tvInfos) > 0 {
			return true, tvInfos[0].MixMediaInfo
		}
		return false, models.MixMediaInfo{}
	} else {
		logger.Panicln("unknown media type")
	}

	return false, models.MixMediaInfo{}
}

// GetTVDetailInfo 注意这里传入的 ID 是总 ID
func (d *Dealers) GetTVDetailInfo(imdbId string, alwaysOnline bool) (*TvDetailInfo, error) {

	if alwaysOnline == true {
		// 就是要实时查询信息来使用
		return d.InsertOrUpdateTVDetailInfo(imdbId)
	}
	// 首先查询私有数据库是否存在
	var tvDetailInfos []models.TvDetailInfo
	dao.Get().Where("imdb_id = ?", imdbId).Find(&tvDetailInfos)
	if len(tvDetailInfos) > 0 {
		// 存在，那么继续查询剩余的信息
		var seasonDetailInfos []models.SeasonDetailInfo
		dao.Get().Where("imdb_id = ?", imdbId).Find(&seasonDetailInfos)

		var episodeDetailInfos []models.EpisodeDetailInfo
		dao.Get().Where("imdb_id = ?", imdbId).Find(&episodeDetailInfos)

		var outTvDetailInfo TvDetailInfo
		outTvDetailInfo.TvDetailInfo = tvDetailInfos[0]
		outTvDetailInfo.SeasonDetailInfos = seasonDetailInfos
		outTvDetailInfo.EpisodeDetailInfos = episodeDetailInfos

		return &outTvDetailInfo, nil
	} else {
		// 不存在，那么就从 TMDB 获取
		return d.InsertOrUpdateTVDetailInfo(imdbId)
	}
}

// InsertOrUpdateTVDetailInfo 从 TMDB 获取剧集的详细信息
func (d *Dealers) InsertOrUpdateTVDetailInfo(imdbId string) (*TvDetailInfo, error) {

	info, _, err := d.GetMediaInfo(imdbId, tmdb_api.Imdb, false)
	if err != nil {
		return nil, err
	}
	intVar := 0
	intVar, err = strconv.Atoi(info.TmdbId)
	if err != nil {
		return nil, err
	}
	// 查询 TV Detail
	options := make(map[string]string)
	tmdbTvDetails, err := d.tmdbHelper.GetClient().GetTVDetails(intVar, options)
	if err != nil {
		return nil, err
	}
	tvDetailInfo := models.NewTvDetailInfo(imdbId, tmdbTvDetails)
	// 从其中解析出 SeasonDetailInfo 出来
	seasonDetailInfos := models.NewSeasonDetailInfos(imdbId, tmdbTvDetails)
	// 查询集的信息
	episodeDetailInfos := make([]models.EpisodeDetailInfo, 0)
	for _, seasonDetailInfo := range seasonDetailInfos {

		var nowSeasonInfo *tmdb.TVSeasonDetails
		nowSeasonInfo, err = d.tmdbHelper.GetClient().GetTVSeasonDetails(intVar, seasonDetailInfo.SeasonNumber, options)
		if err != nil {
			return nil, err
		}

		nowEpisodeDetailInfos := models.NewEpisodeDetailInfos(imdbId, nowSeasonInfo)
		episodeDetailInfos = append(episodeDetailInfos, nowEpisodeDetailInfos...)
	}

	err = dao.Get().Transaction(func(tx *gorm.DB) error {

		// 保存 TVDetailInfo
		var queryTvDetailInfos []models.TvDetailInfo
		tx.Where("imdb_id = ?", imdbId).Find(&queryTvDetailInfos)
		if len(queryTvDetailInfos) > 0 {
			// 更新
			if err := tx.Model(&queryTvDetailInfos[0]).Updates(tvDetailInfo).Error; err != nil {
				return err
			}
		} else {
			// 插入
			if err := tx.Create(&tvDetailInfo).Error; err != nil {
				return err
			}
		}
		// 保存 SeasonDetailInfo，存在则更新，不存在则插入
		for _, seasonDetailInfo := range seasonDetailInfos {

			var querySeasonDetailInfos []models.SeasonDetailInfo
			tx.Where("imdb_id = ? and season_number = ?", imdbId, seasonDetailInfo.SeasonNumber).Find(&querySeasonDetailInfos)
			if len(querySeasonDetailInfos) > 0 {
				// 存在，那么更新
				if err := tx.Model(&querySeasonDetailInfos[0]).Updates(&seasonDetailInfo).Error; err != nil {
					return err
				}
			} else {
				// 不存在，那么插入
				if err := tx.Create(&seasonDetailInfo).Error; err != nil {
					return err
				}
			}
		}
		// 保存 EpisodeDetailInfo，存在则更新，不存在则插入
		for _, episodeDetailInfo := range episodeDetailInfos {
			var queryEpisodeDetailInfos []models.EpisodeDetailInfo
			tx.Where("imdb_id = ? and season_number = ? and episode_number = ?", imdbId, episodeDetailInfo.SeasonNumber, episodeDetailInfo.EpisodeNumber).Find(&queryEpisodeDetailInfos)
			if len(queryEpisodeDetailInfos) > 0 {
				// 存在，那么更新
				if err := tx.Model(&queryEpisodeDetailInfos[0]).Updates(&episodeDetailInfo).Error; err != nil {
					return err
				}
			} else {
				// 不存在，那么插入
				if err := tx.Create(&episodeDetailInfo).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	var outTvDetailInfo TvDetailInfo
	outTvDetailInfo.TvDetailInfo = *tvDetailInfo
	outTvDetailInfo.SeasonDetailInfos = seasonDetailInfos
	outTvDetailInfo.EpisodeDetailInfos = episodeDetailInfos

	return &outTvDetailInfo, nil
}

// GetMovieDetailInfo 从 TMDB 获取电影的详细信息
func (d *Dealers) GetMovieDetailInfo(imdbId string, alwaysOnline bool) (*MovieDetailInfo, error) {

	if alwaysOnline == true {
		return d.InsertOrUpdateMovieDetailInfo(imdbId)
	}
	// 首先查询私有数据库是否存在
	var movieDetailInfos []models.MovieDetailInfo
	dao.Get().Where("imdb_id = ?", imdbId).Find(&movieDetailInfos)
	if len(movieDetailInfos) > 0 {
		// 存在
		var outMovieDetailInfo MovieDetailInfo
		outMovieDetailInfo.MovieDetailInfo = movieDetailInfos[0]
		return &outMovieDetailInfo, nil
	} else {
		// 不存在，那么就从 TMDB 获取
		return d.InsertOrUpdateMovieDetailInfo(imdbId)
	}
}

// InsertOrUpdateMovieDetailInfo 从 TMDB 获取电影的详细信息
func (d *Dealers) InsertOrUpdateMovieDetailInfo(imdbId string) (*MovieDetailInfo, error) {

	var outMovieDetailInfo MovieDetailInfo
	info, _, err := d.GetMediaInfo(imdbId, tmdb_api.Imdb, false)
	if err != nil {
		return nil, err
	}
	intVar := 0
	intVar, err = strconv.Atoi(info.TmdbId)
	if err != nil {
		return nil, err
	}
	// 查询 TV Detail
	options := make(map[string]string)
	tmdbMovieDetails, err := d.tmdbHelper.GetClient().GetMovieDetails(intVar, options)
	if err != nil {
		return nil, err
	}
	err = dao.Get().Transaction(func(tx *gorm.DB) error {
		// 保存 MovieDetailInfo
		var queryMovieDetailInfos []models.MovieDetailInfo
		tx.Where("imdb_id = ?", imdbId).Find(&queryMovieDetailInfos)
		if len(queryMovieDetailInfos) > 0 {
			// 存在，那么更新
			queryMovieDetailInfos[0].ReleaseDate = tmdbMovieDetails.ReleaseDate
			queryMovieDetailInfos[0].Status = tmdbMovieDetails.Status

			if err := tx.Updates(&queryMovieDetailInfos[0]).Error; err != nil {
				return err
			}

			outMovieDetailInfo.MovieDetailInfo = queryMovieDetailInfos[0]
		} else {
			// 不存在，那么插入
			needCreateMovieDetailInfo := models.NewMovieDetailInfo(imdbId, tmdbMovieDetails.ReleaseDate, tmdbMovieDetails.Status)
			if err := tx.Create(needCreateMovieDetailInfo).Error; err != nil {
				return err
			}
			outMovieDetailInfo.MovieDetailInfo = *needCreateMovieDetailInfo
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &outMovieDetailInfo, nil
}
