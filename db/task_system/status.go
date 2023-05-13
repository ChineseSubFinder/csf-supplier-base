package task_system

// Status 任务的状态枚举类型
type Status int

// 数据库中的状态只能是 NotStart、 NoUploaded、Finished、Canceled、Failed
const (
	// NoUploaded 任务状态：未上传
	NoUploaded Status = iota + 1
	// NoAudited 任务状态：未审核
	NoAudited
	// NotStart 任务状态：未开始
	NotStart
	// Running 任务状态：进行中
	Running
	// Finished 任务状态：已完成
	Finished
	// Canceled 任务状态：已取消，人为取消任务则认为任务已取消
	Canceled
	// Failed 任务状态：已失败，超过一定时间没有处理完毕则认为失败
	Failed
)
