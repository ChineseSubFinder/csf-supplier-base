package pkg

import (
	"github.com/ChineseSubFinder/csf-supplier-base/pkg/common"
	PTN "github.com/middelink/go-parse-torrent-name"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
)

// GetVideoInfoFromFileName 从文件名推断文件信息，这个应该是次要方案，优先还是从 nfo 文件获取这些信息
func GetVideoInfoFromFileName(fileName string) (*PTN.TorrentInfo, error) {

	parse, err := PTN.Parse(fileName)
	if err != nil {
		return nil, err
	}
	compile, err := regexp.Compile(regFixTitle2)
	if err != nil {
		return nil, err
	}
	match := compile.ReplaceAllString(parse.Title, "")
	match = strings.TrimRight(match, "")
	parse.Title = match

	return parse, nil
}

func GetIMDBIdFromIMDBUrl(url string) string {

	regIMDBId := regexp.MustCompile(`tt\d+`)
	match := regIMDBId.FindStringSubmatch(url)
	if len(match) == 1 {
		return match[0]
	}
	return ""
}

func GetDouBanIdFromDouBanUrl(inUrl string) string {

	baseUrlInfo, err := url.Parse(inUrl)
	if err != nil {
		return ""
	}
	nowPath := strings.Trim(baseUrlInfo.Path, "/")
	id := strings.ReplaceAll(nowPath, "subject", "")
	id = strings.Trim(id, "/")
	return id
}

// IsFakeBDMVWorked 传入的是伪造的不存在的蓝光结构的视频全路径，如果是就返回 true 和 id.bdmv 的绝对路径 和 STREAM 绝对路径
func IsFakeBDMVWorked(fakseVideFPath string) (bool, string, string) {

	rootDir := filepath.Dir(fakseVideFPath)

	CERDir := filepath.Join(rootDir, "CERTIFICATE")
	BDMVDir := filepath.Join(rootDir, "BDMV")
	STREAMDir := filepath.Join(BDMVDir, "STREAM")
	idBDMVFPath := filepath.Join(CERDir, common.FileBDMV)

	if IsDir(CERDir) == true && IsDir(BDMVDir) == true && IsFile(idBDMVFPath) == true {
		return true, idBDMVFPath, STREAMDir
	}

	return false, "", ""
}

const (
	// 去除特殊字符，把特殊字符都写进去
	regFixTitle2 = "[~!@#$%^&*:()\\+\\-=|{}';'\\[\\].<>/?~！@#￥%……&*（）——+|{}【】'；”“’。、？]"
)
