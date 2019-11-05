package model

import "time"

//
type CompanyUserBlacklist struct {
	ID	int	`gorm:"primary_key" json:" - "` //
	CompanyId	int	`json:"company_id"` //
	ParticipantId	string	`json:"participant_id"` //
	PublisherId	string	`json:"publisher_id"` //
	CreatedAt	time.Time	`json:"created_at"` //
	UpdatedAt	time.Time	`json:"updated_at"` //
	DeletedAt	int	`json:"deleted_at"` //
	WorkId	int	`json:"work_id"` //
	Type	int	`json:"type"` //
}


func (cub CompanyUserBlacklist) TableName() string {
	return "company_user_blacklist"
}

