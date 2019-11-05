package model

import "time"

//
type UserWeChatAuthorize struct {
	ID	int	`gorm:"primary_key" json:"id"` //
	UserId	string	`json:"user_id"` //
	Mobile	string	`json:"mobile"` //
	DeletedAt	int	`json:"deletedAt"` //
	CreatedAt	time.Time	`json:"createdAt"` //
	UpdatedAt	time.Time	`json:"updatedAt"` //
}

func (w UserWeChatAuthorize) TableName() string {
	return "user_wechat_authorize"
}