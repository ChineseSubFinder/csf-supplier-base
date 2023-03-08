package pkg

import (
	"os"
	"path/filepath"
)

// GetSubFixCacheFolderByName 获取缓存的文件夹，没有则新建
func GetSubFixCacheFolderByName(folderName string) (string, error) {
	rootPath, err := GetRootSubFixCacheFolder()
	if err != nil {
		return "", err
	}
	tmpFolderFullPath := filepath.Join(rootPath, folderName)
	if IsDir(tmpFolderFullPath) == false {
		err = os.MkdirAll(tmpFolderFullPath, os.ModePerm)
		if err != nil {
			return "", err
		}
	}
	return tmpFolderFullPath, nil
}

// GetRootSubFixCacheFolder 在程序的根目录新建，字幕时间校正的缓存文件夹
func GetRootSubFixCacheFolder() (string, error) {

	nowProcessRoot, err := os.Getwd()
	if err != nil {
		return "", err
	}
	nowProcessRoot = filepath.Join(nowProcessRoot, cacheRootFolderName, SubFixCacheFolder)
	if IsDir(nowProcessRoot) == false {
		err = os.MkdirAll(nowProcessRoot, os.ModePerm)
		if err != nil {
			return "", err
		}
	}
	return nowProcessRoot, err
}

const (
	cacheRootFolderName = "cache" // 缓存文件夹总名称
)

const (
	SubFixCacheFolder = "CSF-SubFixCache" // 字幕时间校正的缓存文件夹，一般可以不清理
)
