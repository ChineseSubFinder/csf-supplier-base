package search

import (
	"github.com/ChineseSubFinder/csf-supplier-base/db/models"
	"github.com/WQGroup/logger"
	"github.com/pkg/errors"
	"strconv"
)

// VideoNameSearchKeywordMaker 拼接视频搜索的 title 和 年份，考虑到 2020 年开始搜索带有 Year 效率才高
func VideoNameSearchKeywordMaker(title string, year string) string {
	iYear, err := strconv.Atoi(year)
	if err != nil {
		// 允许的错误
		logger.Errorln("VideoNameSearchKeywordMaker", "year to int", err)
		iYear = 0
	}
	searchKeyword := title
	if iYear >= 2020 {
		searchKeyword = searchKeyword + " " + year
	}

	return searchKeyword
}

// KeyWordSelect keyWordType cn, 中文， en，英文，org，原始名称
func KeyWordSelect(mediaInfo *models.MixMediaInfo, l Language) (string, error) {

	keyWord := ""
	switch l {
	case CN:
		keyWord = mediaInfo.NameCn
		if keyWord == "" {
			return "", errors.New("TitleCn is empty")
		}
	case EN:
		keyWord = mediaInfo.NameEn
		if keyWord == "" {
			return "", errors.New("TitleEn is empty")
		}
	case Org:
		keyWord = mediaInfo.NameOrg
		if keyWord == "" {
			return "", errors.New("OriginalTitle is empty")
		}
	default:
		return "", errors.New("keyWordType is not cn, en, org")
	}

	return keyWord, nil
}

type Language int

const (
	CN Language = iota + 1
	EN
	Org
)

func (l Language) String() string {
	switch l {
	case CN:
		return "cn"
	case EN:
		return "en"
	case Org:
		return "org"
	default:
		return "Unknown"
	}
}
