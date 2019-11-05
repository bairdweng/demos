package job

import (
	"iQuest/app/constant"
	"iQuest/app/model"
	"iQuest/db"
	"math"
	"time"

	"github.com/jinzhu/gorm"
)

type JobMemberList struct {
	Id            int       `json:"id"`
	WorkId        int       `json:"work_id"`
	Progress      int       `json:"progress"`
	Name          string    `json:"name"`
	ServiceTypeId int       `json:"service_type_id"`
	CreatedAt     time.Time `json:"created_at"`
}

type JobCondition struct {
	PageNum   int32
	PageSize  int32
	UserId    string
	CompanyId int
}

func GetJobMemberJoinWork(condition JobCondition) (map[string]interface{}, error) {
	var job_member_list []JobMemberList
	var page_info map[string]interface{}
	var count int32

	if err := db.Get().Table("job_member").Where("job_member.progress not in (?)", []string{"1", "2", "9"}).Where("work.status in (?)", []string{"1", "3"}).Where("participant_id = (?)", condition.UserId).Where("job_member.company_id = (?)", condition.CompanyId).Select("*").Joins("join work on work.id = job_member.work_id").Find(&job_member_list).Count(&count).Error; err != nil {
		return nil, err
	}

	page := condition.PageNum
	page_size := condition.PageSize
	if condition.PageNum == 0 {
		page = 1
	}

	if condition.PageSize == 0 {
		page_size = 10
	}
	offset := (page - 1) * page_size
	if err := db.Get().Table("job_member").Limit(page_size).Offset(offset).Where("job_member.progress not in (?)", []string{"1", "2", "9"}).Where("work.status in (?)", []string{"1", "3"}).Where("participant_id = (?)", condition.UserId).Where("job_member.company_id = (?)", condition.CompanyId).Select("job_member.id, job_member.work_id, job_member.progress, job_member.created_at, work.name, work.service_type_id").Joins("join work on work.id = job_member.work_id").Find(&job_member_list).Error; err != nil {
		return nil, err
	}

	page_info = make(map[string]interface{})
	page_info["total"] = count
	page_info["page_num"] = page
	page_info["page_size"] = page_size
	page_info["last_page"] = int32(math.Ceil(float64(count) / float64(page_size)))
	page_info["item"] = job_member_list

	return page_info, nil
}

// GetJobMembersByIds 获取进行中member
func GetJobMembersByIds(workID int, ids []int) ([]*model.Member, error) {
	var members []*model.Member
	query := db.Get().Unscoped().Where("work_id = ? and progress = ?", workID,constant.STATUS_APPROVE).Where("deleted_at = 0 or deleted_at is null")
	if len(ids) > 0 {
		query = query.Where("id in (?)", ids)
	}

	query.Preload("JobSettlementLogs", func(db *gorm.DB) *gorm.DB {
		return db.Order("id desc")
	}).Preload("ParticipantUser", func(db *gorm.DB) *gorm.DB {
		return db.Where("deleted_at = 0 or deleted_at is null").Unscoped()
	}).Find(&members)
	return members, nil
}
