package models

import (
	"github.com/ChineseSubFinder/csf-supplier-base/pkg/types/language"
	"path/filepath"
	"runtime"
	"strings"
)

type UserShareSubtitleInfo struct {
	ID               uint   `gorm:"primarykey"`
	TelegramUserID   int64  `gorm:"column:telegram_user_id;type:bigint;index;not null;"`                            // Telegram 用户 ID
	ImdbId           string `gorm:"column:imdb_id;type:char(20);index;not null"`                                    // IMDB ID，电影的 IMDB ID，电视剧的主 IMDB ID
	SubSha256        string `gorm:"column:sub_sha256;type:char(64);index"`                                          // 文件的 SHA256
	Title            string `gorm:"column:title;type:varchar(255);not null" json:"title"`                           // 字幕标题
	Language         int    `gorm:"column:language;type:int;index;not null"`                                        // 字幕语言,MyLanguage 参考内部的这个值
	SaveRelativePath string `gorm:"column:save_relative_path;type:varchar(255);not null" json:"save_relative_path"` // 保存的相对路径，包含文件名，会处理后一定是一个具体的字幕文件 /movie/2020/12/12/177838.srt
	Score            int    `gorm:"column:score;type:int;not null;default:0" json:"score"`                          // 评分，参考 subtitle_mark.Score
	MarkType         int    `gorm:"column:mark_type;type:int;not null;default:0" json:"mark_type"`                  // 标记类型，参考 subtitle_mark.MarkType
	UploadTime       int64  `gorm:"column:upload_time;type:bigint;not null;default:0" json:"upload_time"`           // 上传时间
}

func (s *UserShareSubtitleInfo) Ext() string {
	return filepath.Ext(s.SaveRelativePath)
}

func (s *UserShareSubtitleInfo) Lang() language.MyLanguage {
	return language.MyLanguage(s.Language)
}

func (s *UserShareSubtitleInfo) GetTitle(needFilterKeyWords []string) string {

	nowTitle := s.Title
	// 因为下载的文件名中，可能包含一些 [zmk.pw] 这样的信息，需要剔除
	for _, keyWord := range needFilterKeyWords {
		nowTitle = strings.ReplaceAll(nowTitle, keyWord, "")
	}
	return nowTitle
}

func (s *UserShareSubtitleInfo) R2StoreKey() string {

	// 将本地的 SaveRelativePath 路径转换为 r2 存储的 key

	// 这里有个梗，因为现在的数据都是在 Windows 上获取上传的，那么路径都是 Windows 的 \\ 格式
	// 在 Linux 上会有问题，所以如果是非 Windows 系统这里需要将 \\ 转换为 /
	// 判断当前的系统是否是 Windows

	nowSaveRelativePath := ""
	if runtime.GOOS == "windows" {
		// 当前是 Windows 系统
		nowSaveRelativePath = s.SaveRelativePath
	} else {
		// 当前是非 Windows 系统
		nowSaveRelativePath = strings.ReplaceAll(s.SaveRelativePath, "\\", "/")
	}

	orgRDirPath := filepath.Dir(nowSaveRelativePath)
	nowSubtitleExt := filepath.Ext(nowSaveRelativePath)
	return strings.ReplaceAll(filepath.Join(orgRDirPath, s.SubSha256+nowSubtitleExt), "\\", "/")
}
