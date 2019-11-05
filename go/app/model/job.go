package model

import "time"

//
type Job struct {
	ID               int         `gorm:"primary_key" json:" - "`           //
	WorkId           int         `gorm:"work_id" json:"work_id,omitempty"` //
	Category         int         `json:"category,omitempty"`               //
	PayStatus        int         `json:"pay_status,omitempty"`             //
	Progress         int         `json:"progress,omitempty"`               //
	Quota            int         `json:"quota,omitempty"`                  //
	SingleRewardMax  float64     `json:"single_reward_max,omitempty"`      //
	SingleRewardMin  float64     `json:"single_reward_min,omitempty"`      //
	IsCanComment     int         `json:"is_can_comment,omitempty"`         //
	IsNeedProof      int         `json:"is_need_proof,omitempty"`          //
	ProofDescription string      `json:"proof_description,omitempty"`      //
	ProofType        int         `json:"proof_type,omitempty"`             //
	SettlementRule   string      `json:"settlement_rule,omitempty"`        //
	TemplateId       int32       `json:"template_id,omitempty"`            //
	Remark           string      `json:"remark,omitempty"`                 //
	Extend           string      `json:"extend,omitempty"`                 //
	CreatedAt        time.Time   `json:"created_at"`                       //
	UpdatedAt        time.Time   `json:"updated_at"`                       //
	DeletedAt        int         `json:"deleted_at"`                       //
	Template         JobTemplate `gorm:"ForeignKey:TemplateId"`
}
