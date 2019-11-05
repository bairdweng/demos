package model

import "time"

//
type JobTemplateDelay struct {
	ID                 int       `gorm:"primary_key" json:" - "`       //
	ContractNo         string    `json:"contract_no,omitempty"`        //
	CompanyId          int32     `json:"company_id,omitempty"`         //
	ServiceCompanyId   int32     `json:"service_company_id,omitempty"` //
	BeginTimestamp     int       `json:"begin_timestamp,omitempty"`    //
	EndTimestamp       int       `json:"end_timestamp,omitempty"`      //
	IsNeedProcessBegin int       `json:"is_process_begin,omitempty"`   //
	IsNeedProcessEnd   int       `json:"is_process_end,omitempty"`     //
	Extend             string    `json:"extend,omitempty"`             //
	BizContent         string    `json:"biz_content,omitempty"`        //
	Remark             string    `json:"remark,omitempty"`             //
	CreatedAt          time.Time `json:"created_at"`                   //
	UpdatedAt          time.Time `json:"updated_at"`                   //
	DeletedAt          int       `json:"deleted_at"`                   //
}

//是否能发岗
func (p *JobTemplateDelay) CanPublishJob() bool {
	return int64(p.BeginTimestamp) <= time.Now().Unix() && int64(p.EndTimestamp) >  time.Now().Unix()
}
