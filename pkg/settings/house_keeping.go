package settings

type HouseKeepingConfig struct {
	SubsSaveRootDirPath string // 字幕存储的根目录，后续在数据库中只会存储相对路径
	TmpRootDirPath      string // 缓存目录
}
