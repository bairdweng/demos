package task

import "time"

type KpiUploadLog struct {
	ID        int64      `gorm:"primary_key" json:" - "` //
	CompanyId int64      `json:"company_id"`             //企业客户ID
	Appid         string     `json:"appid"`                     //应用标识
	FileUrl   string     `json:"file_url"`               //文件地址
	FileName  string     `json:"file_name"`              //文件名称
	Hash      string     `json:"hash"`                   //文件md5 32位小写
	CreatedAt *time.Time `json:"created_at"`             //
	UpdatedAt *time.Time `json:"updated_at"`             //
}

func (KpiUploadLog) TableName() string {
	return "iq_kpi_upload_log"
}
