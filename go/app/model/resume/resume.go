package resume

import (
	"iQuest/app/model"
)

//简历选项

type Resume struct {
	model.BaseModel
	Name         string `json:"name" gorm:"comment:'简历项名称';not null"`           //简历项名称
	Status       int    `json:"status" gorm:"comment:'状态 1，正常，0禁用';default:1"`  //状态 1，正常，0禁用
	CategoryId   int    `json:"category_id" gorm:"comment:'所属类别ID';not null"`   //所属类别ID
	CategoryName string `json:"category_name" gorm:"comment:'所属类别名称';not null"` //所属类别名称
	Type         int    `json:"type" gorm:"comment:'类型 1不可填写，2可填写';default:1"`  //类型 1不可填写，2可填写
	BelongToId   int    `json:"belong_to_id" gorm:"comment:'所属其他项ID';default:0"`
	IsDefault    int 	`json:"is_default" gorm:"comment:'是否是默认';default:0"`
}

func (m Resume) TableName() string {
	return "resume"
}
