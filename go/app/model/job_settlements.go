package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"iQuest/app/graphql/model"
	"time"
)

// JobSettlements 岗位结算记录总表
type JobSettlements struct {
	BaseModel
	BatchID     string     `gorm:"comment:'批次id';not null" json:"batch_id"`          //批次id
	WorkID      int        `gorm:"comment:'岗位id';not null" json:"work_id"`           //岗位id
	Amount      float64    `gorm:"comment:'绩效总金额';type:decimal(11,2)" json:"amount"` //绩效总金额
	SettleCount int        `gorm:"comment:'绩效总条数'" json:"settle_count"`              //绩效总条数
	Work        model.Work `gorm:"foreignkey:WorkID;association_foreignkey:ID"`
}

// BeforeCreate 生成sn
func (settle *JobSettlements) BeforeCreate(scope *gorm.Scope) error {
	var count int
	scope.DB().Model(settle).Where("work_id = ?", settle.WorkID).Count(&count)
	t := time.Now()
	nowTime := fmt.Sprintf("%d%d%d%d%d%d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	scope.SetColumn("batch_id", fmt.Sprintf("%s%06d%02d", nowTime, settle.WorkID, count+1))
	return nil
}
