package download_center

import (
	"github.com/ChineseSubFinder/csf-supplier-base/db/dao"
	"github.com/ChineseSubFinder/csf-supplier-base/db/models"
	"github.com/WQGroup/logger"
)

/*
	下载中心的目标是把一个电影或者电视剧的所有资源下载到本地
	最多是把连续剧是季进行区分存储
	后续需要配合，整理的模块来处理这些下载好的资源
*/
type Center struct {
}

func NewCenter() *Center {
	return &Center{}
}

// MarkDownloaded 标记这个 url 的文件已经下载过了
func (c Center) MarkDownloaded(site Site, downloadedInfo models.DownloadedInfo) {

	if site == ZiMuKu {

		if c.IsDownloaded(site, downloadedInfo.FullUrlPath) == true {
			// 已经插入了，那么就跳过
			logger.Debugln("Center MarkDownloaded site, Skip:", site, "url:", downloadedInfo.FullUrlPath)
			return
		}
		tmpInfo := models.DownloadedInfoZiMuKu{
			DownloadedInfo: downloadedInfo,
		}
		dao.Get().Create(&tmpInfo)
	} else {
		logger.Panicln("Center MarkDownloaded site not support:", site)
	}
}

// IsDownloaded url 的文件是否已经下载过了
func (c Center) IsDownloaded(site Site, fullUrlPath string) bool {

	if site == ZiMuKu {
		zimuku := models.NewDownloadedInfoZiMuKu(fullUrlPath)
		var downloadedInfo []models.DownloadedInfoZiMuKu
		dao.Get().Where("url_path_uid = ?", zimuku.UrlPathUID).Find(&downloadedInfo)
		if len(downloadedInfo) > 0 {
			return true
		}
	} else {
		logger.Panicln("Center IsDownloaded site not support:", site)
	}

	return false
}

type Site int

const (
	ZiMuKu Site = iota + 1
	SubHD
)

func (s Site) String() string {
	switch s {
	case ZiMuKu:
		return "ZiMuKu"
	case SubHD:
		return "SubHD"
	default:
		return "Unknown"
	}
}
