package tmdb_api

import (
	"fmt"
	"github.com/ChineseSubFinder/csf-supplier-base/pkg/common"
	"github.com/ChineseSubFinder/csf-supplier-base/pkg/settings"
	"github.com/allanpk716/rod_helper"
	tmdb "github.com/cyruzin/golang-tmdb"
	"github.com/sirupsen/logrus"
	"strconv"
	"sync"
)

type TMDBApi struct {
	l          *logrus.Logger
	apiKey     string
	tmdbClient *tmdb.Client
	locker     sync.Mutex
}

func NewTmdbHelper(l *logrus.Logger, apiKey string) (*TMDBApi, error) {

	tmdbClient, err := tmdb.Init(apiKey)
	if err != nil {
		err = fmt.Errorf("error initializing tmdb client: %s", err)
		return nil, err
	}

	t := TMDBApi{
		l:          l,
		apiKey:     apiKey,
		tmdbClient: tmdbClient,
	}
	t.ReSetClientConfig()
	return &t, nil
}

func (t *TMDBApi) Alive() bool {

	t.locker.Lock()
	defer t.locker.Unlock()
	options := make(map[string]string)
	options["language"] = "en-US"
	searchMulti, err := t.tmdbClient.GetSearchMulti("Dexter", options)
	if err != nil {
		t.l.Errorln("GetSearchMulti", err)
		return false
	}
	t.l.Infoln("Tmdb Api is Alive", searchMulti.TotalResults)
	return true
}

// GetInfo 获取视频的信息 idType: imdb_id or tmdb_id
func (t *TMDBApi) GetInfo(iD string, idType IdType, isMovieOrSeries, isQueryEnOrCNInfo bool) (outFindByID *tmdb.FindByID, err error) {

	t.locker.Lock()
	defer t.locker.Unlock()
	// 查询的参数
	options := make(map[string]string)
	if isQueryEnOrCNInfo == true {
		options["language"] = "en-US"
	} else {
		options["language"] = "zh-CN"
	}
	if idType == Imdb {

		options["external_source"] = Imdb.String()
		outFindByID, err = t.tmdbClient.GetFindByID(iD, options)
		if err != nil {
			return nil, fmt.Errorf("error getting tmdb info by id = %s: %s", iD, err)
		}
	} else if idType == Tmdb {

		intVar, err := strconv.Atoi(iD)
		if err != nil {
			return nil, fmt.Errorf("error converting tmdb id = %s to int: %s", iD, err)
		}

		if isMovieOrSeries == true {
			movieDetails, err := t.tmdbClient.GetMovieDetails(intVar, options)
			if err != nil {
				return nil, fmt.Errorf("error getting tmdb movie details by id = %s: %s", iD, err)
			}
			outFindByID = &tmdb.FindByID{
				MovieResults: []struct {
					Adult            bool    `json:"adult"`
					BackdropPath     string  `json:"backdrop_path"`
					GenreIDs         []int64 `json:"genre_ids"`
					ID               int64   `json:"id"`
					OriginalLanguage string  `json:"original_language"`
					OriginalTitle    string  `json:"original_title"`
					Overview         string  `json:"overview"`
					PosterPath       string  `json:"poster_path"`
					ReleaseDate      string  `json:"release_date"`
					Title            string  `json:"title"`
					Video            bool    `json:"video"`
					VoteAverage      float32 `json:"vote_average"`
					VoteCount        int64   `json:"vote_count"`
					Popularity       float32 `json:"popularity"`
				}{
					{
						Adult:            movieDetails.Adult,
						BackdropPath:     movieDetails.BackdropPath,
						ID:               movieDetails.ID,
						OriginalLanguage: movieDetails.OriginalLanguage,
						OriginalTitle:    movieDetails.OriginalTitle,
						Overview:         movieDetails.Overview,
						PosterPath:       movieDetails.PosterPath,
						ReleaseDate:      movieDetails.ReleaseDate,
						Title:            movieDetails.Title,
						Video:            movieDetails.Video,
						VoteAverage:      movieDetails.VoteAverage,
						VoteCount:        movieDetails.VoteCount,
						Popularity:       movieDetails.Popularity,
					},
				},
			}
		} else {
			tvDetails, err := t.tmdbClient.GetTVDetails(intVar, options)
			if err != nil {
				return nil, fmt.Errorf("error getting tmdb tv details by id = %s: %s", iD, err)
			}
			outFindByID = &tmdb.FindByID{
				TvResults: []struct {
					OriginalName     string   `json:"original_name"`
					ID               int64    `json:"id"`
					Name             string   `json:"name"`
					VoteCount        int64    `json:"vote_count"`
					VoteAverage      float32  `json:"vote_average"`
					FirstAirDate     string   `json:"first_air_date"`
					PosterPath       string   `json:"poster_path"`
					GenreIDs         []int64  `json:"genre_ids"`
					OriginalLanguage string   `json:"original_language"`
					BackdropPath     string   `json:"backdrop_path"`
					Overview         string   `json:"overview"`
					OriginCountry    []string `json:"origin_country"`
					Popularity       float32  `json:"popularity"`
				}{
					{
						OriginalName:     tvDetails.OriginalName,
						ID:               tvDetails.ID,
						Name:             tvDetails.Name,
						VoteCount:        tvDetails.VoteCount,
						VoteAverage:      tvDetails.VoteAverage,
						FirstAirDate:     tvDetails.FirstAirDate,
						PosterPath:       tvDetails.PosterPath,
						OriginalLanguage: tvDetails.OriginalLanguage,
						BackdropPath:     tvDetails.BackdropPath,
						Overview:         tvDetails.Overview,
						OriginCountry:    tvDetails.OriginCountry,
						Popularity:       tvDetails.Popularity,
					},
				},
			}
		}

	}

	return outFindByID, nil
}

// ConvertId 目前仅仅支持 TMDB ID 转 IMDB ID，iD：TMDB ID，idType：tmdb
func (t *TMDBApi) ConvertId(iD string, idType IdType, isMovieOrSeries bool) (convertIdResult *ConvertIdResult, err error) {

	t.locker.Lock()
	defer t.locker.Unlock()
	if idType == Imdb {
		return nil, fmt.Errorf("imdb id type is not supported")
	} else if idType == Tmdb {
		var intVar int
		intVar, err = strconv.Atoi(iD)
		if err != nil {
			return nil, fmt.Errorf("error converting tmdb id = %s to int: %s", iD, err)
		}
		options := make(map[string]string)
		if isMovieOrSeries == true {
			movieExternalIDs, err := t.tmdbClient.GetMovieExternalIDs(intVar, options)
			if err != nil {
				return nil, err
			}
			convertIdResult = &ConvertIdResult{
				ImdbID: movieExternalIDs.IMDbID,
				TmdbID: iD,
			}

			return convertIdResult, nil
		} else {
			tvExternalIDs, err := t.tmdbClient.GetTVExternalIDs(intVar, options)
			if err != nil {
				return nil, err
			}

			convertIdResult = &ConvertIdResult{
				ImdbID: tvExternalIDs.IMDbID,
				TmdbID: iD,
				TvdbID: fmt.Sprintf("%d", tvExternalIDs.TVDBID),
			}

			return convertIdResult, nil
		}
	} else {
		return nil, fmt.Errorf("id type is not supported: " + idType.String())
	}
}

// GetPopularMovie 获取热门电影的列表
func (t *TMDBApi) GetPopularMovie(pageIndex string, isQueryEnOrCNInfo bool) (moviePopular *tmdb.MoviePopular, err error) {

	t.locker.Lock()
	defer t.locker.Unlock()
	// 查询的参数
	options := make(map[string]string)
	if isQueryEnOrCNInfo == true {
		options["language"] = "en-US"
	} else {
		options["language"] = "zh-CN"
	}
	options["page"] = pageIndex
	moviePopular, err = t.tmdbClient.GetMoviePopular(options)
	if err != nil {
		return
	}

	return
}

// GetPopularTV 获取热门连续剧的列表
func (t *TMDBApi) GetPopularTV(pageIndex string, isQueryEnOrCNInfo bool) (tvPopular *tmdb.TVPopular, err error) {

	t.locker.Lock()
	defer t.locker.Unlock()
	// 查询的参数
	options := make(map[string]string)
	if isQueryEnOrCNInfo == true {
		options["language"] = "en-US"
	} else {
		options["language"] = "zh-CN"
	}
	options["page"] = pageIndex
	tvPopular, err = t.tmdbClient.GetTVPopular(options)
	if err != nil {
		return
	}

	return
}

// GetTopRatedMovie 获取最高评分电影的列表
func (t *TMDBApi) GetTopRatedMovie(pageIndex string, isQueryEnOrCNInfo bool) (movieTopRated *tmdb.MovieTopRated, err error) {

	t.locker.Lock()
	defer t.locker.Unlock()
	// 查询的参数
	options := make(map[string]string)
	if isQueryEnOrCNInfo == true {
		options["language"] = "en-US"
	} else {
		options["language"] = "zh-CN"
	}
	options["page"] = pageIndex
	movieTopRated, err = t.tmdbClient.GetMovieTopRated(options)
	if err != nil {
		return
	}

	return
}

// GetTopRatedTV 获取热门连续剧的列表
func (t *TMDBApi) GetTopRatedTV(pageIndex string, isQueryEnOrCNInfo bool) (tvTopRated *tmdb.TVTopRated, err error) {

	t.locker.Lock()
	defer t.locker.Unlock()
	// 查询的参数
	options := make(map[string]string)
	if isQueryEnOrCNInfo == true {
		options["language"] = "en-US"
	} else {
		options["language"] = "zh-CN"
	}
	options["page"] = pageIndex
	tvTopRated, err = t.tmdbClient.GetTVTopRated(options)
	if err != nil {
		return
	}

	return
}

// GetTrending 趋势
func (t *TMDBApi) GetTrending(mediaType common.MediaType, timeWindow TimeWindow) (getTrending *tmdb.Trending, err error) {

	t.locker.Lock()
	defer t.locker.Unlock()
	getTrending, err = t.tmdbClient.GetTrending(mediaType.String(), timeWindow.String())
	if err != nil {
		return
	}

	return
}

func (t *TMDBApi) ReSetClientConfig() {
	t.locker.Lock()
	defer t.locker.Unlock()
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
		return
	}
	t.tmdbClient.SetClientConfig(*restyClient.GetClient())
	t.tmdbClient.SetClientAutoRetry()
}

func (t *TMDBApi) GetClient() *tmdb.Client {
	return t.tmdbClient
}
