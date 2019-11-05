package resolver

import (
	"context"
	"encoding/json"
	"errors"
	"iQuest/app/graphql/model"
	model2 "iQuest/app/model"
	resume2 "iQuest/app/model/resume"
	"iQuest/app/model/user"
	"iQuest/app/service/resume"
	"iQuest/db"
)

func (r *queryResolver) GetResumeCategory(ctx context.Context, workID int) ([]*model.CategoryType, error) {

	// user := ctx.Value("user").(user.SessionUser)
	// userId := user.UniqueUserId

	// if userId == "" {
	// 	return nil, errors.New("请先进行实名认证")
	// }
	// workID 1132 有数据嘛。
	userId := "123456"

	var work model2.Work
	if err := db.Get().Unscoped().First(&work, workID).Error; err != nil {
		return nil, err
	}

	type categoryJson struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	}

	var categoryJsonList []categoryJson
	categoryStr := work.Resume
	if categoryStr == "" {
		return nil, errors.New("该岗位未配置简历！")
	}

	json.Unmarshal([]byte(categoryStr), &categoryJsonList)
	var categoryIdList []int
	for i := 0; i < len(categoryJsonList); i++ {
		categoryIdList = append(categoryIdList, categoryJsonList[i].Id)
	}
	//userId := "123456"

	//categoryIdList := []int{1, 2, 3, 4, 5, 6}
	categoryList, err := resume.GetResumeCategory(categoryIdList, userId)
	if err != nil {
		return nil, err
	}

	res := []*model.CategoryType{}

	categorysCount := len(categoryList)
	for i := 0; i < categorysCount; i++ {
		res = append(res, &categoryList[i])
	}

	return res, nil
}

func (r *queryResolver) GetResumeByCategory(ctx context.Context, categoryID int, belongID int) ([]*model.ResumeType, error) {
	user := ctx.Value("user").(user.SessionUser)
	userId := user.UniqueUserId

	if userId == "" {
		return nil, errors.New("请先进行实名认证")
	}

	//userId := "123456"
	resumeList, err := resume.GetResumeByCategory(categoryID, belongID, userId)
	if err != nil {
		return nil, err
	}

	res := []*model.ResumeType{}

	categorysCount := len(resumeList)
	for i := 0; i < categorysCount; i++ {
		res = append(res, &resumeList[i])
	}

	return res, nil
}
func (r *mutationResolver) CreateUserResume(ctx context.Context, categoryInput []*model.CategoryInput, workID int) (*bool, error) {
	user := ctx.Value("user").(user.SessionUser)
	userId := user.UniqueUserId

	if userId == "" {
		return nil, errors.New("请先进行实名认证")
	}

	//userId := "123456"
	var categoryInputList []model.CategoryInput
	count := len(categoryInput)
	for i := 0; i < count; i++ {
		categoryInputList = append(categoryInputList, *categoryInput[i])
	}
	res, err := resume.CreateUserResume(categoryInputList, userId, workID)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *queryResolver) GetWorkUserResume(ctx context.Context, workID int, userID string) ([]*model.CategoryType, error) {
	userId := userID
	categoryTypeList, err := resume.GetWorkUserResume(workID, userId)
	if err != nil {
		return nil, err
	}

	res := []*model.CategoryType{}

	count := len(categoryTypeList)
	for i := 0; i < count; i++ {
		res = append(res, &categoryTypeList[i])
	}

	return res, nil
}

func (r *queryResolver) IsNeedResume(ctx context.Context, workID int) (*bool, error) {
	user := ctx.Value("user").(user.SessionUser)
	userId := user.UniqueUserId
	res := false
	if userId == "" {

		return &res, nil
	}

	var work model2.Work
	if err := db.Get().Unscoped().First(&work, workID).Error; err != nil {
		return nil, err
	}
	if work.Resume == "" {
		return &res, nil
	}
	var count int
	if err := db.Get().Unscoped().Model(resume2.WorkUserResume{}).Where("user_id = ? and work_id = ?", userId, workID).Count(&count).Error; err != nil {
		return nil, err
	}

	if count == 0 {
		res = true
	}

	return &res, nil

}
