package models

import (
	"crypto/sha256"
	"fmt"
	"github.com/WQGroup/logger"
	"gorm.io/gorm"
	"net/url"
)

type DownloadedInfoZiMuKu struct {
	gorm.Model
	DownloadedInfo `gorm:"embedded"`
}

func NewDownloadedInfoZiMuKu(urlFullPathStr string) *DownloadedInfoZiMuKu {

	baseUrlInfo, err := url.Parse(urlFullPathStr)
	if err != nil {
		logger.Panicln("NewDownloadedInfoZiMuKu urlStr parse error", err)
	}
	UrlPathUID := fmt.Sprintf("%x", sha256.Sum256([]byte(baseUrlInfo.Path)))

	return &DownloadedInfoZiMuKu{
		DownloadedInfo: DownloadedInfo{
			FullUrlPath: urlFullPathStr,
			UrlPathUID:  UrlPathUID,
			UrlPath:     baseUrlInfo.Path,
		},
	}
}
