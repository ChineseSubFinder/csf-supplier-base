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

// ToString 将任务类型转换为字符串
func (t TaskType) ToString() string {
	switch t {
	case AudioToSubtitle:
		return "AudioToSubtitle"
	case SplitSubtitle:
		return "SplitSubtitle"
	case Translation:
		return "Translation"
	case MergeSubtitle:
		return "MergeSubtitle"
	default:
		return "NoType"
	}
}
