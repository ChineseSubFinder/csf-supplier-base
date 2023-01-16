package common

const (
	SubTypeASS = "ass"
	SubTypeSSA = "ssa"
	SubTypeSRT = "srt"
	// 这些都是文字字幕
	SubExtASS = ".ass"
	SubExtSSA = ".ssa"
	SubExtSRT = ".srt"
	// 图片字幕
	SubExtBDN = ".bdn"
	SubExtSST = ".sst"
	// 蓝光字幕
	SubExtSUP = ".sup"
)

type SubtitleType int

const (
	Characters   SubtitleType = iota + 1 // 文字类型
	Picture                              // 图片类型
	BluRay                               // 蓝光类型
	NotSupported                         // 不支持的类型
)
