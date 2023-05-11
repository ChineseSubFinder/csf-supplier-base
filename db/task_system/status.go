package task_system

// Status 任务的状态枚举类型
type Status int

const (
	// NoAudited 任务状态：未审核
	NoAudited Status = iota + 1
	// NotStart 任务状态：未开始
	NotStart
	// Running 任务状态：进行中
	Running
	// Finished 任务状态：已完成
	Finished
	// Canceled 任务状态：已取消
	Canceled
	// Failed 任务状态：已失败
	Failed
)
