package models

type HouseKeepingError struct {
	HouseKeeping
	SaveRelativePath string `gorm:"column:save_relative_path;type:varchar(255);not null" json:"save_relative_path"`    // 保存的相对路径，包含文件名，可能是一个压缩包，也可能是一个字幕文件 /movie/2020/12/12/177838.srt
	UnzipError       bool   `gorm:"column:unzip_error;type:tinyint(1);not null;default:0" json:"unzip_error"`          // 是否是解压缩失败
	SubtitleExtType  int    `gorm:"column:subtitle_ext_type;type:int(11);not null;default:0" json:"subtitle_ext_type"` // 里面字幕类型，ass srt ssa sup bdn sst 等
}
