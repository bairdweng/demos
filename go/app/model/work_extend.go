package model

import "time"

//
type WorkExtend struct {
	ID	int	`gorm:"primary_key" json:"id"` //
	WorkId	int	`json:"workId"` //
	AppId	string	`json:"appId"` //
	FocusCount	int	`json:"focusCount"` //
	LikeCount	int	`json:"likeCount"` //
	ShareCount	int	`json:"shareCount"` //
	CommentCount	int	`json:"comment_count"` //
	Extend	string	`json:"extend"` //
	DeletedAt	int	`json:"deletedAt"` //
	CreatedAt	time.Time	`json:"createdAt"` //
	UpdatedAt	time.Time	`json:"updatedAt"` //
}

func (w WorkExtend) TableName() string {
	return "work_extend"
}