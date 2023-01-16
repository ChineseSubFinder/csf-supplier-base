package common

type IntervalConfig struct {
	MovieConfig MovieConfig
	TvConfig    TvConfig
}

type MovieConfig struct {
	IsHot                int // 如果在近期的热门列表中，每隔多少分钟检查一次
	IsTop                int // 如果在Top列表中，每隔多少分钟检查一次
	RecentTwoYearsNoSub  int // 两年内的电影，且没有字幕的，每隔多少分钟检查一次
	RecentTwoYearsHasSub int // 两年内的电影，且有字幕的，每隔多少分钟检查一次
	ManyYears            int // 已经很多年的电影（不是近两年的），每隔多少分钟检查一次
}

type TvConfig struct {
	End             int // 已经完结的连续剧，每隔多少分钟检查一次
	Wait4NextSeason int // 还在更新的连续剧，每隔多少分钟检查一次
	OnTheAir        int // 正在连载连续剧，每隔多少分钟检查一次
}
