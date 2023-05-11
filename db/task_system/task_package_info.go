package task_system

import "gorm.io/gorm"

type TaskPackageInfo struct {
	gorm.Model
	ImdbId         string `gorm:"column:imdb_id;type:char(20);index;not null"` // IMDB ID
	IsMovie        bool   `gorm:"column:is_movie;type:tinyint(1);index;not null;default:0"`
	Season         int    `gorm:"column:season;type:int;index;not null"`
	Episode        int    `gorm:"column:episode;type:int;index;not null"`
	TelegramUserID int64  `gorm:"column:telegram_user_id;type:bigint;not null;uniqueIndex;"` // Telegram 用户 ID
	PackageID      string `gorm:"column:package_id;type:char(64);uniqueIndex;not null"`      // 任务包 ID
	Status         Status `gorm:"column:status;type:tinyint unsigned;index;not null"`        // 任务包的状态

	IsAudioOrSRT bool   `gorm:"column:is_audio_or_srt;type:tinyint(1);index;not null;default:0"` // 是音频还是字幕
	SubSha256    string `gorm:"column:sub_sha256;type:char(64);index;not null"`                  // 文件的 SHA256
	FileSize     int    `gorm:"column:file_size;type:int;not null"`                              // 文件大小，单位：字节

	AudioSrcLanguage   string `gorm:"column:audio_src_language;type:varchar(10);not null"`  // 音频的源语言
	TranslatedLanguage string `gorm:"column:translated_language;type:varchar(10);not null"` // 期望的翻译后的语言
}
