package job

import (
	"iQuest/app/model"
	"iQuest/db"
)

// CreateJobTemplate 创建岗位模版
func CreateJobTemplate(job *model.JobTemplate) (*model.JobTemplate, error) {
	if err := db.Get().Create(&job).Error; err != nil {
		return nil, err
	}

	return job, nil
}
