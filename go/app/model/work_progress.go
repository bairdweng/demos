package model

import "time"

//
type WorkProgress struct {
	ID	int	`gorm:"primary_key" json:"id"` //
	AppId	string	`json:"appId"` //
	ParticipantId	string	`json:"participantId"` //
	PublisherId	string	`json:"publisherId"` //
	WorkId	int	`json:"workId"` //
	Type	string	`json:"type"` //
	Extend	string	`json:"extend"` //
	CreatedAt	time.Time	`json:"createdAt"` //
	UpdatedAt	time.Time	`json:"updatedAt"` //
}

func (w WorkProgress) TableName() string {
	return "work_progress"
}
