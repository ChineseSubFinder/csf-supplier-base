package task_system

// TaskType 定义任务的类型枚举类型
type TaskType int

const (
	// NoType 任务类型：未定义
	NoType TaskType = iota + 1
	// AudioToSubtitle 任务类型：音频转字幕
	AudioToSubtitle
	// SplitSubtitle 任务类型：拆分字幕
	SplitSubtitle
	// Translation 任务类型：翻译
	Translation
	// MergeSubtitle 任务类型：合并翻译后的字幕
	MergeSubtitle
)
