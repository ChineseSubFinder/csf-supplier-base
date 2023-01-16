package models

import "gorm.io/gorm"

type SubtitleMovie struct {
	gorm.Model
	SubtitleInfo `gorm:"embedded"`
}

type OrderSubtitleMovie []SubtitleMovie

func (d OrderSubtitleMovie) Len() int {
	return len(d)
}

func (d OrderSubtitleMovie) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

func (d OrderSubtitleMovie) Less(i, j int) bool {

	priorityI := d[i].Score*float32(d[i].Votes) + float32(d[i].DownloadTimes)
	priorityJ := d[j].Score*float32(d[j].Votes) + float32(d[j].DownloadTimes)

	return priorityI < priorityJ
}
