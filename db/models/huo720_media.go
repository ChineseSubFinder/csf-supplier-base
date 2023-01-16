package models

import (
	"crypto/sha256"
	"fmt"
	"github.com/WQGroup/logger"
	"gorm.io/gorm"
	"net/url"
)

type Huo720Media struct {
	gorm.Model
	UrlPathUID string `gorm:"column:url_path_uid;primary_key;type:varchar(64);not null" json:"url_path_uid"` // 详情的 URL 的 PATH 计算 SHA256，/subject/3042261/
	IMDBid     string `gorm:"column:imdb_id;primary_key;type:varchar(20);not null" json:"imdb_id"`           // IMDB ID
	UrlPath    string `gorm:"column:url_path;type:varchar(255);not null" json:"url_path"`                    // 详情的 URL 的 PATH，/subject/3042261/
	UrlScheme  string `gorm:"column:url_scheme;type:varchar(255);not null" json:"url_scheme"`                // 详情的 URL 的 scheme，https 或者是 http
	UrlHost    string `gorm:"column:url_host;type:varchar(255);not null" json:"url_host"`                    // 详情的 URL 的 HOST，movie.douban.com
	IsMovie    bool   `gorm:"column:is_movie;type:bool;not null;default:false" json:"is_movie"`
}

func NewHuo720Media(urlStr string, imdbId string, isMove bool) *Huo720Media {

	baseUrlInfo, err := url.Parse(urlStr)
	if err != nil {
		logger.Panicln("NewHuo720Media urlStr parse error", err)
	}
	UrlPathUID := fmt.Sprintf("%x", sha256.Sum256([]byte(baseUrlInfo.Path)))

	return &Huo720Media{
		UrlPathUID: UrlPathUID,
		IMDBid:     imdbId,
		UrlPath:    baseUrlInfo.Path,
		UrlScheme:  baseUrlInfo.Scheme,
		UrlHost:    baseUrlInfo.Host,
		IsMovie:    isMove,
	}
}

func (h *Huo720Media) Url() string {
	return fmt.Sprintf("%s://%s%s", h.UrlScheme, h.UrlHost, h.UrlPath)
}
