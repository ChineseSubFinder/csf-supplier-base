package task_system

import "gorm.io/gorm"

type TaskInfo struct {
	gorm.Model
	TaskType        TaskType `gorm:"column:task_type;type:tinyint unsigned;index;not null"` // 任务的类型
	PackageID       string   `gorm:"column:package_id;type:char(64);index;not null"`        // 任务所属的任务包 ID
	TaskIndex       int      `gorm:"column:task_index;type:int;index;not null"`             // 任务在任务包中的索引, 从 0 开始
	Status          Status   `gorm:"column:status;type:tinyint unsigned;index;not null"`    // 任务的状态
	SrcDataRPath    string   `gorm:"column:src_data_r_path;type:varchar(255);not null"`     // 源任务数据的相对路径，相对于 R2 存储
	FinishDataRPath string   `gorm:"column:finish_data_r_path;type:varchar(255);not null"`  // 这个任务完成后，存储的数据的相对路径，相对于 R2 存储
	DataVersion     string   `gorm:"column:data_version;type:varchar(255);not null"`        // 任务数据的版本
}
