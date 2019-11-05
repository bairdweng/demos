package work

import (
	"errors"
	"iQuest/app/model"
	"iQuest/db"
	"math"
	//"time"
)

type WorkModel struct {
	ID              int     `gorm:"primary_key" json:"id" `
	AppId           string  `json:"app_id,omitempty"`
	CompanyId       int32   `json:"company_id,omitempty"`
	CompanyName     string  `json:"company_name,omitempty"`
	JobId           int64   `json:"job_id,omitempty"`
	MediaCoverUrl   string  `json:"media_cover_url,omitempty"`
	MediaUrls       string  `json:"media_urls,omitempty"`
	Name            string  `json:"name,omitempty"`
	Requirement     string  `json:"requirement,omitempty"`
	ServiceTypeId   int32   `json:"service_type_id,omitempty"`
	ServiceTypeName string  `json:"service_type_name,omitempty"`
	WorkType        int32   `json:"work_type,omitempty"`
	SingleRewardMin float64 `json:"single_reward_min,omitempty"`
	SingleRewardMax float64 `json:"single_reward_max,omitempty"`
}

type WorkListCondition struct {
	CompanyId int
	UserId    string
	Status    int
	PageNum   int
	PageSize  int
}

// GetByID 获取work
func GetByID(id int) (*model.Work, error) {
	var work model.Work
	var err error
	query := db.Get().Unscoped().Where("deleted_at = 0 or deleted_at is null").First(&work, id)
	if query.Error != nil {
		err = query.Error
	}
	if query.RecordNotFound() {
		err = errors.New("数据不存在")
	}

	return &work, err
}

func GetListByCondition(condition WorkListCondition) (map[string]interface{}, error) {
	var work_list []WorkModel
	var count int32

	work_db := db.Get()

	if condition.Status != 0 {
		work_db = work_db.Where("status = (?)", condition.Status)
	}

	if condition.UserId != "" {
		work_db = work_db.Joins("join job_member on work.id=job_member.work_id").Where("job_member.participant_id = (?)", condition.UserId)
	}
	work_db_list := work_db
	if err := work_db.Table("work").
		Where("work.deleted_at = (?)", 0).
		Where("job.deleted_at = (?)", 0).
		Joins("join job on work.id=job.work_id").
		Select("*").
		Unscoped().
		Count(&count).
		Error; err != nil {
		return nil, err
	}

	//if err := work_db.Model(&model.Work{}).Where("deleted_at = (?)", 0).
	//	Preload("Job", func(db *gorm.DB) *gorm.DB {
	//		return db.Unscoped().Where("deleted_at = (?)", 0)
	//	}).
	//	Unscoped().
	//	Count(&count).Error; err != nil {
	//	return nil, err
	//}

	page := condition.PageNum
	page_size := condition.PageSize
	if condition.PageNum == 0 {
		page = 1
	}

	if condition.PageSize == 0 {
		page_size = 10
	}
	offset := (page - 1) * page_size

	if err := work_db_list.Table("work").
		Limit(page_size).
		Offset(offset).
		Where("work.deleted_at = (?)", 0).
		Where("job.deleted_at = (?)", 0).
		Joins("join job on work.id=job.work_id").
		Select("work.id, work.work_type, work.company_id, work.name, work.requirement, work.media_cover_url, work.media_urls, job.id as job_id, job.single_reward_min, job.single_reward_max,job_template.app_id, job_template.company_name, job_template.service_type_id, job_template.service_type_name").
		Joins("join job_template on job.template_id=job_template.id").
		Unscoped().
		Find(&work_list).
		Error; err != nil {
		return nil, err
	}

	out_put_works := []interface{}{}

	for i := 0; i < len(work_list); i++ {
		work_null_map := make(map[string]interface{})

		work_null_map["work_id"] = work_list[i].ID
		work_null_map["app_id"] = work_list[i].AppId
		work_null_map["work_type"] = work_list[i].WorkType
		work_null_map["company_id"] = work_list[i].CompanyId
		work_null_map["company_name"] = work_list[i].CompanyName
		work_null_map["service_type_id"] = work_list[i].ServiceTypeId
		work_null_map["service_type_name"] = work_list[i].ServiceTypeName
		work_null_map["name"] = work_list[i].Name
		work_null_map["requirement"] = work_list[i].Requirement
		work_null_map["media_cover_url"] = work_list[i].MediaCoverUrl
		work_null_map["media_urls"] = work_list[i].MediaUrls
		work_null_map["job_id"] = work_list[i].JobId
		work_null_map["single_reward_min"] = work_list[i].SingleRewardMin
		work_null_map["single_reward_max"] = work_list[i].SingleRewardMax

		out_put_works = append(out_put_works, work_null_map)

	}

	page_info := make(map[string]interface{})
	page_info["total"] = count
	page_info["page_num"] = page
	page_info["page_size"] = page_size
	page_info["last_page"] = int32(math.Ceil(float64(count) / float64(page_size)))
	page_info["item"] = out_put_works

	return page_info, nil
}

// CreateWork 创建岗位
func CreateWork(work *model.Work) (*model.Work, error) {
	if err := db.Get().Create(&work).Error; err != nil {
		return nil, err
	}

	return work, nil
}
