package models

import "gorm.io/gorm"

// ManualHandling 需要人工处理的数据
type ManualHandling struct {
	gorm.Model
	Stage        Stage  `gorm:"column:stage;type:tinyint;not null" json:"stage"`                      // 阶段
	Url          string `gorm:"column:url;type:varchar(255);not null" json:"url"`                     // 需要处理的 URL
	UniqueString string `gorm:"column:unique_string;type:varchar(255);not null" json:"unique_string"` // 唯一字符串,算作备用字段，StageQueryInfo 才会用上
	Remark       string `gorm:"column:remark;type:varchar(255)" json:"remark"`                        // 备注
	IsProcessed  bool   `gorm:"column:is_processed;type:tinyint;not null" json:"is_processed"`        // 是否已经处理
}

func NewManualHandling(stage Stage, url string, remark string, isProcessed bool) *ManualHandling {
	return &ManualHandling{Stage: stage, Url: url, Remark: remark, IsProcessed: isProcessed}
}

type Stage int

const (
	StageSubscriber   Stage = iota + 1 // 订阅阶段
	StageHotList                       // 热门列表
	StageDownloader                    // 下载阶段
	StageHouseKeeping                  // 整理阶段
	StageQueryInfo                     // 查询信息阶段
	StageTopList                       // Top 列表，可能是 Top 250 类似的列表，这些列表可能比较固定
	StageTrendingList                  // Trending 列表，近期的趋势列表，这个的结果会插入到 HotList 中
)

func (s Stage) String() string {
	switch s {
	case StageSubscriber:
		return "订阅阶段"
	case StageHotList:
		return "热门列表"
	case StageDownloader:
		return "下载阶段"
	case StageHouseKeeping:
		return "整理阶段"
	case StageQueryInfo:
		return "查询信息阶段"
	default:
		return "未知阶段"
	}
}
