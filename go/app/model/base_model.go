package model

import (
	"time"
)

//BaseModel 基础模型
type BaseModel struct {
	ID        int        `gorm:"primary_key;comment:'主键'" json:"id"`           //主键
	CreatedAt time.Time  `gorm:"comment:'创建时间'" json:"created_at"`             //创建时间
	UpdatedAt time.Time  `gorm:"comment:'更新时间'" json:"updated_at"`             //更新时间
	DeletedAt *time.Time `gorm:"comment:'删除时间'" sql:"index" json:"deleted_at"` //删除时间
}
