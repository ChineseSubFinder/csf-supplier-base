package models

import (
	"crypto/sha256"
	"fmt"
	"github.com/ChineseSubFinder/csf-supplier-base/pkg/settings"
	"github.com/WQGroup/logger"
	"net/url"
	"path"
)

type ZiMuKuInfo struct {
	IMDBId     string `gorm:"column:imdb_id;unique;type:varchar(20);not null" json:"imdb_id"`                // IMDB ID
	UrlPathUID string `gorm:"column:url_path_uid;primary_key;type:varchar(64);not null" json:"url_path_uid"` // 具体一个字幕 URL 的 PATH 计算 SHA256，/detail/177838.html
	UrlPath    string `gorm:"column:url_path;type:varchar(255);not null" json:"url_path"`                    // 具体一个字幕 URL 的 PATH，/detail/177838.html
}

func NewZiMuKuInfo(imdbId, fullUrl string) *ZiMuKuInfo {

	baseUrlInfo, err := url.Parse(fullUrl)
	if err != nil {
		logger.Panicln("NewZiMuKuInfo urlStr parse error", err)
	}
	UrlPathUID := fmt.Sprintf("%x", sha256.Sum256([]byte(baseUrlInfo.Path)))

	z := &ZiMuKuInfo{
		IMDBId:     imdbId,
		UrlPathUID: UrlPathUID,
		UrlPath:    baseUrlInfo.Path,
	}

	return z
}

func (z ZiMuKuInfo) GetFullUrl() string {

	baseUrl := settings.Get().ZiMuKuConfig.SiteRootUrl
	// 拼接 url 信息
	return path.Join(baseUrl, z.UrlPath)
}
