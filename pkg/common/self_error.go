package common

import "github.com/pkg/errors"

var (
	ZiMuKuSearchKeyWordStep0DetailPageUrlNotFound   = errors.New("zimuku search keyword step0 not found, detail page url")
	ZiMuKuSearchKeyWordStep1IMDBIDNotFound          = errors.New("zimuku search keyword step1 not found, imdb id")
	ZiMuKuSearchKeyWordStep1IMDBAndDouBanIDNotFound = errors.New("zimuku search keyword step1 not found, imdb and douban id")
	ZiMuKuSearchKeyWordStep1DouBanIDNotFound        = errors.New("zimuku search keyword step1 not found, douban id")
	ZiMuKuSearchKeyWordStep1IMDBIDNotMatch          = errors.New("zimuku search keyword step1 not found, imdb id not match")
	ZiMuKuSearchKeyWordStep1LocalMediaInfoNotFound  = errors.New("zimuku search keyword step1 not found, local media info")
	ZiMuKuDownloadUrlStep2NotFound                  = errors.New("zimuku download url step2 not found")
	ZiMuKuDownloadUrlStep3NotFound                  = errors.New("zimuku download url step3 not found")
	ZiMuKuDownloadUrlDownFileFailed                 = errors.New("zimuku download url DownFile failed")
	ZiMuKuDownloadUrlDownFileOverLimit              = errors.New("zimuku download url DownFile over limit")
)

var (
	MediaInfoDealerIMDBEpsId2TVIdCanNotGetImdbId = errors.New("EpsId2TVId can not get imdb id")
	MediaInfoDealerTMDBEpsId2TVIdNotFound        = errors.New("EpsId2TVId not found media info from tmdb")

	MediaInfoDealerTMDBNotFoundGetMediaInfo = errors.New("GetMediaInfo not found media info from tmdb")

	MediaInfoDealerTMDBFoundMovieAndSeriesInfoDontKnowWhichOneIs = errors.New("Found movie and series info from tmdb, dont know which one is")
)

var (
	UnZipError = errors.New("unzip error")
)

var (
	SubtitleExtTypeNotSupported = errors.New("subtitle ext type not supported")
	SubtitleExtTypeIsPicture    = errors.New("subtitle ext type is picture")
	SubtitleExtTypeIsBluRay     = errors.New("subtitle ext type is blu ray")
)

var (
	GetSeasonAndEpisodeFromSubFileNameError = errors.New("get season and episode from sub file name error")
)
