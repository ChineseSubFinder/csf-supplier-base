package models

import (
	"fmt"
	"github.com/ChineseSubFinder/csf-supplier-base/pkg"
	"github.com/ChineseSubFinder/csf-supplier-base/pkg/types/language"
	"github.com/pkg/errors"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type SubtitleInfo struct {
	SubSha256            string  `gorm:"column:sub_sha256;type:char(64)"`                                                        // 文件的 SHA256
	ImdbId               string  `gorm:"column:imdb_id;type:char(20);index;not null"`                                            // IMDB ID，电影的 IMDB ID，电视剧的主 IMDB ID
	Title                string  `gorm:"column:title;type:varchar(255);not null" json:"title"`                                   // 字幕标题
	Language             int     `gorm:"column:language;type:int;not null"`                                                      // 字幕语言,MyLanguage 参考内部的这个值
	SaveRelativePath     string  `gorm:"column:save_relative_path;type:varchar(255);not null" json:"save_relative_path"`         // 保存的相对路径，包含文件名，会处理后一定是一个具体的字幕文件 /movie/2020/12/12/177838.srt
	Score                float32 `gorm:"column:score;type:float;not null;default:0" json:"score"`                                // 评分
	Votes                int     `gorm:"column:votes;type:int;not null;default:0" json:"votes"`                                  // 投票数
	DownloadTimes        int     `gorm:"column:download_times;type:int;not null;default:0" json:"download_times"`                // 下载次数
	UploadTime           int64   `gorm:"column:upload_time;type:bigint;not null;default:0" json:"upload_time"`                   // 上传时间
	SubtitlesComesFrom   string  `gorm:"column:subtitles_comes_from;type:varchar(255);not null" json:"subtitles_comes_from"`     // 字幕来源，这里不是字幕站，而是制作的人和组
	Upload2R2            bool    `gorm:"column:upload2r2;type:tinyint(1);not null;default:0" json:"upload2r2"`                   // 是否上传到 r2
	Upload2R2Time        int64   `gorm:"column:upload2r2_time;type:bigint;not null;default:0" json:"upload2r2_time"`             // 上传到 r2 的时间
	Upload2R2Result      string  `gorm:"column:upload2r2_result;type:varchar(255);not null" json:"upload2r2_result"`             // 上传到 r2 的结果
	Upload2CloudDb       bool    `gorm:"column:upload2cloud_db;type:tinyint(1);not null;default:0" json:"upload2cloud_db"`       // 是否上传到云数据库，这样就可以在云端搜索了
	Upload2CloudDbTime   int64   `gorm:"column:upload2cloud_db_time;type:bigint;not null;default:0" json:"upload2cloud_db_time"` // 上传到云数据库的时间
	Upload2CloudDbResult string  `gorm:"column:upload2cloud_db_result;type:varchar(255);not null" json:"upload2cloud_db_result"` // 上传到云数据库的结果
}

func (s *SubtitleInfo) Ext() string {
	return filepath.Ext(s.SaveRelativePath)
}

func (s *SubtitleInfo) Lang() language.MyLanguage {
	return language.MyLanguage(s.Language)
}

func (s *SubtitleInfo) R2StoreKey() string {

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

func (s *SubtitleInfo) GetSubtitleData(saveRootDirPath string) ([]byte, error) {
	// 从存储的路径进行相对路径的拼接，得到完整的路径
	subtitleFullPath := filepath.Join(saveRootDirPath, s.SaveRelativePath)
	// 判断这个文件是否存在，如果存在则读取文件内容
	if pkg.IsFile(subtitleFullPath) == false {
		return nil, errors.New(fmt.Sprintf("subtitle file not exist, path: %s", subtitleFullPath))
	}
	return ioutil.ReadFile(subtitleFullPath)
}

func (s *SubtitleInfo) MarkUpload2R2() {
	s.Upload2R2 = true
	s.Upload2R2Time = time.Now().Unix()
	s.Upload2R2Result = ""
}

func (s *SubtitleInfo) MarkUpload2R2Failed(failedStr string) {
	s.Upload2R2 = false
	s.Upload2R2Time = time.Now().Unix()
	s.Upload2R2Result = failedStr[:250] // 避免越界
}

func (s *SubtitleInfo) MarkUpload2CloudDb() {
	s.Upload2CloudDb = true
	s.Upload2CloudDbTime = time.Now().Unix()
	s.Upload2CloudDbResult = ""
}

func (s *SubtitleInfo) MarkUpload2CloudDbFailed(failedStr string) {
	s.Upload2CloudDb = false
	s.Upload2CloudDbTime = time.Now().Unix()
	s.Upload2CloudDbResult = failedStr[:250] // 避免越界
}

func (s *SubtitleInfo) GetTitle(needFilterKeyWords []string) string {

	nowTitle := s.Title
	// 因为下载的文件名中，可能包含一些 [zmk.pw] 这样的信息，需要剔除
	for _, keyWord := range needFilterKeyWords {
		nowTitle = strings.ReplaceAll(nowTitle, keyWord, "")
	}
	return nowTitle
}
