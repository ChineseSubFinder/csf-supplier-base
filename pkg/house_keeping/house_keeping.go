package house_keeping

import (
	"fmt"
	"github.com/ChineseSubFinder/csf-supplier-base/db/dao"
	"github.com/ChineseSubFinder/csf-supplier-base/db/models"
	"github.com/ChineseSubFinder/csf-supplier-base/pkg"
	"github.com/ChineseSubFinder/csf-supplier-base/pkg/archive_helper"
	"github.com/ChineseSubFinder/csf-supplier-base/pkg/common"
	"github.com/ChineseSubFinder/csf-supplier-base/pkg/settings"
	"github.com/ChineseSubFinder/csf-supplier-base/pkg/sub_parser_hub"
	subparser "github.com/ChineseSubFinder/csf-supplier-base/pkg/sub_parser_hub/sub_parser"
	"github.com/ChineseSubFinder/csf-supplier-base/pkg/sub_parser_hub/sub_parser/ass"
	"github.com/ChineseSubFinder/csf-supplier-base/pkg/sub_parser_hub/sub_parser/srt"
	"github.com/WQGroup/logger"
	"github.com/allanpk716/rod_helper"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type HouseKeeping struct {
	subDownloadedSaveRootPath string                       // 字幕下载后保存的根目录
	subParserHub              *sub_parser_hub.SubParserHub // 字幕解析器
}

func NewHouseKeeping(subDownloadedSaveRootPath string) *HouseKeeping {

	var houseKeeping HouseKeeping
	houseKeeping.subDownloadedSaveRootPath = subDownloadedSaveRootPath
	houseKeeping.subParserHub = sub_parser_hub.NewSubParserHub(
		logger.GetLogger(),
		ass.NewParser(logger.GetLogger()),
		srt.NewParser(logger.GetLogger()))

	return &houseKeeping
}

func (h HouseKeeping) GetSourceIds(sourceId SourceId) []models.DownloadedInfoZiMuKu {

	var downloadInfoZiMuKus []models.DownloadedInfoZiMuKu

	switch sourceId {
	case DownloadedInfoZiMuKu:
		/*
			首先需要从数据库（download_info_zimuku）中，获取现在有下载哪些字幕
			按主键 ID 去分批获取，然后一个个处理，处理完毕后，需要存储到以下几个表中：
			1. HouseKeeping		这个是记录处理到第几个 ID 了
			2. subtitle_movie	如果是电影的字幕处理完毕就写入到这里
			3. subtitle_tv		如果是电视剧的字幕处理完毕就写入到这里
		*/
		// 查询最后一个处理的 ID
		var houseKeepings []models.HouseKeeping
		dao.Get().Order("downloaded_id desc").Where("which_site = ?", models.ZiMuKu.Index()).Limit(1).Find(&houseKeepings)
		var processStartId uint
		if len(houseKeepings) < 1 {
			// 没有记录，从 1 开始，全新的开始
			processStartId = 0
		} else {
			// 需要查询 Error 列表中，最后一个处理的 ID
			// 查询最后一个处理的 ID
			var houseKeepingErrors []models.HouseKeepingError
			dao.Get().Order("downloaded_id desc").Where("which_site = ?", models.ZiMuKu.Index()).Limit(1).Find(&houseKeepingErrors)

			if len(houseKeepingErrors) < 1 {
				processStartId = houseKeepings[0].DownloadedSubId
			} else {
				if houseKeepings[0].DownloadedSubId < houseKeepingErrors[0].DownloadedSubId {
					processStartId = houseKeepingErrors[0].DownloadedSubId
				} else {
					processStartId = houseKeepings[0].DownloadedSubId
				}
			}
		}
		// processStartId = 523
		// 然后从这个 ID 开始，每次查询 10 个进行处理
		dao.Get().Where("id > ? AND id <= ?", processStartId, processStartId+10).Find(&downloadInfoZiMuKus)

	case HouseKeepingError:
		/*
			说明：这个是处理 HouseKeeping 表中的错误数据
			如果处理成功了，那么就需要删除这个数据
		*/
		// 查询最后一个处理的 ID
		var houseKeepingErrors []models.HouseKeepingError
		dao.Get().Order("downloaded_id desc").Where("unzip_error = ? AND which_site = ?", true, models.ZiMuKu.Index()).Limit(1).Find(&houseKeepingErrors)
		var processStartId uint
		if len(houseKeepingErrors) < 1 {
			// 没有记录，从 1 开始，全新的开始
			logger.Infoln("HouseKeepingError 没有记录，退出")
			return nil
		} else {
			processStartId = houseKeepingErrors[0].DownloadedSubId
		}

		dao.Get().Where("id = ?", processStartId).Limit(1).Find(&downloadInfoZiMuKus)
	default:
		logger.Panicln("Not Supported SourceId", sourceId)
	}

	return downloadInfoZiMuKus
}

func (h HouseKeeping) Process(site models.WhichSite, sourceId SourceId, downloadInfoZiMuKus []models.DownloadedInfoZiMuKu) {

	subsSaveRootDirPath := ""

	var err error
	var houseKeepings []models.HouseKeeping

	switch site {
	case models.ZiMuKu:
		subsSaveRootDirPath = settings.Get().ZiMuKuConfig.SubsSaveRootDirPath
	default:
		logger.Panicln("Not Supported Site", site)
	}

	for index, downloadedInfoZiMuKu := range downloadInfoZiMuKus {

		logger.Infoln("--------------------------------------------")
		logger.Infoln("Start Process", index+1, "/", len(downloadInfoZiMuKus), "ID:", downloadedInfoZiMuKu.ID, downloadedInfoZiMuKu.Title)
		// 判断这个是否已经处理过了，确保一定是没有处理过的
		if sourceId != HouseKeepingError {
			dao.Get().Where("downloaded_id = ? AND which_site = ?", downloadedInfoZiMuKu.ID, site.Index()).Find(&houseKeepings)
			if len(houseKeepings) > 0 {
				logger.Infoln("Already Processed")
				continue
			}
		}
		// 这个文件在本地要存在
		nowSubFileFPath := filepath.Join(h.subDownloadedSaveRootPath, downloadedInfoZiMuKu.SaveRelativePath)
		if pkg.IsFile(nowSubFileFPath) == false {
			// 如果这个文件不存在
			logger.Errorln("File Not Exist", nowSubFileFPath)
			// 暂时不按这个逻辑处理，因为可能配置填写错误，那么就导致忽略了一大堆本来存在而跳过不处理的
			continue
		}
		// 解析这个数据库的信息
		dSubInfo := downloadedInfoZiMuKu.Info()
		err = h.organizeSubs(site, subsSaveRootDirPath, downloadedInfoZiMuKu.ID, downloadedInfoZiMuKu.DownloadedInfo, dSubInfo)
		if err != nil {

			// 这个是否已经记录过了
			var houseKeepingError models.HouseKeepingError
			var houseKeepingErrors []models.HouseKeepingError
			dao.Get().Where("downloaded_id = ? AND which_site = ?", downloadedInfoZiMuKu.ID, site.Index()).Find(&houseKeepingErrors)
			if len(houseKeepingErrors) > 0 {
				// 已经存在了，更新
				houseKeepingError = houseKeepingErrors[0]
				houseKeepingError.ProcessTime = time.Now().Unix()
				houseKeepingError.SaveRelativePath = downloadedInfoZiMuKu.SaveRelativePath
			} else {
				houseKeepingError.DownloadedSubId = downloadedInfoZiMuKu.ID
				houseKeepingError.WhichSite = site.Index()
				houseKeepingError.ProcessTime = time.Now().Unix()
				houseKeepingError.SaveRelativePath = downloadedInfoZiMuKu.SaveRelativePath
			}

			if errors.Is(err, common.UnZipError) == true {
				// 如果是解压的错误，那么就更新错误的数据库，标记这个文件就算错误了，后续也需要跳过，不在找出来处理
				logger.Infoln("Skip", common.UnZipError, downloadedInfoZiMuKu.Title)
				// 存入错误的数据库中
				houseKeepingError.UnzipError = true
			} else if errors.Is(err, common.SubtitleExtTypeIsPicture) == true {

				logger.Infoln("Skip", common.SubtitleExtTypeIsPicture, downloadedInfoZiMuKu.Title)
				houseKeepingError.SubtitleExtType = (int)(common.Picture)

			} else if errors.Is(err, common.SubtitleExtTypeIsBluRay) == true {

				logger.Infoln("Skip", common.SubtitleExtTypeIsBluRay, downloadedInfoZiMuKu.Title)
				houseKeepingError.SubtitleExtType = (int)(common.BluRay)
			} else if errors.Is(err, common.SubtitleExtTypeNotSupported) == true {

				logger.Infoln("Skip", common.SubtitleExtTypeNotSupported, downloadedInfoZiMuKu.Title)
				houseKeepingError.SubtitleExtType = (int)(common.NotSupported)
			} else {
				logger.Infoln("Skip this one", downloadedInfoZiMuKu.Title, "Insert 2 HouseKeepingError DB")
			}
			// 存入错误的数据库中
			err = dao.Get().Save(&houseKeepingError).Error
			if err != nil {
				logger.Errorln("Insert 2 HouseKeepingError DB Error", downloadedInfoZiMuKu.Title, err)
			}
			continue
		}
		// 如果是处理上面错误的数据库中的数据，那么就可以删除了
		if sourceId == HouseKeepingError {
			logger.Infoln("Delete HouseKeepingError DB ID", downloadedInfoZiMuKu.ID)

			dao.Get().Where("downloaded_id = ?", downloadedInfoZiMuKu.ID).Delete(&models.HouseKeepingError{})
		}
	}
}

// organizeSubs 将下载的字幕进行一次组织，获取字幕的信息，然后存储
func (h HouseKeeping) organizeSubs(site models.WhichSite, subsSaveRootDirPath string, downloadDBId uint, downloadedInfo models.DownloadedInfo, dSubInfo models.DSubInfo) error {

	logger.Infoln("organizeSubs Start:", downloadedInfo.Title)
	var err error
	tmpDir := rod_helper.GetTmpFolderByName(settings.Get().HouseKeepingConfig.TmpRootDirPath, dSubInfo.ImdbId)
	defer func() {
		err = os.RemoveAll(tmpDir)
		if err != nil {
			logger.Errorln("RemoveAll", tmpDir, "Error", err)
		}
	}()
	// 注意这里的字幕网站的缓存路径也不一样
	nowFileFPath := downloadedInfo.SubFileFPath(subsSaveRootDirPath)
	nowExt := filepath.Ext(nowFileFPath)
	nowExt = strings.ToLower(nowExt)
	var fullSeason bool
	var sourceFileSha256 string
	var subFromZipFile bool // 这个字幕是通过 zip 解压出来的，因为有时候，解压出来的字幕是没得 Season 和 Eps 信息的，那么就需要从这个 zip 包的名称来解析
	subFileFullPaths := make([]string, 0)

	if nowExt != archive_helper.Zip.String() &&
		nowExt != archive_helper.Tar.String() &&
		nowExt != archive_helper.Rar.String() &&
		nowExt != archive_helper.SevenZ.String() {

		// 排除不是压缩包的情况，这里就是字幕或者是其他的文件
		// 是否是受支持的字幕类型
		subExtType := sub_parser_hub.IsSubExtWanted(nowExt)
		if subExtType == common.Characters {
			// 通过
		} else if subExtType == common.Picture {
			return common.SubtitleExtTypeIsPicture
		} else if subExtType == common.BluRay {
			return common.SubtitleExtTypeIsBluRay
		} else {
			return common.SubtitleExtTypeNotSupported
		}
		// 这里就是文字型的字幕
		//sourceFileSha256, err = pkg.GetFileSHA256String(nowFileFPath)
		//if err != nil {
		//	logger.Errorln("GetFileSHA256String", nowFileFPath, "Error", err)
		//	return err
		//}
		subFileFullPaths = append(subFileFullPaths, nowFileFPath)
	} else {

		subFromZipFile = true
		// 那么这里就应该都是压缩包的类型的文件了
		fullSeason, _, _, err = pkg.GetSeasonAndEpisodeFromSubFileName(filepath.Base(nowFileFPath))
		if err != nil {

			// 这里忽略解析不出来的情况，仅仅是判断是否是季度包就行了
			if errors.Is(err, common.GetSeasonAndEpisodeFromSubFileNameError) == false {
				logger.Errorln("GetSeasonAndEpisodeFromSubFileName", nowFileFPath, "Error", err)
				return err
			}
		}
		sourceFileSha256, err = pkg.GetFileSHA256String(nowFileFPath)
		if err != nil {
			logger.Errorln("GetFileSHA256String", nowFileFPath, "Error", err)
			return err
		}

		// 这里就是受支持的压缩文件
		err = archive_helper.UnArchiveFileEx(nowFileFPath, tmpDir, true, true)
		// 解压完成后，遍历受支持的字幕列表，加入缓存列表
		if err != nil {
			logger.Errorln("UnArchiveFileEx", nowFileFPath, err)
			return common.UnZipError
		}
		// 搜索这个目录下的所有符合字幕格式的文件
		var subSearchResult *sub_parser_hub.SearchSubResult
		subSearchResult, err = sub_parser_hub.SearchMatchedSubFileByDir(tmpDir)
		if err != nil {
			logger.Errorln("SearchMatchedSubFileByDir", nowFileFPath, err)
			return err
		}
		// 那么应该期待搜索出文字类型的字幕，最差也应该是图片类型或者蓝光的字幕
		if len(subSearchResult.Get(common.Characters)) > 0 {

			// 找到了文字类型的字幕
			subFileFullPaths = append(subFileFullPaths, subSearchResult.Get(common.Characters)...)
		} else if len(subSearchResult.Get(common.Picture)) > 0 {
			// 图片类型
			return common.SubtitleExtTypeIsPicture
		} else if len(subSearchResult.Get(common.BluRay)) > 0 {
			// 蓝光的字幕
			return common.SubtitleExtTypeIsBluRay
		} else {
			// 没有找到
			return common.SubtitleExtTypeNotSupported
		}
	}

	for i, subFileFullPath := range subFileFullPaths {

		logger.Infoln("Start Process Sub File", i+1, "/", len(subFileFullPaths), subFileFullPath)
		// 判断这个是电影还是电视剧
		if dSubInfo.IsMovie == true {
			// 这里就是电影
			err = h.processMovie(site, subFileFullPath, downloadDBId, downloadedInfo, dSubInfo)
			if err != nil {
				logger.Errorln("processMovie Error", err)
				continue
			}
		} else {
			// 连续剧的情况
			err = h.processTV(site, fullSeason, sourceFileSha256,
				subFileFullPath, downloadDBId, downloadedInfo, dSubInfo,
				subFromZipFile, nowFileFPath)
			if err != nil {
				logger.Errorln("processTV Error", err)
				continue
			}
		}
	}

	return nil
}

func (h HouseKeeping) processMovie(site models.WhichSite, subFileFullPath string, downloadDBId uint, downloadedInfo models.DownloadedInfo, dSubInfo models.DSubInfo) error {

	var err error
	var bok bool
	var subFileInfo *subparser.FileInfo
	bok, subFileInfo, err = h.subParserHub.DetermineFileTypeFromFile(subFileFullPath, site)
	if err != nil {
		logger.Errorln("DetermineFileTypeFromFile", subFileFullPath, err)
		return err
	}
	if bok == false {
		logger.Errorln("DetermineFileTypeFromFile", subFileFullPath, " == false")
		return errors.New("Not Supported")
	}
	// 创建电影的文件夹目录
	desSaveMovieDirPath := filepath.Join(settings.Get().HouseKeepingConfig.SubsSaveRootDirPath, common.Movie.String())
	if pkg.IsDir(desSaveMovieDirPath) == false {
		err = os.MkdirAll(desSaveMovieDirPath, os.ModePerm)
		if err != nil {
			logger.Errorln("MkdirAll", desSaveMovieDirPath, err)
			return err
		}
	}
	// 创建这个电影的目录
	desSaveMovieSubDirPath := filepath.Join(desSaveMovieDirPath, dSubInfo.ImdbId)
	if pkg.IsDir(desSaveMovieSubDirPath) == false {
		err = os.MkdirAll(desSaveMovieSubDirPath, os.ModePerm)
		if err != nil {
			logger.Errorln("MkdirAll", desSaveMovieSubDirPath, err)
			return err
		}
	}
	// 计算这个字幕的 sha256 值
	var nowSubSHA256 string
	nowSubSHA256, err = pkg.GetFileSHA256String(subFileFullPath)
	if err != nil {
		logger.Errorln("GetFileSHA256String", subFileFullPath, err)
		return err
	}
	// 构建这个字幕的存储位置信息
	saveFileName := nowSubSHA256 + filepath.Ext(filepath.Base(subFileFullPath))
	desSaveMovieSubFilePath := filepath.Join(desSaveMovieSubDirPath, saveFileName)
	err = pkg.CopyFile(subFileFullPath, desSaveMovieSubFilePath)
	if err != nil {
		logger.Errorln("CopyFile", subFileFullPath, desSaveMovieSubFilePath, err)
		return err
	}
	// 存入数据库
	var subtitleMovie models.SubtitleMovie
	subtitleMovie.SubSha256 = nowSubSHA256
	subtitleMovie.ImdbId = dSubInfo.ImdbId
	subtitleMovie.Title = filepath.Base(subFileFullPath)
	subtitleMovie.Language = (int)(subFileInfo.Lang)
	relPath, err := filepath.Rel(settings.Get().HouseKeepingConfig.SubsSaveRootDirPath, desSaveMovieSubFilePath)
	if err != nil {
		logger.Errorln("filepath.Rel", desSaveMovieSubFilePath, err)
		return err
	}
	subtitleMovie.SaveRelativePath = relPath
	subtitleMovie.Score = downloadedInfo.Score
	subtitleMovie.Votes = downloadedInfo.Votes
	subtitleMovie.DownloadTimes = downloadedInfo.DownloadTimes
	subtitleMovie.UploadTime = downloadedInfo.UploadTime
	subtitleMovie.SubtitlesComesFrom = downloadedInfo.SubtitlesComesFrom
	// 这个字幕是否已经存储到数据库过
	var subtitleMovies []models.SubtitleMovie
	var houseKeepings []models.HouseKeeping
	err = dao.Get().Transaction(func(tx *gorm.DB) error {

		// subtitle_movie 先要存储这个数据库的信息
		tx.Where("sub_sha256 = ? AND imdb_id = ?", nowSubSHA256, dSubInfo.ImdbId).Find(&subtitleMovies)
		if len(subtitleMovies) < 1 {
			// 没有存储过，那么就存储
			if err := tx.Create(&subtitleMovie).Error; err != nil {
				return err
			}
		}
		// HouseKeeping 然后才是这个数据库的信息
		tx.Where("downloaded_id = ? AND which_site = ?", downloadDBId, site.Index()).Find(&houseKeepings)
		if len(houseKeepings) < 1 {
			// 没有存储过，那么就存储
			var tmpHouseKeeping models.HouseKeeping
			tmpHouseKeeping.DownloadedSubId = downloadDBId
			tmpHouseKeeping.WhichSite = (int)(site)
			tmpHouseKeeping.ProcessTime = time.Now().Unix()
			if err := tx.Create(&tmpHouseKeeping).Error; err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		logger.Errorln("Transaction", err)
		return err
	}

	return nil
}

func (h HouseKeeping) processTV(site models.WhichSite, inFullSeason bool, fullSeasonFileSha256 string,
	subFileFullPath string, downloadDBId uint, downloadedInfo models.DownloadedInfo, dSubInfo models.DSubInfo,
	subFromZipFile bool, zipFullPath string) error {

	// inFullSeason 外层解析出来的 Season 信息，可以作为参考
	// 从解压的文件名称推断 Season 和 Episode 信息
	var err error
	var bok bool
	var subFileInfo *subparser.FileInfo
	var isFullSeason bool
	var nowSubFileSeason, nowSubFileEps int
	var cantParseSubtitleName bool
	// 解析单个字幕的 season 和 episode 信息
	isFullSeason, nowSubFileSeason, nowSubFileEps, err = pkg.GetSeasonAndEpisodeFromSubFileName(filepath.Base(subFileFullPath))
	if err != nil {
		/*
			这里还有其他的命名规则，比如
				X:\subtitles\zimuku\tv\tt3032476\6\[zmk.pw]Better.Call.Saul.S06E09\chs.srt
			那么，就再向上一层获取这个目录的名称，作为 season 和 episode 的信息，但是如果向上一层就到 ttxxxxx IMDB ID 文件夹，就依然是解析错误了
			就丢到 unknown 文件夹里面去
		*/
		parentFolder := filepath.Base(filepath.Dir(subFileFullPath))
		if strings.HasPrefix(parentFolder, "tt") == false {
			// 那么再次解析这个文件夹的名称
			_, nowSubFileSeason, nowSubFileEps, err = pkg.GetSeasonAndEpisodeFromSubFileName(parentFolder)
			if err != nil {
				// 说明不符合上面的情况，直接丢到 unknown 文件夹里面去
				if errors.Is(err, common.GetSeasonAndEpisodeFromSubFileNameError) == false {
					logger.Errorln("GetSeasonAndEpisodeFromSubFileName", subFileFullPath, err)
					return err
				}
				// 还是有可能是解析不出来的，那么这类文件就要放在季度文件夹中，但是需要额外的一个文件夹存储，unknown 文件夹
				cantParseSubtitleName = true
			}
		} else {

			// 那么如果这个字幕是从一个 zip 解压出来的，就可以判断这个 zip 的名称是否可以得到 Season 和 Episode 信息
			if subFromZipFile == true {

				zipName := filepath.Base(zipFullPath)
				_, nowSubFileSeason, nowSubFileEps, err = pkg.GetSeasonAndEpisodeFromSubFileName(zipName)
				if err != nil {
					// 说明不符合上面的情况，直接丢到 unknown 文件夹里面去
					if errors.Is(err, common.GetSeasonAndEpisodeFromSubFileNameError) == false {
						logger.Errorln("GetSeasonAndEpisodeFromSubFileName", subFileFullPath, err)
						return err
					}
					// 还是有可能是解析不出来的，那么这类文件就要放在季度文件夹中，但是需要额外的一个文件夹存储，unknown 文件夹
					cantParseSubtitleName = true
				}
			} else {
				// 说明不符合上面的情况，直接丢到 unknown 文件夹里面去
				if errors.Is(err, common.GetSeasonAndEpisodeFromSubFileNameError) == false {
					logger.Errorln("GetSeasonAndEpisodeFromSubFileName", subFileFullPath, err)
					return err
				}
				// 还是有可能是解析不出来的，那么这类文件就要放在季度文件夹中，但是需要额外的一个文件夹存储，unknown 文件夹
				cantParseSubtitleName = true
			}
		}
	}
	// 这里就不应该解析出 full season 的信息
	if isFullSeason == true {
		logger.Errorln("GetSeasonAndEpisodeFromSubFileName", subFileFullPath, "isFullSeason == true")
		return errors.New("GetSeasonAndEpisodeFromSubFileName isFullSeason == true")
	}

	bok, subFileInfo, err = h.subParserHub.DetermineFileTypeFromFile(subFileFullPath, site)
	if err != nil {
		logger.Errorln("DetermineFileTypeFromFile", subFileFullPath, err)
		return err
	}
	if bok == false {
		logger.Errorln("DetermineFileTypeFromFile", subFileFullPath, " == false")
		return errors.New("Not Supported")
	}
	// 创建连续剧的文件夹目录
	desSaveTVDirPath := filepath.Join(settings.Get().HouseKeepingConfig.SubsSaveRootDirPath, common.TV.String())
	if pkg.IsDir(desSaveTVDirPath) == false {
		err = os.MkdirAll(desSaveTVDirPath, os.ModePerm)
		if err != nil {
			logger.Errorln("MkdirAll", desSaveTVDirPath, err)
			return err
		}
	}
	// 创建这个连续剧的目录
	desSaveTVSubDirPath := filepath.Join(desSaveTVDirPath, dSubInfo.ImdbId)
	if pkg.IsDir(desSaveTVSubDirPath) == false {
		err = os.MkdirAll(desSaveTVSubDirPath, os.ModePerm)
		if err != nil {
			logger.Errorln("MkdirAll", desSaveTVSubDirPath, err)
			return err
		}
	}
	/*
		创建这一集的对应的 Season 目录
		这里传入的是一集的字幕文件，理论上可以从字幕文件名称中得到 nowSubFileSeason、nowSubFileEps
	*/
	desSaveTVSeasonSubDirPath := ""
	if cantParseSubtitleName == true {
		// 如果没有这个信息，那么就直接放入 unknown 文件夹
		desSaveTVSeasonSubDirPath = filepath.Join(desSaveTVSubDirPath, unknownFolderName)
	} else {
		// 正常解析出来这个信息了，那么就创建这个季度的文件夹
		desSaveTVSeasonSubDirPath = filepath.Join(desSaveTVSubDirPath, fmt.Sprintf("%d", nowSubFileSeason))
	}
	if pkg.IsDir(desSaveTVSeasonSubDirPath) == false {
		err = os.MkdirAll(desSaveTVSeasonSubDirPath, os.ModePerm)
		if err != nil {
			logger.Errorln("MkdirAll", desSaveTVSeasonSubDirPath, err)
			return err
		}
	}
	/*
		这里需要区分两种情况：
		1. 是全季字幕的情况，那么就需要使用这个季度包的 sha256 的值，创建一个文件夹，将全季的字幕都放入进去
		2. 不是全季字幕的情况，将这一集的字幕放入这个季度的文件夹进去
	*/
	if inFullSeason == true {
		// 1. 是全季字幕的情况
		// 在新建一层文件夹，将这个季度的字幕都放入进去
		desSaveTVSeasonSubDirPath = filepath.Join(desSaveTVSeasonSubDirPath, fullSeasonFileSha256)
		if pkg.IsDir(desSaveTVSeasonSubDirPath) == false {
			err = os.MkdirAll(desSaveTVSeasonSubDirPath, os.ModePerm)
			if err != nil {
				logger.Errorln("MkdirAll", desSaveTVSeasonSubDirPath, err)
				return err
			}
		}
	}
	// 计算这个字幕的 sha256 值
	var nowSubSHA256 string
	nowSubSHA256, err = pkg.GetFileSHA256String(subFileFullPath)
	if err != nil {
		logger.Errorln("GetFileSHA256String", subFileFullPath, err)
		return err
	}
	// 构建这个字幕的存储位置信息
	saveFileName := nowSubSHA256 + filepath.Ext(filepath.Base(subFileFullPath))
	desSaveMovieSubFilePath := filepath.Join(desSaveTVSeasonSubDirPath, saveFileName)
	err = pkg.CopyFile(subFileFullPath, desSaveMovieSubFilePath)
	if err != nil {
		logger.Errorln("CopyFile", subFileFullPath, desSaveMovieSubFilePath, err)
		return err
	}

	// 存入数据库
	var subtitleTV models.SubtitleTV

	subtitleTV.IsFullSeason = inFullSeason
	subtitleTV.FullSeasonSha256 = fullSeasonFileSha256
	subtitleTV.Season = nowSubFileSeason
	subtitleTV.Episode = nowSubFileEps
	subtitleTV.CantParseName = cantParseSubtitleName

	subtitleTV.SubSha256 = nowSubSHA256
	subtitleTV.ImdbId = dSubInfo.ImdbId
	subtitleTV.Title = filepath.Base(subFileFullPath)
	subtitleTV.Language = (int)(subFileInfo.Lang)
	relPath, err := filepath.Rel(settings.Get().HouseKeepingConfig.SubsSaveRootDirPath, desSaveMovieSubFilePath)
	if err != nil {
		logger.Errorln("filepath.Rel", desSaveMovieSubFilePath, err)
		return err
	}
	subtitleTV.SaveRelativePath = relPath
	subtitleTV.Score = downloadedInfo.Score
	subtitleTV.Votes = downloadedInfo.Votes
	subtitleTV.DownloadTimes = downloadedInfo.DownloadTimes
	subtitleTV.UploadTime = downloadedInfo.UploadTime
	subtitleTV.SubtitlesComesFrom = downloadedInfo.SubtitlesComesFrom
	// 这个字幕是否已经存储到数据库过
	var subtitleTVs []models.SubtitleTV
	var houseKeepings []models.HouseKeeping
	err = dao.Get().Transaction(func(tx *gorm.DB) error {

		// subtitle_movie 先要存储这个数据库的信息
		tx.Where("sub_sha256 = ? AND imdb_id = ?", nowSubSHA256, dSubInfo.ImdbId).Find(&subtitleTVs)
		if len(subtitleTVs) < 1 {
			// 没有存储过，那么就存储
			if err := tx.Create(&subtitleTV).Error; err != nil {
				return err
			}
		} else {
			// 已经存储过，那么就更新
			if err = tx.Model(&subtitleTVs[0]).Updates(&subtitleTV).Error; err != nil {
				return err
			}
		}
		// HouseKeeping 然后才是这个数据库的信息
		tx.Where("downloaded_id = ? and which_site = ?", downloadDBId, site.Index()).Find(&houseKeepings)
		if len(houseKeepings) < 1 {
			// 没有存储过，那么就存储
			var tmpHouseKeeping models.HouseKeeping
			tmpHouseKeeping.DownloadedSubId = downloadDBId
			tmpHouseKeeping.WhichSite = (int)(site)
			tmpHouseKeeping.ProcessTime = time.Now().Unix()
			if err := tx.Create(&tmpHouseKeeping).Error; err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		logger.Errorln("Transaction", err)
		return err
	}

	return nil
}

const unknownFolderName = "unknown"
