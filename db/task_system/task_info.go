package task_system

import (
	"time"
)

type TaskInfo struct {
	ID        int64     `gorm:"column:id;type:bigint;primary_key;auto_increment;not null"` // 任务的 ID
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;index;not null"`            // 任务创建的时间
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;index;not null"`            // 任务更新的时间

	TaskType        TaskType `gorm:"column:task_type;type:tinyint unsigned;index;not null"` // 任务的类型
	PackageID       string   `gorm:"column:package_id;type:char(64);index;not null"`        // 任务所属的任务包 ID
	TaskIndex       int      `gorm:"column:task_index;type:int;index;not null"`             // 任务在任务包中的索引, 从 0 开始
	Status          Status   `gorm:"column:status;type:tinyint unsigned;index;not null"`    // 任务的状态
	SrcDataRPath    string   `gorm:"column:src_data_r_path;type:varchar(255);not null"`     // 源任务数据的相对路径，相对于 R2 存储
	FinishDataRPath string   `gorm:"column:finish_data_r_path;type:varchar(255);not null"`  // 这个任务完成后，存储的数据的相对路径，相对于 R2 存储
	FinishDataSize  int      `gorm:"column:finish_data_size;type:int;not null"`             // 这个任务完成后，存储的数据的大小，单位：字节

	NeedSplit bool `gorm:"column:need_split;type:tinyint(1);not null;default:0"` // 是否需要拆分的任务，比如字幕的翻译，就需要分解成多个任务
	SplitDone bool `gorm:"column:split_done;type:tinyint(1);not null;default:0"` // 是否已经拆分完成

	NeedMerge bool `gorm:"column:need_merge;type:tinyint(1);not null;default:0"` // 是否需要合并的任务，比如字幕的翻译，就需要合并成一个任务
	MergeDone bool `gorm:"column:merge_done;type:tinyint(1);not null;default:0"` // 是否已经合并完成

	DataVersion string `gorm:"column:data_version;type:varchar(255);not null"` // 任务数据的版本

	SrcDataFileUrl    string `gorm:"-" json:"src_data_file_url"`    // 源任务数据的下载 URL
	FinishDataFileUrl string `gorm:"-" json:"finish_data_file_url"` // 这个任务完成后，存储的数据的上传 URL
}

func NewTaskInfo(taskType TaskType, packageID string,
	taskIndex int, status Status, dataVersion string) *TaskInfo {
	return &TaskInfo{TaskType: taskType, PackageID: packageID,
		TaskIndex: taskIndex, Status: status, DataVersion: dataVersion}
}
