package commonly

import (
	"time"
)

type CommonlyUsed struct {
	Id        int64      `gorm:"primary_key" json:"id"`
	CompanyId int64      `json:"company_id"` //企业id
	UserId    int64      `json:"company_id"` //用户id
	AppId     int64      `json:"app_id"`     //appId
	CardNo    string     `json:"card_no"`    //身份证号码
	Education string     `json:"education"`  //学历
	Remark    string     `json:"remark"`     //备注
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// CommonlyUsedPersonnel 用户模型
type CommonlyUsedPersonnel struct {
	ID          int       `gorm:"column:id;primary_key" json:"id"`
	AppID       string    `gorm:"column:app_id" json:"app_id"`
	Mobile      string    `gorm:"column:mobile" json:"mobile"`
	Name        string    `gorm:"column:name" json:"name"`
	Remark      *string   `gorm:"column:remark" json:"remark,omitempty"`
	SigningTime *int      `gorm:"column:signing_time" json:"signing_time,omitempty"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`
	UserID      string    `gorm:"column:user_id" json:"user_id"`
	Address     *string   `gorm:"column:address" json:"address,omitempty"`
	Avatar      *string   `gorm:"column:avatar" json:"avatar,omitempty"`
	BankNo      string    `gorm:"column:bank_no" json:"bank_no"`
	CardNo      string    `gorm:"column:card_no" json:"card_no"`
	CompanyID   int       `gorm:"column:company_id" json:"company_id"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	DeletedAt   *int      `gorm:"column:deleted_at" json:"deleted_at,omitempty"`
	Education   *string   `gorm:"column:education" json:"education,omitempty"`
}
