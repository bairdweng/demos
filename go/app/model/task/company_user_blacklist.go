package task

import "time"

//
type CompanyUserBlacklist struct {
	ID            int64      `gorm:"primary_key" json:" - "` //
	CompanyId     int64      `json:"company_id"`             //企业客户ID
	PublisherId   int64      `json:"publisher_id"`           //发布者ID
	ParticipantId int64      `json:"participant_id"`         //C端用户ID
	Appid         string     `json:"appid"`                     //应用标识
	CreatedAt     *time.Time `json:"created_at"`             //
	UpdatedAt     *time.Time `json:"updated_at"`             //
	DeletedAt     *time.Time `json:"deleted_at"`             //
}

func (CompanyUserBlacklist) TableName() string {
	return "iq_company_user_blacklist"
}
