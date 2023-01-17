package need_crawling

import (
	"github.com/ChineseSubFinder/csf-supplier-base/db/dao"
	"github.com/ChineseSubFinder/csf-supplier-base/db/models"
	"github.com/ChineseSubFinder/csf-supplier-base/pkg/common"
	"github.com/ChineseSubFinder/csf-supplier-base/pkg/media_info_dealer"
	"github.com/ChineseSubFinder/csf-supplier-base/pkg/search_subtitles"
	"github.com/pkg/errors"
	"time"
)

// Movie 是否需要爬取，在判断这些复杂条件前， 一般会针对没有爬取过的优先进行一次，能执行这个函数的，可以认为已经爬取过一次了，那么就需要按里面的逻辑来
func Movie(dealer *media_info_dealer.Dealers, imdbId string, isHotToday, isTop bool, config common.IntervalConfig, alwaysOnline bool) (bool, error) {

	// 获取电影的信息
	movieDetailInfo, err := dealer.GetMovieDetailInfo(imdbId, alwaysOnline)
	if err != nil {
		return false, err
	}
	// 获取这个电影最后的爬取时间
	var movies []models.Movie
	dao.Get().Where("imdb_id = ?", imdbId).Find(&movies)
	if len(movies) == 0 {
		return false, errors.New("movie not found")
	}
	lastCrawlingTime := movies[0].LastCrawlingTime
	// int64 转换为 time.Time
	lastCrawlingTimeTime := time.Unix(lastCrawlingTime, 0)

	if isTop == true {
		// Top List 中的，也就是盖棺定论的
		jugTime := time.Duration(config.MovieConfig.IsTop) * time.Minute
		if time.Now().After(lastCrawlingTimeTime.Add(jugTime)) == true {
			return true, nil
		} else {
			return false, nil
		}
	}

	// 优先判断这个电影是否是近两年上映的
	airDate := movieDetailInfo.GetReleaseDate()
	j2YearsDate := airDate.AddDate(2, 0, 0)

	if j2YearsDate.After(time.Now()) == true {

		if isHotToday == true {
			// 如果是今日热门
			jugTime := time.Duration(config.MovieConfig.IsHot) * time.Minute
			if time.Now().After(lastCrawlingTimeTime.Add(jugTime)) == true {
				return true, nil
			} else {
				return false, nil
			}
		}

		// 电影是近两年上映的
		// 这个电影是否已经下载了字幕，且有中文字幕
		hasChineseSub := false
		subtitleMovies, err := search_subtitles.Movie(imdbId)
		if err != nil {
			return false, err
		}
		if len(subtitleMovies) > 0 {
			// 有字幕，那么需要看是否有中文字幕
			// 其中是否有中文字幕
			for _, movieSubtitle := range subtitleMovies {
				if movieSubtitle.Lang().HasChinese() == true {
					hasChineseSub = true
					break
				}
			}
		}

		jugTime := time.Duration(config.MovieConfig.RecentTwoYearsNoSub) * time.Minute
		if hasChineseSub == false {
			// 没有中文字幕
			jugTime = time.Duration(config.MovieConfig.RecentTwoYearsNoSub) * time.Minute
		} else {
			// 有中文字幕
			jugTime = time.Duration(config.MovieConfig.RecentTwoYearsHasSub) * time.Minute
		}
		if time.Now().After(lastCrawlingTimeTime.Add(jugTime)) == true {
			return true, nil
		} else {
			return false, nil
		}

	} else {
		// 电影不是近两年上映的，一般来说字幕该有就有了
		// 是否达到了爬取的时间间隔(那么爬取的时间间隔是 ManyYears)
		if time.Now().After(lastCrawlingTimeTime.Add(time.Duration(config.MovieConfig.ManyYears)*time.Minute)) == true {
			return true, nil
		} else {
			return false, nil
		}
	}
}

// TV 是否需要爬取，在判断这些复杂条件前， 一般会针对没有爬取过的优先进行一次，能执行这个函数的，可以认为已经爬取过一次了，那么就需要按里面的逻辑来
func TV(dealer *media_info_dealer.Dealers, imdbId string, isHotToday, isTop bool, config common.IntervalConfig, alwaysOnline bool) (bool, error) {

	// 获取连续剧的信息
	tvDetailInfo, err := dealer.GetTVDetailInfo(imdbId, alwaysOnline)
	if err != nil {
		return false, err
	}
	// 获取这个连续剧最后的爬取时间
	var tvs []models.Tv
	dao.Get().Where("imdb_id = ?", imdbId).Find(&tvs)
	if len(tvs) == 0 {
		return false, errors.New("tv not found")
	}
	lastCrawlingTime := tvs[0].LastCrawlingTime
	// int64 转换为 time.Time
	lastCrawlingTimeTime := time.Unix(lastCrawlingTime, 0)

	if isTop == true {
		// Top List 中的，也就是盖棺定论的
		jugTime := time.Duration(config.MovieConfig.IsTop) * time.Minute
		if time.Now().After(lastCrawlingTimeTime.Add(jugTime)) == true {
			return true, nil
		} else {
			return false, nil
		}
	}
	/*
		大方向要分为：
		1. 已经完结的、被砍掉的、被取消的
		2. 未完结的，
			1. 正在连载的，正在更新的
			2. 等待更新的，下一季的
	*/
	if tvDetailInfo.InProduction == false {
		// 被砍掉的、被取消的
		if time.Now().After(lastCrawlingTimeTime.Add(time.Duration(config.TvConfig.End)*time.Minute)) == true {
			return true, nil
		} else {
			return false, nil
		}
	} else {
		// 未完结的
		if tvDetailInfo.IsAir() == true {
			// 正在连载的，正在更新的
			if time.Now().After(lastCrawlingTimeTime.Add(time.Duration(config.TvConfig.OnTheAir)*time.Minute)) == true {
				return true, nil
			} else {
				return false, nil
			}
		} else {
			// 等待更新的，下一季的
			if time.Now().After(lastCrawlingTimeTime.Add(time.Duration(config.TvConfig.Wait4NextSeason)*time.Minute)) == true {
				return true, nil
			} else {
				return false, nil
			}
		}
	}
}
