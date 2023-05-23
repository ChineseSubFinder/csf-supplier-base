package task_system

import (
	"time"
)

type TaskPackageInfo struct {
	ID        int64     `gorm:"column:id;type:bigint;primary_key;auto_increment;not null"` // 任务的 ID
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;index;not null"`            // 任务创建的时间
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;index;not null"`            // 任务更新的时间

	ImdbId         string `gorm:"column:imdb_id;type:char(20);index;not null"` // IMDB ID
	IsMovie        bool   `gorm:"column:is_movie;type:tinyint(1);index;not null;default:0"`
	Season         int    `gorm:"column:season;type:int;index;not null"`
	Episode        int    `gorm:"column:episode;type:int;index;not null"`
	TelegramUserID int64  `gorm:"column:telegram_user_id;type:bigint;not null;index;"`  // Telegram 用户 ID
	PackageID      string `gorm:"column:package_id;type:char(64);uniqueIndex;not null"` // 任务包 ID
	Status         Status `gorm:"column:status;type:tinyint unsigned;index;not null"`   // 任务包的状态

	IsAudioOrSRT  bool   `gorm:"column:is_audio_or_srt;type:tinyint(1);index;not null;default:0"` // 是音频还是字幕
	FileSha256    string `gorm:"column:file_sha256;type:char(64);index;not null"`                 // 文件的 SHA256
	FileSize      int    `gorm:"column:file_size;type:int;not null"`                              // 文件大小，单位：字节
	SrcDataRPath  string `gorm:"column:src_data_r_path;type:varchar(255);not null"`               // 源任务数据的相对路径，相对于 R2 存储
	SrcDataSize   int    `gorm:"column:src_data_size;type:int;not null"`                          // 源任务数据的大小，单位：字节
	UploadFileUrl string `gorm:"-" json:"upload_file_url"`                                        // 上传文件的 URL
	Token         string `gorm:"-" json:"token"`                                                  // 缓存的 token信息

	AudioSrcLanguage   string `gorm:"column:audio_src_language;type:varchar(10);not null"`  // 音频的源语言
	TranslatedLanguage string `gorm:"column:translated_language;type:varchar(10);not null"` // 期望的翻译后的语言
}

func NewTaskPackageInfo(imdbId string, isMovie bool,
	season int, episode int,
	telegramUserID int64,
	packageID string, status Status,
	isAudioOrSRT bool, fileSha256 string, fileSize int,
	audioSrcLanguage string, translatedLanguage string) *TaskPackageInfo {
	return &TaskPackageInfo{
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
		ImdbId: imdbId, IsMovie: isMovie,
		Season: season, Episode: episode,
		TelegramUserID: telegramUserID, PackageID: packageID, Status: status,
		IsAudioOrSRT: isAudioOrSRT, FileSha256: fileSha256, FileSize: fileSize,
		AudioSrcLanguage: audioSrcLanguage, TranslatedLanguage: translatedLanguage,
	}
}

func (ti *TaskPackageInfo) SetStatus(status Status) {
	ti.Status = status
	ti.UpdatedAt = time.Now()
}
