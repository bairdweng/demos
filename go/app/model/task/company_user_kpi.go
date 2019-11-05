package task

import "time"

//
type CompanyUserKpi struct {
	ID             int64      `gorm:"primary_key" json:" - "` //
	CompanyId      int64      `json:"company_id"`             //企业客户ID
	PublisherId    int64      `json:"publisher_id"`           //发布者ID
	ParticipantId  int64      `json:"participant_id"`         //C端用户ID
	Appid         string     `json:"appid"`                     //应用标识
	TaskId         int64      `json:"task_id"`                //任务(职位)ID
	Amount         float64    `json:"amount"`                 //金额
	KpiCoefficient float32    `json:"kpi_coefficient"`        //绩效系数% (eg:85.22%)
	IsPaid         int        `json:"is_paid"`                //是否已发放成功,0未发放,1成功
	PayAt          *time.Time `json:"pay_at"`                 //
	Extend         string     `json:"extend"`                 //冗余字段
	CreatedAt      *time.Time `json:"created_at"`             //
	UpdatedAt      *time.Time `json:"updated_at"`             //
}

func (CompanyUserKpi) TableName() string {
	return "iq_company_user_kpi"
}
