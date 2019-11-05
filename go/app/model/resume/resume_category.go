package resume

import (
	"iQuest/app/model"
)

//简历类别   `json:"sn" gorm:"comment:'结算流水号';not null"`
type ResumeCategory struct {
	model.BaseModel
	Name   string `json:"name" gorm:"comment:'类别名称';not null"`           //类别名称
	Status int    `json:"status" gorm:"comment:'状态 1，正常，0禁用';default:1"` //状态 1，正常，0禁用
	Type int `json:"type" gorm:"comment:'类型 1单选，2多选';default:1"`  //类型 1单选，2多选
}

func (m ResumeCategory) TableName() string {
	return "resume_category"
}
