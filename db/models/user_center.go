package models

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

type UserCenter struct {
	gorm.Model
	TelegramUserID int64  `gorm:"column:telegram_user_id;type:bigint;not null;uniqueIndex;"` // Telegram 用户 ID
	ApiKey         string `gorm:"column:api_key;type:varchar(128);not null;uniqueIndex;"`    // 密钥,uuid
	KeyInfo
}

type KeyInfo struct {
	ActiveTime                   int64 `gorm:"column:active_time;type:bigint;not null"`                            // 什么时候激活的，续期的时候，这个时间也会更新
	ExpirationTime               int64 `gorm:"column:expiration_time;type:bigint;not null"`                        // 过期时间（具体那个时间），如果过期，那么这个 API Key 就无法使用
	DailyDownloadLimit           int   `gorm:"column:daily_download_limit;type:int;not null;default:50"`           // 每天的下载限制
	DailyUploadLimit             int   `gorm:"column:daily_upload_limit;type:int;not null;default:50"`             // 每天的上传限制
	DailyMachineTranslationLimit int   `gorm:"column:daily_machine_translation_limit;type:int;not null;default:0"` // 每天的机翻限制
}

func NewUserCenter(telegramUserID int64, apikey string) *UserCenter {
	return &UserCenter{
		TelegramUserID: telegramUserID,
		ApiKey:         apikey,
		KeyInfo:        *NewKeyInfo(50, 100),
	}
}

func NewKeyInfo(dailyDownloadLimit, dailyUploadLimit int) *KeyInfo {

	nowTime := time.Now()
	return &KeyInfo{
		ActiveTime:         nowTime.Unix(),
		ExpirationTime:     nowTime.AddDate(1, 0, 0).Unix(),
		DailyDownloadLimit: dailyDownloadLimit,
		DailyUploadLimit:   dailyUploadLimit,
	}
}

func (k KeyInfo) GetActiveTime() time.Time {
	return time.Unix(k.ActiveTime, 0)
}

func (k KeyInfo) GetExpirationTime() time.Time {
	return time.Unix(k.ExpirationTime, 0)
}

func (k KeyInfo) GetDailyDownloadLimit() string {
	return fmt.Sprintf("%d", k.DailyDownloadLimit)
}
