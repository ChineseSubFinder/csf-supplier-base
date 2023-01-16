package imdb_info_center

import (
	"compress/gzip"
	"fmt"
	"github.com/ChineseSubFinder/csf-supplier/pkg"
	"github.com/ChineseSubFinder/csf-supplier/pkg/settings"
	"github.com/ChineseSubFinder/csf-supplier/pkg/struct_json"
	"github.com/WQGroup/logger"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func DownloadFiles(tmpDir string, downloadList []string) {

	logger.Infof("Download Files start, save at: %s", tmpDir)
	var wg sync.WaitGroup
	for _, item := range downloadList {
		nowUrl := item
		wg.Add(1)
		go download(&wg, item, tmpDir, nowUrl)
	}
	wg.Wait()
	logger.Infof("Download Finished")
}

func DecompressFiles(tmpDir string) {
	defer logger.Infoln("Decompress Finished")
	logger.Infoln("DecompressFiles start")

	for _, filename := range settings.Get().ImdbInfoCenterConfig.Files {

		logger.Infoln("decompress file:", filename)
		fileFPath := filepath.Join(tmpDir, filename)
		_, err := unpackGzipFile(fileFPath)
		if err != nil {
			logger.Panicln("[ERROR] Unzip file:", err)
		}
	}
}

func ClearDownloadedAndDecompressFiles(tmpDir string) {

	for _, filename := range settings.Get().ImdbInfoCenterConfig.Files {

		// 下载的压缩包
		fileFPath := filepath.Join(tmpDir, filename)
		if pkg.IsFile(fileFPath) == true {
			logger.Infoln("clear file:", fileFPath)
			_ = os.Remove(fileFPath)
		}
		// 解压出来的文件
		decompressFilePath := strings.ReplaceAll(fileFPath, filepath.Ext(fileFPath), "")
		if pkg.IsFile(decompressFilePath) == true {
			logger.Infoln("clear file:", decompressFilePath)
			_ = os.Remove(decompressFilePath)
		}
	}
}

func download(wg *sync.WaitGroup, name, tmpDir string, url string) {

	defer func() {
		logger.Infoln("Download Finished : ", name)
		wg.Done()
	}()

	logger.Infoln("Download file : ", name, "save at:", tmpDir)

	resp, err := http.Get(url)
	if err != nil {
		logger.Panicf("%s: %v", name, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("non-200 status: %s", resp.Status)
		logger.Panicf("%s: %v", name, err)
	}
	// create dest
	destName := filepath.Base(url)
	dest, err := os.Create(filepath.Join(tmpDir, destName))
	if err != nil {
		err = fmt.Errorf("can't create %s: %v", destName, err)
		logger.Panicf("%s: %v", name, err)
	}
	_, err = io.Copy(dest, resp.Body)
	if closeErr := dest.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		logger.Panicf("%s: %v", name, err)
	}
}

func unpackGzipFile(gzFilePath string) (int64, error) {

	dstFilePath := strings.ReplaceAll(gzFilePath, filepath.Ext(gzFilePath), "")
	gzFile, err := os.Open(gzFilePath)
	if err != nil {
		return 0, fmt.Errorf("failed to open file %s for unpack: %s", gzFilePath, err)
	}
	dstFile, err := os.OpenFile(dstFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		return 0, fmt.Errorf("failed to create destination file %s for unpack: %s", dstFilePath, err)
	}

	ioReader, ioWriter := io.Pipe()

	go func() { // goroutine leak is possible here
		gzReader, _ := gzip.NewReader(gzFile)
		// it is important to close the writer or reading from the other end of the
		// pipe or io.copy() will never finish
		defer func() {
			gzFile.Close()
			gzReader.Close()
			ioWriter.Close()
		}()

		io.Copy(ioWriter, gzReader)
	}()

	written, err := io.Copy(dstFile, ioReader)
	if err != nil {
		return 0, err // goroutine leak is possible here
	}
	ioReader.Close()
	dstFile.Close()

	return written, nil
}

const DownloadImdbDataInfo = "imdb_data_cache_info.json"

type ImdbDataCacheInfo struct {
	DownloadedTime int64
}

func NewImdbDataCacheInfo(downloadedTime int64) *ImdbDataCacheInfo {
	return &ImdbDataCacheInfo{DownloadedTime: downloadedTime}
}

func GetDownloadedCacheTime(cacheRootDirPath string) *ImdbDataCacheInfo {

	saveFPath := filepath.Join(cacheRootDirPath, DownloadImdbDataInfo)
	if pkg.IsFile(saveFPath) == false {
		// 需要保存一个新的
		info := ImdbDataCacheInfo{
			DownloadedTime: time.Now().Unix(),
		}
		err := struct_json.ToFile(saveFPath, info)
		if err != nil {
			logger.Panicln("save imdb data cache info failed: ", err)
		}
		return &info
	} else {
		// 如果存在，那么就直接读取
		info := ImdbDataCacheInfo{}
		err := struct_json.ToStruct(saveFPath, &info)
		if err != nil {
			logger.Panicln("read imdb data cache info failed: ", err)
		}
		return &info
	}
}

func SetDownloadedCacheTime(cacheRootDirPath string, info *ImdbDataCacheInfo) {
	saveFPath := filepath.Join(cacheRootDirPath, DownloadImdbDataInfo)
	err := struct_json.ToFile(saveFPath, *info)
	if err != nil {
		logger.Panicln("save imdb data cache info failed: ", err)
	}
}
