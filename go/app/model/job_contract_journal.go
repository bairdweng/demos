package model

import "time"

//
type JobContractJournal struct {
	ID	int	`gorm:"primary_key" ` //
	ContractNo	string	`json:"contract_no"` //
	CompanyId	int32	`json:"company_id"` //
	ServiceCompanyId	int32	`json:"service_company_id"` //
	BeginTimestamp	int	`json:"begin_timestamp"` //
	EndTimestamp	int	`json:"end_timestamp"` //
	ActiveTimestamp	int	`json:"active_timestamp"` //
	IsHandled	int	`json:"is_handled"` //0未处理 1已处理
	Extend	string	`json:"extend"` //税筹推过来的原数据
	Remark	string	`json:"remark"` //
	CreatedAt	time.Time	`json:"created_at"` //
	UpdatedAt	time.Time	`json:"updated_at"` //
	DeletedAt	int	`json:"deleted_at"` //
}

func (m JobContractJournal) TableName() string {
	return "job_contract_journal"
}
