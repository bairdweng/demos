package model

import (
	"github.com/jinzhu/gorm"
	"iQuest/app/Api"
	"iQuest/app/constant"
	"iQuest/db"
	"strconv"
	"time"
)

//
type Work struct {
	ID            int       `gorm:"primary_key" json:" - "`   //
	AppId         string    `json:"app_id,,omitempty"`        //
	CompanyId     int32     `json:"companyId,,omitempty"`     //
	UserId        string    `json:"userId,,omitempty"`        //
	ServiceTypeId int32     `json:"serviceTypeId,,omitempty"` //
	WorkType      int       `json:"workType,,omitempty"`      //
	Name          string    `json:"name,,omitempty"`          //
	Requirement   string    `json:"requirement,,omitempty"`   //
	PayType       int       `json:"payType,,omitempty"`       //
	Duration      int       `json:"duration,,omitempty"`      //
	EndAt         int       `json:"endAt,,omitempty"`         //
	Source        int       `json:"source,,omitempty"`        //
	Status        int       `json:"status,,omitempty"`        //
	Type          int       `json:"type,,omitempty"`          //
	IsPublic      int       `json:"isPublic,,omitempty"`      //
	MediaCoverUrl string    `json:"mediaCoverUrl,,omitempty"` //
	MediaUrls     string    `json:"mediaUrls,,omitempty"`     //
	Resume     string    `json:"resume,,omitempty"`     //
	Extend        string    `json:"extend,,omitempty"`        //
	UpdatedAt     time.Time `json:"updatedAt,,omitempty"`     //
	CreatedAt     time.Time `json:"createdAt,,omitempty"`     //
	DeletedAt     int       `json:"deletedAt,,omitempty"`     //
	Members       []Member  `gorm:"foreignkey:WorkId"`
	Job           Job       `gorm:"foreignkey:WorkId"`
}

func (w *Work) CanJoins(userId string, offset, limit int ,appId *string) (int, *[]Work, error) {
	var total = 0
	var works = make([]Work, 0)
	var err error

	var joinedIds []int64
	err = db.Get().Model(&Member{}).Unscoped().Where("participant_id = ? and progress in (?)", userId, constant.JOB_STATUS_ING ).Pluck("work_id", &joinedIds).Error
	if err != nil {
		return total, &works, err
	}

	//获取身份证
	userIdInt,_ := strconv.Atoi(userId)
	ver, err := Api.FindRealnameInfoByUserid(int64(userIdInt))
	if err != nil {
		return total, &works, err
	}
	//获取签约列表
	SignList,_ := Api.GetSignInfo(ver.Data.CredentialsNo)
	var companyList []string
	if len(SignList.Data) > 0 {
		for _, value := range SignList.Data {
			companyList = append(companyList, value.CompanyId)
		}
	}

	var companyBlackList []string
	err = db.Get().Model(&CompanyUserBlacklist{}).Unscoped().Where("participant_id = ?", userId).Pluck("company_id", &companyBlackList).Error
	if len(companyBlackList) <= 0 {
		companyBlackList = []string{"0"}
	}

	/*
		薪鸟小程序中
		1.未签约的公司,可以看薪鸟的
		2.和薪鸟指定公司签约的,可以看薪鸟指定公司的
	*/
	if appId != nil &&  *appId == constant.XINNIAO_APPID{
		//if len(SignList.Data) == 0 {
			//未签约
			for i := 0; i < len(constant.XINNIAO_COMPANY_ID_ARR); i++ {
				tempId := strconv.Itoa(int(constant.XINNIAO_COMPANY_ID_ARR[i]))
				companyList = append(companyList, tempId)
			}
		//}

	}

	db := db.Get().Model(w).
		Unscoped().
		Where("status = ? and is_public = ? ", constant.WORK_STATUS_NORMAL, constant.WORK_STATUS_NORMAL).
		Where("company_id in (?)", companyList).
		Where("company_id not in (?)", companyBlackList).
		Preload("Job").
		Count(&total).
		Order("id desc").
		Offset(offset).
		Limit(limit)
	if len(joinedIds) > 0 {
		db = db.Not("id", joinedIds)
	}

	err = db.Find(&works).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return total, &works, err
	}
	return total, &works, nil

}
