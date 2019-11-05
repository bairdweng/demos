package import_logs

import (
	logs "iQuest/app/model/import_logs"
	"iQuest/db"
	"math"
)

type GetListCondition struct {
	PageNum  int32
	PageSize int32
	FileHash string
	Status   int
}

func GetListWithPage(condition GetListCondition) (map[string]interface{}, error) {
	var import_logs []logs.ImportLogs
	var page_info map[string]interface{}
	var count int32

	if err := db.Get().Where("file_hash = (?) and status = (?)", condition.FileHash, condition.Status).Find(&import_logs).Count(&count).Error; err != nil {
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
	if err := db.Get().Limit(page_size).Offset(offset).Where("file_hash = (?) and status = (?)", condition.FileHash, condition.Status).Find(&import_logs).Error; err != nil {
		return nil, err
	}

	page_info = make(map[string]interface{})

	page_info["total"] = count
	page_info["page_num"] = page
	page_info["page_size"] = page_size
	page_info["last_page"] = int32(math.Ceil(float64(count) / float64(page_size)))
	page_info["item"] = import_logs

	count_condition := GetCountCondition{
		FileHash: condition.FileHash,
	}
	count_info, err := GetCountByCondition(count_condition)
	if err != nil {
		return nil, err
	}
	page_info["file_all_count"] = count_info["file_all_count"]
	page_info["success_count"] = count_info["success_count"]

	return page_info, nil
}

func Create(logs *logs.ImportLogs) error {

	if err := db.Get().Create(&logs).Error; err != nil {
		return err
	}
	return nil
}

type GetCountCondition struct {
	FileHash string
}

func GetCountByCondition(condition GetCountCondition) (map[string]interface{}, error) {
	var import_logs []logs.ImportLogs

	var file_count int32
	if err := db.Get().Where("file_hash = (?)", condition.FileHash).Find(&import_logs).Count(&file_count).Error; err != nil {
		return nil, err
	}

	var success_count int32
	if err := db.Get().Where("file_hash = (?) and status = 1", condition.FileHash).Find(&import_logs).Count(&success_count).Error; err != nil {
		return nil, err
	}

	page_info := make(map[string]interface{})
	page_info["file_all_count"] = file_count
	page_info["success_count"] = success_count

	return page_info, nil
}
