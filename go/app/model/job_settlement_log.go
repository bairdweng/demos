package model

import (
	"fmt"
	"iQuest/app/model/commonly"

	"github.com/jinzhu/gorm"
)

// JobSettlementLog 岗位结算记录
type JobSettlementLog struct {
	BaseModel
	SettleID       string                         `gorm:"comment:'批次id';index:settle_id" json:"settle_id"`    //批次id
	WorkID         int                            `gorm:"comment:'岗位id';not null" json:"work_id"`             //岗位id
	MemberID       int                            `gorm:"comment:'参加id';not null" json:"member_id"`           //参加id
	UserID         string                         `gorm:"comment:'参加者用户ID';not null" json:"user_id"`          //参加者用户ID
	OperatorUserID string                         `gorm:"comment:'后台操作人ID';not null" json:"operator_user_id"` //后台操作人ID
	SN             string                         `gorm:"comment:'结算流水号';not null" json:"sn"`                 //结算流水号
	Amount         float64                        `gorm:"comment:'金额';type:decimal(11,2)" json:"amount"`      //金额
	File           string                         `gorm:"comment:'文件'" json:"file"`                           //文件
	ConfirmAt      int                            `gorm:"comment:'确认时间'" json:"confirm_at"`                   //确认时间
	Status         int                            `gorm:"comment:'状态(0:待确认,1:已确认)';default:0" json:"status"`  //状态
	User           commonly.CommonlyUsedPersonnel `gorm:"foreignkey:UserID;association_foreignkey:UserID"`
}

// BeforeCreate 生成sn
func (log *JobSettlementLog) BeforeCreate(scope *gorm.Scope) error {
	var count int
	scope.DB().Model(log).Where("work_id = ?", log.WorkID).Where("member_id = ?", log.MemberID).Count(&count)
	scope.SetColumn("sn", fmt.Sprintf("%05d%03d%02d", log.WorkID, log.MemberID, count+1))
	return nil
}
