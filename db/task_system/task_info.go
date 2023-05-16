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
	SrcDataSize     int      `gorm:"column:src_data_size;type:int;not null"`                // 源任务数据的大小，单位：字节
	FinishDataRPath string   `gorm:"column:finish_data_r_path;type:varchar(255);not null"`  // 这个任务完成后，存储的数据的相对路径，相对于 R2 存储
	FinishDataSize  int      `gorm:"column:finish_data_size;type:int;not null"`             // 这个任务完成后，存储的数据的大小，单位：字节

	PartIndex int `gorm:"column:part_index;type:int;not null"` // 任务的分片索引，从 0 开始
	PartCount int `gorm:"column:part_count;type:int;not null"` // 任务的分片总数

	DataVersion string `gorm:"column:data_version;type:varchar(255);not null"` // 任务数据的版本

	SrcDataFileUrl    string `gorm:"-" json:"src_data_file_url"`    // 源任务数据的下载 URL
	FinishDataFileUrl string `gorm:"-" json:"finish_data_file_url"` // 这个任务完成后，存储的数据的上传 URL
}

func NewTaskInfo(taskType TaskType, packageID string,
	taskIndex int, status Status, dataVersion string) *TaskInfo {
	return &TaskInfo{
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
		TaskType: taskType, PackageID: packageID,
		TaskIndex: taskIndex, Status: status, DataVersion: dataVersion}
}

func (t *TaskInfo) SetStatus(status Status) {
	t.Status = status
	t.UpdatedAt = time.Now()
}
