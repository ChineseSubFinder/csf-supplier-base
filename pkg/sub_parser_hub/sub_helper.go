package sub_parser_hub

import (
	"github.com/ChineseSubFinder/csf-supplier-base/pkg/archive_helper"
	"github.com/ChineseSubFinder/csf-supplier-base/pkg/common"
	"github.com/ChineseSubFinder/csf-supplier-base/pkg/filter"
	"github.com/WQGroup/logger"
	"os"
	"path/filepath"
	"strings"
)

// ChangeVideoExt2SubExt 检测 Name，如果是视频的后缀名就改为字幕的后缀名
func ChangeVideoExt2SubExt(subInfos []*SubInfo) {
	for x, info := range subInfos {
		tmpSubFileName := info.Name
		// 如果后缀名是下载字幕目标的后缀名  或者 是压缩包格式的，则跳过
		if strings.Contains(tmpSubFileName, info.Ext) == true || archive_helper.IsWantedArchiveExtName(tmpSubFileName) == true {

		} else {
			subInfos[x].Name = tmpSubFileName + info.Ext
		}
	}
}

// SearchMatchedSubFileByDir 搜索符合后缀名的视频文件，排除 Sub_SxE0 这样的文件夹中的文件
func SearchMatchedSubFileByDir(inDir string) (*SearchSubResult, error) {

	result := NewSearchSubResult()
	err := filepath.Walk(inDir, func(path string, info os.FileInfo, err error) error {

		if info.IsDir() == true {
			return nil
		}
		// 这里就是文件了
		if filter.SkipFileInfo(info, filepath.Base(filepath.Dir(path))) == true {
			return nil
		}

		switch IsSubExtWanted(filepath.Ext(info.Name())) {
		case common.Characters:
			result.Add(common.Characters, path)
		case common.Picture:
			result.Add(common.Picture, path)
		case common.BluRay:
			result.Add(common.BluRay, path)
		}

		return nil
	})

	return result, err
}

// SearchMatchedSubFileByOneVideo 搜索这个视频当前目录下匹配的字幕
func SearchMatchedSubFileByOneVideo(oneVideoFullPath string) (*SearchSubResult, error) {

	dir := filepath.Dir(oneVideoFullPath)
	fileName := filepath.Base(oneVideoFullPath)
	fileName = strings.ToLower(fileName)
	fileName = strings.ReplaceAll(fileName, filepath.Ext(fileName), "")

	result := NewSearchSubResult()
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {

		if info.IsDir() == true {
			return nil
		}
		// 这里就是文件了
		if filter.SkipFileInfo(info, filepath.Base(filepath.Dir(path))) == true {
			return nil
		}
		// 判断的时候用小写的，后续重命名的时候用原有的名称
		nowFileName := strings.ToLower(info.Name())

		// 字幕文件名应该包含 视频文件名（无后缀）
		if strings.HasPrefix(nowFileName, fileName) == false {
			return nil
		}
		switch IsSubExtWanted(filepath.Ext(info.Name())) {
		case common.Characters:
			result.Add(common.Characters, path)
		case common.Picture:
			result.Add(common.Picture, path)
		case common.BluRay:
			result.Add(common.BluRay, path)
		}

		return nil
	})

	return result, err
}

type SearchSubResult struct {
	charactersSubtitles []string
	picturesSubtitles   []string
	blueRaySubtitles    []string
}

func NewSearchSubResult() *SearchSubResult {

	return &SearchSubResult{
		charactersSubtitles: make([]string, 0),
		picturesSubtitles:   make([]string, 0),
		blueRaySubtitles:    make([]string, 0),
	}
}

func (s *SearchSubResult) Add(iType common.SubtitleType, inSubFPath string) {
	switch iType {
	case common.Characters:
		s.charactersSubtitles = append(s.charactersSubtitles, inSubFPath)
	case common.Picture:
		s.picturesSubtitles = append(s.picturesSubtitles, inSubFPath)
	case common.BluRay:
		s.blueRaySubtitles = append(s.blueRaySubtitles, inSubFPath)
	default:
		logger.Panicln("Add SubtitleType Error")
	}
}

func (s *SearchSubResult) Get(iType common.SubtitleType) []string {
	switch iType {
	case common.Characters:
		return s.charactersSubtitles
	case common.Picture:
		return s.picturesSubtitles
	case common.BluRay:
		return s.blueRaySubtitles
	default:
		logger.Panicln("Get SubtitleType Error")
	}
	return nil
}
