package filter

import (
	"github.com/WQGroup/logger"
	"os"
	"path/filepath"
	"strings"
)

func SkipFileInfo(curFile os.FileInfo) bool {

	// 跳过不符合的文件，比如 MAC OS 下可能有缓存文件，见 #138
	if curFile.Size() < 1000 {
		logger.Debugln("curFile.Size() < 1000:", curFile.Name())
		return true
	}

	if curFile.Size() == 4096 && strings.HasPrefix(curFile.Name(), "._") == true {
		logger.Debugln("curFile.Size() == 4096 && Prefix Name == ._*", curFile.Name())
		return true
	}
	// 跳过预告片，见 #315
	if strings.HasSuffix(strings.ReplaceAll(curFile.Name(), filepath.Ext(curFile.Name()), ""), "-trailer") == true {
		logger.Debugln("curFile Name has -trailer:", curFile.Name())
		return true
	}

	return false
}
