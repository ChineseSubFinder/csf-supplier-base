package models

import "gorm.io/gorm"

type HouseKeeping struct {
	gorm.Model
	DownloadedSubId uint  `gorm:"column:downloaded_id;type:int;not null"`                                // 已下载的 id，需要对应相应的表的存储内容的 ID，相当于记录处理到哪个了
	WhichSite       int   `gorm:"column:which_site;type:int;not null"`                                   // 1: ZiMuKu
	ProcessTime     int64 `gorm:"column:process_time;type:bigint;not null;default:0" json:"upload_time"` // 处理的时间
}

type WhichSite int

const (
	ZiMuKu WhichSite = iota + 1
	Other            // 不確定或者無所謂的，就設置這個
)

func (w WhichSite) Index() int {
	return int(w)
}
