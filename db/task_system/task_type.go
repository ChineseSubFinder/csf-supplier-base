package task_system

// TaskType 定义任务的类型枚举类型
type TaskType int

const (
	// NoType 任务类型：未定义
	NoType TaskType = iota + 1
	// AudioToSubtitle 任务类型：音频转字幕
	AudioToSubtitle
	// Translation 任务类型：翻译
	Translation
)
