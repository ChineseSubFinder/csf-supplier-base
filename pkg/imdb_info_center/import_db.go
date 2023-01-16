package imdb_info_center

import (
	"bytes"
	"fmt"
	"github.com/ChineseSubFinder/csf-supplier/pkg"
	"github.com/ChineseSubFinder/csf-supplier/pkg/imdb_info_center/dao"
	models2 "github.com/ChineseSubFinder/csf-supplier/pkg/imdb_info_center/models"
	"github.com/ChineseSubFinder/csf-supplier/pkg/settings"
	"github.com/WQGroup/logger"
	"github.com/valyala/tsvreader"
	"gorm.io/gorm"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func ImportAll(tmpDir string) {

	var wg sync.WaitGroup
	for _, filename := range settings.Get().ImdbInfoCenterConfig.Files {

		// 下载的压缩包
		fileFPath := filepath.Join(tmpDir, filename)
		// 解压出来的文件
		decompressFilePath := strings.ReplaceAll(fileFPath, filepath.Ext(fileFPath), "")
		if pkg.IsFile(decompressFilePath) == true {
			// 需要导入
			wg.Add(1)
			importTSVData(&wg, decompressFilePath)
		} else {
			// 没有找到导入的文件
			logger.Warningln("Not found decompress file : ", decompressFilePath)
		}
	}
	wg.Wait()
	logger.Infoln("ImportAll Done")
}

func importTSVData(wg *sync.WaitGroup, tsvFPath string) {

	defer func() {
		wg.Done()
	}()

	f, err := os.Open(tsvFPath)
	if err != nil {
		logger.Panicln(err)
	}

	tsvName := tsvFPath2TSVName(tsvFPath)

	allLineCount, err := lineCounter(f)
	if err != nil {
		logger.Panicln(err)
	}
	_, err = f.Seek(0, 0)
	if err != nil {
		logger.Panicf("Encountered error while seek 0 : %v", err)
	}

	r := tsvreader.New(f)
	counter := 0

	err = dao.Get().Transaction(func(tx *gorm.DB) error {

		for r.Next() {

			switch tsvName {
			case TitleBasics:
				{
					titleBasic := models2.TitleBasic{}
					titleBasic.TConst = r.String()
					titleBasic.TitleType = r.String()
					titleBasic.PrimaryTitle = r.String()
					titleBasic.OriginalTitle = r.String()
					adult := r.String()
					titleBasic.StartYear = r.String()
					titleBasic.EndYear = r.String()
					titleBasic.RuntimeMinutes = r.String()
					titleBasic.Genres = r.String()
					// 先要读取，然后再判断是否是 Title 第一行，跳过
					if counter > 0 {
						if isWantedTitleType(titleBasic.TitleType) == false {
							fmt.Fprintln(os.Stderr, counter, "/", allLineCount, "Skip TitleType : ", titleBasic.TitleType)
							counter++
							continue
						}
						counter++
						iAdult, _ := pkg.GetNumber2int(adult)
						if iAdult == 1 {
							titleBasic.IsAdult = true
						} else {
							titleBasic.IsAdult = false
						}
						// 这个 TConst 是否存在，存在则跳过插入
						fmt.Printf("%d / %d - TitleBasics.TConst : %s\n", counter, allLineCount, titleBasic.TConst)
						var tbs []models2.TitleBasic
						tx.Where("tconst = ?", titleBasic.TConst).Find(&tbs)
						if tbs != nil && len(tbs) > 0 {
							// 存在则跳过
							continue
						}
						tx.Save(&titleBasic)
					} else {
						counter++
					}
				}
			case TitleEpisode:
				{
					titleEpisode := models2.TitleEpisode{}
					titleEpisode.TConst = r.String()
					titleEpisode.ParentTConst = r.String()
					titleEpisode.SeasonNumber, _ = pkg.GetNumber2int(r.String())
					titleEpisode.EpisodeNumber, _ = pkg.GetNumber2int(r.String())
					// 先要读取，然后再判断是否是 Title 第一行，跳过
					if counter > 0 {
						counter++
						// 这个 TConst 是否存在，存在则跳过插入
						// 这个 TConst 是否存在，存在则跳过插入
						fmt.Printf("%d / %d - TitleEpisode.TConst : %s\n", counter, allLineCount, titleEpisode.TConst)
						var tbs []models2.TitleEpisode
						tx.Where("tconst = ?", titleEpisode.TConst).Find(&tbs)
						if tbs != nil && len(tbs) > 0 {
							// 存在则跳过
							continue
						}
						tx.Save(&titleEpisode)
					} else {
						counter++
					}
				}
			case TitleRatings:
				{
					titleRating := models2.TitleRatings{}
					titleRating.TConst = r.String()
					titleRating.AverageRating, _ = pkg.GetNumber2Float(r.String())
					titleRating.NumVotes, _ = pkg.GetNumber2int(r.String())
					// 先要读取，然后再判断是否是 Title 第一行，跳过
					if counter > 0 {
						counter++
						// 这个 TConst 是否存在，存在则跳过插入
						// 这个 TConst 是否存在，存在则跳过插入
						fmt.Printf("%d / %d - TitleRatings.TConst : %s\n", counter, allLineCount, titleRating.TConst)
						var tbs []models2.TitleRatings
						tx.Where("tconst = ?", titleRating.TConst).Find(&tbs)
						if tbs != nil && len(tbs) > 0 {
							// 存在则跳过
							continue
						}
						tx.Save(&titleRating)
					} else {
						counter++
					}
				}
			default:
				logger.Panicln("Unknown tsvName : ", tsvName)
			}
		}
		// 提交事务
		return nil
	})
	if err != nil {
		logger.Panicln(err)
	}

	logger.Infoln("Import", tsvName.DeCompressNameString(), "Done.")
}

func lineCounter(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

func tsvFPath2TSVName(tsvFPath string) TSVName {

	tsvNameStr := filepath.Base(tsvFPath)
	tsvName := TitleAkas
	switch tsvNameStr {
	case TitleAkas.DeCompressNameString():
		tsvName = TitleAkas
	case TitleBasics.DeCompressNameString():
		tsvName = TitleBasics
	case TitleCrew.DeCompressNameString():
		tsvName = TitleCrew
	case TitleEpisode.DeCompressNameString():
		tsvName = TitleEpisode
	case TitlePrincipals.DeCompressNameString():
		tsvName = TitlePrincipals
	case TitleRatings.DeCompressNameString():
		tsvName = TitleRatings
	case NameBasics.DeCompressNameString():
		tsvName = NameBasics
	default:
		logger.Panicln("Unknown tsvNameStr : ", tsvNameStr)
	}

	return tsvName
}

// isWantedTitleType 只想要 电影和电视剧 相关的，其他的都不要
func isWantedTitleType(inputTitleType string) bool {

	lowerStr := strings.ToLower(inputTitleType)
	if strings.Contains(lowerStr, "movie") == true {
		return true
	}
	if strings.Contains(lowerStr, "tv") == true {
		return true
	}
	if strings.Contains(lowerStr, "episode") == true {
		return true
	}
	return false
}

type TSVName int

const (
	TitleAkas TSVName = iota + 1
	TitleBasics
	TitleCrew
	TitleEpisode
	TitlePrincipals
	TitleRatings
	NameBasics
)

func (s TSVName) GZNameString() string {
	switch s {
	case TitleAkas:
		return "title.akas.tsv.gz"
	case TitleBasics:
		return "title.basics.tsv.gz"
	case TitleCrew:
		return "title.crew.tsv.gz"
	case TitleEpisode:
		return "title.episode.tsv.gz"
	case TitlePrincipals:
		return "title.principals.tsv.gz"
	case TitleRatings:
		return "title.ratings.tsv.gz"
	case NameBasics:
		return "name.basics.tsv.gz"
	default:
		return "Unknown"
	}
}

func (s TSVName) DeCompressNameString() string {
	return strings.ReplaceAll(s.GZNameString(), filepath.Ext(s.GZNameString()), "")
}
