package manual_handling

import (
	"fmt"
	"github.com/ChineseSubFinder/csf-supplier/internal/dao"
	"github.com/ChineseSubFinder/csf-supplier/internal/models"
)

func Add(stage models.Stage, url, remark string) {

	// 查询是否已经插入这个 URL 信息
	var mhs []models.ManualHandling
	dao.Get().Where("url = ?", url).Find(&mhs)
	if len(mhs) > 0 {
		// 已经存在
		return
	}
	mh := models.NewManualHandling(stage, url, remark, false)
	dao.Get().Save(&mh)
}

func AddImdbIdNotInDb(stage models.Stage, imdbId string, isMovie bool) {

	mediaType := "电影"
	if isMovie == false {
		mediaType = "电视剧"
	}
	needInsertRemark := fmt.Sprintf("%s %s 不在数据库中", mediaType, imdbId)
	// 查询是否已经插入这个 remark 信息
	var mhs []models.ManualHandling
	dao.Get().Where("remark = ?", needInsertRemark).Find(&mhs)
	if len(mhs) > 0 {
		// 已经存在
		return
	}
	mh := models.NewManualHandling(stage, "", needInsertRemark, false)
	dao.Get().Save(&mh)
}
