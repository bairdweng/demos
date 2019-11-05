package model

import (
	"iQuest/app/model/commonly"
	"time"
)

//
type Member struct {
	ID                int                            `gorm:"primary_key" json:" - "`                          //
	WorkId            int                            `gorm:"work_id" json:"work_id,omitempty"`                //
	PublisherId       string                         `gorm:"publisher_id" json:"publisher_id,omitempty"`      //
	ParticipantId     string                         `gorm:"participant_id" json:"participant_id,omitempty"`  //
	Source            int                            `gorm:"source" json:"source,omitempty"`                  //
	Progress          int                            `gorm:"progress" json:"progress,omitempty"`              //
	ProofFileUrl      string                         `gorm:"proof_file_url" json:"proof_file_url,omitempty"`  //
	Reward            float64                        `gorm:"reward" json:"reward,omitempty"`                  //
	ParticipateAt     int                            `gorm:"participarte_at" json:"participate_at,omitempty"` //
	FinishAt          int                            `gorm:"finish_at" json:"finish_at,omitempty"`            //
	Extend            string                         `gorm:"extend" json:"extend,omitempty"`                  //
	CreatedAt         time.Time                      `gorm:"created_at" json:"created_at,omitempty"`          //
	UpdatedAt         time.Time                      `gorm:"updated_at" json:"updated_at,omitempty"`          //
	DeletedAt         int                            `gorm:"deleted_at" json:"deleted_at,omitempty"`          //
	Remark            string                         `gorm:"remark" json:"remark,omitempty"`                  //
	UserData          int                            `gorm:"userData" json:"userData,omitempty"`              //
	Work              Work                           `gorm:"foreignkey:WorkId" `
	ParticipantUser   commonly.CommonlyUsedPersonnel `gorm:"foreignkey:ParticipantId;association_foreignkey:UserID"`
	JobSettlementLogs []JobSettlementLog             `gorm:"foreignkey:MemberID"`
}

func (m Member) TableName() string {
	return "job_member"
}
