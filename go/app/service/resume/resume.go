package resume

import (
	"iQuest/app/graphql/model"
	"iQuest/app/model/resume"
	"iQuest/db"
)

//type CategoryType struct {
//	Id int `json:"id"`
//	Name string `json:"name"`
//	Resume []ResumeType `json:"resume"`
//}
//
//type ResumeType struct {
//	Id int `json:"id"`
//	Name string `json:"name"`
//	IsHas int `json:"is_has"`
//}

func GetResumeByCategory(categoryId int, belongID int, userId string) ([]model.ResumeType, error) {
	var resumes []resume.Resume
	var userResumes []resume.UserResume
	var resumeTypes []model.ResumeType
	var resumeType model.ResumeType

	belongIDList := []int{0, belongID}

	if err := db.Get().Unscoped().Where("category_id = (?) and belong_to_id IN (?) and status = 1", categoryId, belongIDList).Find(&resumes).Error; err != nil {
		return nil, err
	}

	if err := db.Get().Unscoped().Where("user_id = (?) and category_id =(?) and status = 1", userId, categoryId).Find(&userResumes).Error; err != nil {
		//return nil, err
	}

	resumesCount := len(resumes)
	userResumesCount := len(userResumes)

	for i := 0; i < resumesCount; i++ {
		resumeType = model.ResumeType{
			ID:    resumes[i].ID,
			Name:  resumes[i].Name,
			IsHas: 0,
		}

		for j := 0; j < userResumesCount; j++ {
			if resumes[i].ID == userResumes[j].ResumeId {
				resumeType.IsHas = 1

			}
		}
		if resumeType.IsHas == 0 {
			resumeType.IsHas = resumes[i].IsDefault
		}
		resumeTypes = append(resumeTypes, resumeType)

	}

	return resumeTypes, nil

}

func GetResumeCategory(categoryIdList []int, userId string) ([]model.CategoryType, error) {
	var categorys []resume.ResumeCategory
	var userResumes []resume.UserResume
	var categoryTypes []model.CategoryType
	//var categoryType model.CategoryType
	//var resumeType model.ResumeType
	if err := db.Get().Unscoped().Where("id In (?) and status = 1", categoryIdList).Find(&categorys).Error; err != nil {
		return nil, err
	}

	if err := db.Get().Unscoped().Where("user_id = (?) and category_id In (?) and status = 1", userId, categoryIdList).Find(&userResumes).Error; err != nil {
		//return nil, err
	}
	categorysCount := len(categorys)
	userResumesCount := len(userResumes)

	for i := 0; i < categorysCount; i++ {
		categoryType := model.CategoryType{
			ID:   categorys[i].ID,
			Name: categorys[i].Name,
			Type: categorys[i].Type,
		}
		for j := 0; j < userResumesCount; j++ {
			if categorys[i].ID == userResumes[j].CategoryId {
				resumeType := model.ResumeType{
					ID:    userResumes[j].ResumeId,
					Name:  userResumes[j].ResumeName,
					IsHas: 1,
				}

				categoryType.Resume = append(categoryType.Resume, &resumeType)
			}

		}

		//设置默认值
		if len(categoryType.Resume) == 0 {
			var resumesDef resume.Resume

			db.Get().Unscoped().Where(" category_id = (?) and is_default = 1 and status = 1 ", categorys[i].ID).First(&resumesDef)
			if resumesDef.ID != 0 {
				resumeType := model.ResumeType{
					ID:    resumesDef.ID,
					Name:  resumesDef.Name,
					IsHas: 1,
				}
				categoryType.Resume = append(categoryType.Resume, &resumeType)

			}

		}

		categoryTypes = append(categoryTypes, categoryType)

	}
	return categoryTypes, nil

}

func CreateUserResume(categoryTypeList []model.CategoryInput, userId string, workId int) (bool, error) {
	//var resumeList []resume.Resume

	//var userResumeList []resume.UserResume
	var userResume resume.UserResume
	var workUserResume resume.WorkUserResume
	//var resumeTypes []model.ResumeType

	//db.Unscoped().Delete(&order)

	tx := db.Get().Begin()
	categoryTypeListCount := len(categoryTypeList)
	for i := 0; i < categoryTypeListCount; i++ {

		if err := db.Get().Unscoped().Where("category_id = (?) and user_id = (?)", categoryTypeList[i].ID, userId).Delete(resume.UserResume{}).Error; err != nil {
			tx.Rollback()
			return false, err
		}

		if err := db.Get().Unscoped().Where("category_id = (?) and user_id = (?) and work_id = (?)", categoryTypeList[i].ID, userId, workId).Delete(resume.WorkUserResume{}).Error; err != nil {
			tx.Rollback()
			return false, err
		}
		resumeList := categoryTypeList[i].Resume
		resumeCount := len(resumeList)
		for j := 0; j < resumeCount; j++ {
			userResume = resume.UserResume{
				UserID:       userId,
				ResumeId:     resumeList[j].ID,
				ResumeName:   resumeList[j].Name,
				CategoryId:   categoryTypeList[i].ID,
				CategoryName: categoryTypeList[i].Name,
			}
			if err := db.Get().Create(&userResume).Error; err != nil {
				tx.Rollback()
				return false, err
			}

			workUserResume = resume.WorkUserResume{
				UserID:       userId,
				WorkId:       workId,
				ResumeId:     resumeList[j].ID,
				ResumeName:   resumeList[j].Name,
				CategoryId:   categoryTypeList[i].ID,
				CategoryName: categoryTypeList[i].Name,
			}
			if err := db.Get().Create(&workUserResume).Error; err != nil {
				tx.Rollback()
				return false, err
			}
		}

	}

	tx.Commit()
	return true, nil

}

func GetWorkUserResume(workId int, userId string) ([]model.CategoryType, error) {
	var workUserResumeList []resume.WorkUserResume
	var categoryTypeList []model.CategoryType

	if err := db.Get().Unscoped().Where("work_id = (?) and user_id = (?)", workId, userId).Find(&workUserResumeList).Error; err != nil {
		return nil, err
	}

	count := len(workUserResumeList)

	for i := 0; i < count; i++ {
		isContinue := false
		for k := 0; k < len(categoryTypeList); k++ {
			if categoryTypeList[k].ID == workUserResumeList[i].CategoryId {
				isContinue = true
			}
		}
		if isContinue {
			continue
		}

		categoryType := model.CategoryType{
			ID:   workUserResumeList[i].CategoryId,
			Name: workUserResumeList[i].CategoryName,
			Type: 1,
		}
		for j := 0; j < count; j++ {
			if categoryType.ID == workUserResumeList[j].CategoryId {

				resumeType := model.ResumeType{
					ID:    workUserResumeList[j].ResumeId,
					Name:  workUserResumeList[j].ResumeName,
					IsHas: 1,
				}
				categoryType.Resume = append(categoryType.Resume, &resumeType)
			}

		}
		categoryTypeList = append(categoryTypeList, categoryType)
	}

	return categoryTypeList, nil
}
