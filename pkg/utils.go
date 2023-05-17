package pkg

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/ChineseSubFinder/csf-supplier-base/pkg/common"
	"github.com/ChineseSubFinder/csf-supplier-base/pkg/regex_things"
	"github.com/WQGroup/logger"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func Base642IMGFile(imgType, base64String, tmpRootPath string) (string, error) {

	// Base64 解码
	decodeBytes, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		return "", err
	}

	if IsDir(tmpRootPath) == false {
		return "", err
	}
	// 写入文件
	desFPath := filepath.Join(tmpRootPath, RandStringBytesMaskImprSrcSB(32)+"."+imgType)
	err = WriteFile(desFPath, decodeBytes)
	if err != nil {
		return "", err
	}

	return desFPath, nil
}

// IsDir 存在且是文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// IsFile 存在且是文件
func IsFile(filePath string) bool {
	s, err := os.Stat(filePath)
	if err != nil {
		return false
	}
	return !s.IsDir()
}

func GetEpisodeKeyName(season, eps int, zerofill ...bool) string {

	if len(zerofill) < 1 || zerofill[0] == false {
		return "S" + strconv.Itoa(season) + "E" + strconv.Itoa(eps)
	} else {
		return fmt.Sprintf("S%02dE%02d", season, eps)
	}
}

// WriteFileByStrings 写文件
func WriteFileByStrings(desFileFPath string, strs []string) error {
	var err error
	nowDesPath := desFileFPath
	if filepath.IsAbs(nowDesPath) == false {
		nowDesPath, err = filepath.Abs(nowDesPath)
		if err != nil {
			return err
		}
	}
	// 创建对应的目录
	nowDirPath := filepath.Dir(nowDesPath)
	if IsDir(nowDirPath) == false {
		err = os.MkdirAll(nowDirPath, os.ModePerm)
		if err != nil {
			return err
		}
	}
	file, err := os.Create(nowDesPath)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	for _, str := range strs {
		_, err = fmt.Fprintln(file, str)
		if err != nil {
			return err
		}
	}

	return nil
}

// WriteFile 写文件
func WriteFile(desFileFPath string, bytes []byte) error {
	var err error
	nowDesPath := desFileFPath
	if filepath.IsAbs(nowDesPath) == false {
		nowDesPath, err = filepath.Abs(nowDesPath)
		if err != nil {
			return err
		}
	}
	// 创建对应的目录
	nowDirPath := filepath.Dir(nowDesPath)
	if IsDir(nowDirPath) == false {
		err = os.MkdirAll(nowDirPath, os.ModePerm)
		if err != nil {
			return err
		}
	}
	file, err := os.Create(nowDesPath)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	_, err = file.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}

// CopyFile copies a single file from src to dst
func CopyFile(src, dst string) error {
	var err error
	var srcFd *os.File
	var dstFd *os.File
	var srcInfo os.FileInfo

	if srcFd, err = os.Open(src); err != nil {
		return err
	}
	defer func() {
		_ = srcFd.Close()
	}()

	if dstFd, err = os.Create(dst); err != nil {
		return err
	}
	defer func() {
		_ = dstFd.Close()
	}()

	if _, err = io.Copy(dstFd, srcFd); err != nil {
		return err
	}
	if srcInfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcInfo.Mode())
}

// AddBaseUrl 判断传入的 url 是否需要拼接 baseUrl
func AddBaseUrl(baseUrl, url string) string {
	if strings.Contains(url, "://") {
		return url
	}
	return fmt.Sprintf("%s%s", baseUrl, url)
}

// ReplaceSpecString 替换特殊的字符
func ReplaceSpecString(inString string, rep string) string {
	return regex_things.RegMatchSpString.ReplaceAllString(inString, rep)
}

// ReplaceWindowsSpecString 替换 Windows 下的特殊字符
func ReplaceWindowsSpecString(inString, rep string) string {
	return regex_things.RegMathWindowsSpString.ReplaceAllString(inString, rep)
}

func GetNumber2Float(input string) (float32, error) {
	compile := regexp.MustCompile(regex_things.RegGetNumber)
	params := compile.FindStringSubmatch(input)
	if params == nil || len(params) == 0 {
		return -1, errors.New("get number not match")
	}
	fNum, err := strconv.ParseFloat(params[0], 32)
	if err != nil {
		return -1, errors.New("get number ParseFloat error")
	}
	return float32(fNum), nil
}

func GetNumber2int(input string) (int, error) {
	compile := regexp.MustCompile(regex_things.RegGetNumber)
	params := compile.FindStringSubmatch(input)
	if params == nil || len(params) == 0 {
		return -1, errors.New("get number not match")
	}
	fNum, err := strconv.Atoi(params[0])
	if err != nil {
		return -1, errors.New("get number ParseFloat error")
	}
	return fNum, nil
}

// GetSeasonAndEpisodeFromSubFileName 从文件名推断 季 和 集 的信息 Season Episode，这个应该是次要方案，优先还是从 nfo 文件获取这些信息
func GetSeasonAndEpisodeFromSubFileName(videoFileName string) (bool, int, int, error) {

	upperName := strings.ToUpper(videoFileName)
	// 先进行单个 Episode 的匹配
	// Killing.Eve.S02E01.Do.You.Know.How
	var re = regexp.MustCompile(`(?m)[\.\s]S(\d+).*?E(\d+)[\.\s]`)
	matched := re.FindAllStringSubmatch(upperName, -1)
	if matched == nil || len(matched) < 1 {
		// NCISLOSAngelesS06E05.chs.ass
		re = regexp.MustCompile(`(?m)S(\d+).*?E(\d+)`)
		matched = re.FindAllStringSubmatch(upperName, -1)
		if matched == nil || len(matched) < 1 {
			// Killing.Eve.S02.Do.You.Know.How
			// 看看是不是季度字幕打包
			re = regexp.MustCompile(`(?m)[\.\s]S(\d+)[\.\s]`)
			matched = re.FindAllStringSubmatch(upperName, -1)
			if matched == nil || len(matched) < 1 {
				return false, -1, -1, common.GetSeasonAndEpisodeFromSubFileNameError
			}
			season, err := GetNumber2int(matched[0][1])
			if err != nil {
				return false, -1, -1, err
			}
			return true, season, 0, nil
		}
	}

	// 一集的字幕
	season, err := GetNumber2int(matched[0][1])
	if err != nil {
		return false, -1, -1, err
	}
	episode, err := GetNumber2int(matched[0][2])
	if err != nil {
		return false, -1, -1, err
	}

	return false, season, episode, nil

}

// GetFileSHA256String 获取文件的 SHA256 字符串
func GetFileSHA256String(fileFPath string) (string, error) {

	fp, err := os.Open(fileFPath)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = fp.Close()
	}()

	partAll, err := io.ReadAll(fp)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", sha256.Sum256(partAll)), nil
}

func JugRetryTimes(times int) {
	if times > 100 {
		logger.Panicln("retry time to many, break", times)
	}
}

// UploadFile2R2 使用 resty 包，进行文件的上传到 R2 上
func UploadFile2R2(uploadURL string, filePath string) error {

	client := resty.New()
	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	resp, err := client.R().
		SetBody(fileBytes).
		SetContentLength(true).
		Put(uploadURL)
	if err != nil {
		return err
	}

	if resp != nil && len(resp.Body()) > 0 {
		return errors.New(string(resp.Body()))
	}

	return nil
}

// ClearFolder 清空文件夹
func ClearFolder(folderFullPath string) error {
	pathSep := string(os.PathSeparator)
	files, err := os.ReadDir(folderFullPath)
	if err != nil {
		return err
	}
	for _, curFile := range files {
		fullPath := folderFullPath + pathSep + curFile.Name()
		if curFile.IsDir() {
			err = os.RemoveAll(fullPath)
			if err != nil {
				return err
			}
		} else {
			// 这里就是文件了
			err = os.Remove(fullPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func GetMaxSizeFile(path string) string {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return ""
	}
	var maxFile os.FileInfo
	for _, file := range files {
		if maxFile == nil {
			maxFile = file
		} else {
			if file.Size() > maxFile.Size() {
				maxFile = file
			}
		}
	}
	return filepath.Join(path, maxFile.Name())
}

// GetFileSha256AndSize 获取文件的名称，文件的 SHA256 以及 文件的大小
func GetFileSha256AndSize(fileFPath string) (string, string, int, error) {

	if IsFile(fileFPath) == false {
		return "", "", -1, errors.New("file not exist")
	}
	// 获取文件的名称，文件的 SHA256 以及 文件的大小
	fileName := filepath.Base(fileFPath)
	fp, err := os.Open(fileFPath)
	if err != nil {
		return "", "", -1, err
	}
	defer func() {
		_ = fp.Close()
	}()

	partAll, err := io.ReadAll(fp)
	if err != nil {
		return "", "", -1, err
	}
	fileSha256 := fmt.Sprintf("%x", sha256.Sum256(partAll))
	fileSize := len(partAll)

	return fileName, fileSha256, fileSize, nil
}
