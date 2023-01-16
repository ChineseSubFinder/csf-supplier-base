package models

type HotMedia struct {
	IMDBid           string `gorm:"column:imdb_id;unique;primary_key;type:varchar(20);not null"`
	Count            int    `gorm:"column:count;type:int;not null"`                 // 从热门列表中更新的次数
	LastCrawlingTime int64  `gorm:"column:last_crawling_time;type:bigint;not null"` // 上次爬取的时间
	IsHotToday       bool   `gorm:"column:is_hot_today;type:tinyint(1);not null"`   // 是否是今天的热门
	IsNeedSkip       bool   `gorm:"column:is_need_skip;type:tinyint(1);not null"`   // 是否需要跳过
}
