package media_info_dealer

import "github.com/ChineseSubFinder/csf-supplier-base/db/models"

type MediaInfo struct {
	ImdbId           string
	TmdbId           string
	TVdbId           string
	OriginalTitle    string
	OriginalLanguage string  // 视频的原始语言  en zh
	TitleEn          string  // 英文标题
	TitleCn          string  // 中文的标题
	Year             string  // 播出的时间，如果是连续剧是第一次播出的时间 2019-01-01  2022-01-01
	Vote             float32 // 评分
	VoteCount        int     // 评分人数
}

func (m MediaInfo) UpdateInfo(mixMediaInfo *models.MixMediaInfo) {

	mixMediaInfo.NameCn = m.TitleCn
	mixMediaInfo.NameOrg = m.OriginalTitle
	mixMediaInfo.NameEn = m.TitleEn
	mixMediaInfo.TMDBid = m.TmdbId
	mixMediaInfo.TVDBid = m.TVdbId
	mixMediaInfo.ReleaseTime = m.Year
}
