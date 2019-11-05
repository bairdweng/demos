package resume

import "iQuest/app/model"

//用户简历

type WorkUserResume struct {
	model.BaseModel
	WorkId       int    `json:"work_id" gorm:"comment:'关联岗位ID';not null"` //关联岗位ID
	UserID       string `json:"user_id"`
	ResumeId     int    `json:"resume_id" gorm:"comment:'简历项名称id';not null"`    //简历项id
	ResumeName   string `json:"resume_name" gorm:"comment:'简历项名称';not null"`    //简历项名称
	CategoryId   int    `json:"category_id" gorm:"comment:'所属类别ID';not null"`   //所属类别ID
	CategoryName string `json:"category_name" gorm:"comment:'所属类别名称';not null"` //所属类别名称
	Status       int    `json:"status" gorm:"comment:'状态 1，正常，0禁用';default:1"`  //状态 1，正常，0禁用
}

func (m WorkUserResume) TableName() string {
	return "work_user_resume"
}
