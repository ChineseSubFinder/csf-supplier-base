package models

import (
	"github.com/ChineseSubFinder/csf-supplier-base/pkg"
	"github.com/WQGroup/logger"
	"path/filepath"
	"strings"
)

type DownloadedInfo struct {
	FullUrlPath        string  `gorm:"-"`                                                                                  // 全路径 https://zimuku.org/detail/177838.html
	UrlPathUID         string  `gorm:"column:url_path_uid;primary_key;type:varchar(64);not null" json:"url_path_uid"`      // 具体一个字幕 URL 的 PATH 计算 SHA256，/detail/177838.html
	UrlPath            string  `gorm:"column:url_path;type:varchar(255);not null" json:"url_path"`                         // 具体一个字幕 URL 的 PATH，/detail/177838.html
	SaveRelativePath   string  `gorm:"column:save_relative_path;type:varchar(255);not null" json:"save_relative_path"`     // 保存的相对路径，包含文件名，可能是一个压缩包，也可能是一个字幕文件 /movie/2020/12/12/177838.srt
	Score              float32 `gorm:"column:score;type:float;not null;default:0" json:"score"`                            // 评分
	Votes              int     `gorm:"column:votes;type:int;not null;default:0" json:"votes"`                              // 投票数
	DownloadTimes      int     `gorm:"column:download_times;type:int;not null;default:0" json:"download_times"`            // 下载次数
	UploadTime         int64   `gorm:"column:upload_time;type:bigint;not null;default:0" json:"upload_time"`               // 上传时间
	Title              string  `gorm:"column:title;type:varchar(255);not null" json:"title"`                               // 字幕标题
	SubtitlesComesFrom string  `gorm:"column:subtitles_comes_from;type:varchar(255);not null" json:"subtitles_comes_from"` // 字幕来源，这里不是字幕站，而是制作的人和组
}

func (d DownloadedInfo) IsMovie() bool {
	return strings.HasPrefix(d.SaveRelativePath, "movie")
}

func (d DownloadedInfo) Info() DSubInfo {

	var downloadSubInfo DSubInfo
	if d.IsMovie() == true {
		// 电影
		downloadSubInfo.IsMovie = true
		imdbIdFolderName := filepath.Base(filepath.Dir(d.SaveRelativePath))
		if strings.HasPrefix(imdbIdFolderName, "tt") == false {
			// 需要停下来找问题
			logger.Panicln("imdbIdFolderName is not start with tt", d.SaveRelativePath)
		}
		downloadSubInfo.ImdbId = imdbIdFolderName
		downloadSubInfo.SubFileName = filepath.Base(d.SaveRelativePath)
		downloadSubInfo.Season = -1
	} else {
		// 电视剧
		imdbIdFolderName := ""
		downloadSubInfo.IsMovie = false
		// 因为可能可以解析出 Season 的概念，也可能解析不出来 Season 这个信息
		// 那么就需要在这里进行判断，到底有不有 Season 的这个文件夹存在
		// 1. 正常是 ttxxxx/1/subName.srt 这样的格式
		// 2. 但也会有 ttxxxx/subName.srt 这样的格式

		// subSeasonPathDir  --> tv\tt3032476\4   或者 tv\tt3032476
		subSeasonPathDir := filepath.Dir(d.SaveRelativePath)
		subGusSeasonName := filepath.Base(subSeasonPathDir)
		if strings.HasPrefix(subGusSeasonName, "tt") == true {
			// 那么就是情况2，没有 Season 的文件夹
			imdbIdFolderName = filepath.Base(subSeasonPathDir)
			downloadSubInfo.Season = -1
		} else {
			// 再向上一级
			var subGusImdbIdName = filepath.Base(filepath.Dir(subSeasonPathDir))
			if strings.HasPrefix(subGusImdbIdName, "tt") == true {
				imdbIdFolderName = subGusImdbIdName
			} else {
				logger.Panicln("imdbIdFolderName is not start with tt", d.SaveRelativePath)
			}
			// 那么理论上就找到 Season 文件夹的信息了
			seasonFolderName := subGusSeasonName
			number2int, err := pkg.GetNumber2int(seasonFolderName)
			if err != nil {
				logger.Panicln("seasonFolderName is not number", d.SaveRelativePath)
			}
			downloadSubInfo.Season = number2int
		}

		downloadSubInfo.ImdbId = imdbIdFolderName
		downloadSubInfo.SubFileName = filepath.Base(d.SaveRelativePath)
	}

	return downloadSubInfo
}

func (d DownloadedInfo) SubFileFPath(orgSaveRootDirPath string) string {
	return filepath.Join(orgSaveRootDirPath, d.SaveRelativePath)
}

type DSubInfo struct {
	ImdbId      string
	SubFileName string
	IsMovie     bool
	Season      int
}
