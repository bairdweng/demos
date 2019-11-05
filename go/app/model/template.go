package model

import (
	"iQuest/app/constant"
	"iQuest/db"
	"time"
)

//
type JobTemplate struct {
	ID                 int64            `gorm:"primary_key" json:" - "`       //
	AppId              string           `json:"appId,omitempty"`              //
	CompanyId          int32            `json:"companyId,omitempty"`          //
	CompanyName        string           `json:"companyName,omitempty"`        //
	ServiceCompanyId   int32            `json:"serviceCompanyId,omitempty"`   //
	ServiceCompanyName string           `json:"serviceCompanyName,omitempty"` //
	UserId             string           `json:"userId,omitempty"`             //
	WorkType           int              `json:"workType,omitempty"`           //
	SignTemplateId     int              `json:"signTemplateId,omitempty"`     //
	ServiceTypeId      int32            `json:"serviceTypeId,omitempty"`      //
	ServiceTypeName    string           `json:"serviceTypeName,omitempty"`    //
	Name               string           `json:"name,omitempty"`               //
	Requirement        string           `json:"requirement,omitempty"`        //
	SettlementRule     string           `json:"settlementRule,omitempty"`     //
	KpiTemplateUrl     string           `json:"kpiTemplateUrl,omitempty"`     //
	Remark             string           `json:"remark,omitempty"`             //
	Extend             string           `json:"extend,omitempty"`             //
	UpdatedAt          time.Time        `json:"updatedAt,omitempty"`          //
	CreatedAt          time.Time        `json:"createdAt,omitempty"`          //
	DeletedAt          int              `json:"deletedAt,omitempty"`          //
	IsEnable           int              `json:"isEnable,omitempty"`           //
	BizContent         string           `json:"bizContent,omitempty"`         //
	Source             int              `json:"source,omitempty"`             //
	DownloadCode       string           `json:"downloadCode,omitempty"`       //
	PlatformAppid      string           `json:"platformAppid,omitempty"`      //
	DisplayName        string           `json:"displayName,omitempty"`        //
	Jobs               []Job            `gorm:"ForeignKey:TemplateId"`
	ContractNo         string           `form:"contract_no" json:"contractNo,omitempty"`
	Contract           JobTemplateDelay `gorm:"foreignkey:contractNo"`
}

func (j *JobTemplate) BatchPublishWorkByAudit(templateId int64) error {
	return db.Get().Model(j).Where("template_id = ?", templateId).UpdateColumn("status", constant.WORK_STATUS_NORMAL).Error
}
