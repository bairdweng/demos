package resolver

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/fatih/structs"
	"github.com/mitchellh/mapstructure"
	"iQuest/app/constant"
	"iQuest/app/graphql/model"
	"iQuest/app/graphql/prisma"
	gormModel "iQuest/app/model"
	session "iQuest/app/model/user"
	"iQuest/app/service"
	"iQuest/db"
	"log"
	"math"
	"time"

	//session "iQuest/app/model/user"
	"iQuest/library/utils"
)

//创建岗位模板
func (r *mutationResolver) CreateJobTemplate(ctx context.Context, data model.CreateJobTemplateInput, isAuditPass *bool) (*model.JobTemplate, error) {

	isEnable := 0
	if isAuditPass != nil && *isAuditPass {
		if data.Appid == nil || data.SignTemplateID == nil {
			return nil, errors.New("companyId, appId, signTemplateId 都不能为空")
		}
		isEnable = 1
	}

	user := ctx.Value("user")
	if user == nil {
		return nil, errors.New("无登录信息")
	}

	userId := user.(session.SessionUser).UniqueUserId

	//userId := utils.PointInt2PointInt32(data.UserID, 0) //TODO 接口从哪里拿userid

	//判断appid
	appIdValue := data.Appid
	if utils.IsInXinNiaoCompanyIdArr(int32(data.CompanyID)) {
		temp := constant.XINNIAO_APPID
		appIdValue = &temp
	}

	t, err := r.Prisma.CreateJobTemplate(prisma.JobTemplateCreateInput{
		AppId:              appIdValue,
		CompanyId:          utils.Int2PointInt32(data.CompanyID),
		UserId:            &userId,
		WorkType:           constant.WorkTypeJob,
		SignTemplateId:     utils.PointInt2PointInt32(data.SignTemplateID, 0),
		ServiceTypeId:      int32(data.ServiceTypeID),
		ServiceTypeName:    &data.ServiceTypeName,
		SettlementRule:     data.SettlementRule,
		Name:               data.Name,
		Requirement:        data.Requirement,
		KpiTemplateUrl:     data.KpiTemplateURL,
		IsEnable:           int32(isEnable),
		ServiceCompanyId:   utils.Int2PointInt32(data.ServiceCompanyID),
		ServiceCompanyName: &data.ServiceCompanyName,
		DeletedAt:          prisma.Int32(0),
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	var template model.JobTemplate
	if t != nil {
		tmpMap := structs.Map(*t)
		tmpMap["CreatedAt"] = job.DateTimeToTimestamp(t.CreatedAt)
		err = mapstructure.Decode(tmpMap, &template)
		if err != nil {
			return nil, err
		}
	}

	return &template, nil
}

//审核模板回调 (废弃)
func (r *mutationResolver) AuditJobCallback(ctx context.Context, data model.AuditJobTemplateInput) (*model.JobTemplate, error) {

	t, err := r.Prisma.JobTemplate(prisma.JobTemplateWhereUniqueInput{
		ID: utils.Int2PointInt32(data.JobTemplateID),
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	user := ctx.Value("user")
	if user == nil {
		return nil, errors.New("无登录信息")
	}

	userId := user.(session.SessionUser).UniqueUserId

	if userId != *t.UserId {
		return nil, errors.New("无权限")
	}

	var jobTemplate model.JobTemplate
	isEnable := 0
	if data.IsAuditPass {
		isEnable = 1
	}
	template, err := r.Prisma.UpdateJobTemplate(prisma.JobTemplateUpdateParams{
		Where: prisma.JobTemplateWhereUniqueInput{
			ID: utils.Int2PointInt32(data.JobTemplateID),
		},
		Data: prisma.JobTemplateUpdateInput{
			AppId:              &data.Appid,
			CompanyId:          utils.Int2PointInt32(data.CompanyID),
			SignTemplateId:     utils.Int2PointInt32(data.SignTemplateID),
			ServiceCompanyId:   utils.Int2PointInt32(data.ServiceCompanyID),
			ServiceCompanyName: &data.ServiceCompanyName,
			IsEnable:           utils.Int2PointInt32(isEnable),
		},
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	if template != nil {
		tmpMap := structs.Map(*template)
		tmpMap["CreatedAt"] = job.DateTimeToTimestamp(t.CreatedAt)
		err = mapstructure.Decode(tmpMap, &jobTemplate)
		if err != nil {
			return nil, err
		}
	}
	//TODO 自动发放岗位

	return &jobTemplate, nil
}

//删除模板
func (r *mutationResolver) DeleteJobTemplate(ctx context.Context, id int) (bool, error) {

	template, err := r.Prisma.JobTemplate(prisma.JobTemplateWhereUniqueInput{
		ID: utils.Int2PointInt32(id),
	}).Exec(ctx)

	if err != nil {
		return false, err
	}

	user := ctx.Value("user")
	if user == nil {
		return false, errors.New("无登录信息")
	}

	userId := user.(session.SessionUser).UniqueUserId

	if userId != *template.UserId {
		return false, errors.New("无权限")
	}

	now := int32(time.Now().Unix())
	_, err = r.Prisma.UpdateJobTemplate(prisma.JobTemplateUpdateParams{
		Where: prisma.JobTemplateWhereUniqueInput{
			ID: utils.Int2PointInt32(id),
		},
		Data: prisma.JobTemplateUpdateInput{
			DeletedAt: &now,
		},
	}).Exec(ctx)

	if err != nil {
		return false, err
	}

	return true, nil
}

//岗位模板管理
func (r *queryResolver) JobTemplates(ctx context.Context, pageNumber int, pageItem int, search *model.SearchJobTemplateInput) (*model.JobTemplatePagination, error) {

	condition := prisma.JobTemplateWhereInput{
		//IsEnable: utils.Int2PointInt32(constant.JobTemplateEnable),
	}

	user := ctx.Value("user")
	if user == nil {
		return nil, errors.New("无登录信息")
	}

	//userId := user.(session.SessionUser).UserID
	companyId := user.(session.SessionUser).CompanyID

	//condition.UserId = utils.Int322PointString(userId)
	condition.CompanyId = &companyId
	condition.DeletedAt = prisma.Int32(0)
	if search != nil {
		if search.ID != nil {
			condition.ID = utils.PointInt2PointInt32(search.ID, 0)
		}

		if search.ServiceTypeID != nil {
			condition.ServiceTypeId = utils.PointInt2PointInt32(search.ServiceTypeID, 0)
		}

		if search.Name != nil {
			condition.NameContains = search.Name
		}
	}

	//condition.DeletedAtNot = nil

	skip := (pageNumber - 1) * pageItem

	orderBy := prisma.JobTemplateOrderByInputIDDesc
	templates, err := r.Prisma.JobTemplates(&prisma.JobTemplatesParams{
		Where:   &condition,
		First:   utils.Int2PointInt32(pageItem),
		Skip:    utils.Int2PointInt32(skip),
		OrderBy: &orderBy,
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	var jobTemplates []*model.JobTemplate
	for _, item := range templates {
		var jobTemplate model.JobTemplate
		tmpMap := structs.Map(item)
		tmpMap["CreatedAt"] = job.DateTimeToTimestamp(item.CreatedAt)
		err = mapstructure.Decode(tmpMap, &jobTemplate)
		if err != nil {
			return nil, errors.New("模板转换失败:" + err.Error())
		}

		//获取合同信息
		contract := gormModel.JobContractJournal{}
		err = db.Get().Unscoped().Where("contract_no = ? and is_handled = ? and remark = '' ", jobTemplate.ContractNo, constant.JobJournalHandled).Last(&contract).Error
		if err != nil {
			return nil, errors.New("查找合同信息失败:" + err.Error())
		}
		jobTemplate.ContractStartDate = &contract.BeginTimestamp
		jobTemplate.ContractEndDate = &contract.EndTimestamp


		jobTemplates = append(jobTemplates, &jobTemplate)
	}

	var totalItem, totalPage int
	totalItem, totalPage = 0, 0
	jobTemplateAggregate, _ := r.Prisma.JobTemplatesConnection(&prisma.JobTemplatesConnectionParams{
		Where: &condition,
	}).Aggregate(ctx)

	if jobTemplateAggregate != nil {
		totalItem = int(jobTemplateAggregate.Count)
	}

	totalPage = int(math.Ceil(float64(totalItem) / float64(pageItem)))

	return &model.JobTemplatePagination{
		PageInfo: &model.PageInfo{
			TotalItem: totalItem,
			TotalPage: totalPage,
		},
		Items: jobTemplates,
	}, nil

}

func (r *queryResolver) JobTemplate(ctx context.Context, id int) (*model.JobTemplateInfo, error) {

	template, err := r.Prisma.JobTemplate(prisma.JobTemplateWhereUniqueInput{
		ID: utils.Int2PointInt32(id),
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	tmpMap := structs.Map(template)
	if _, ok := tmpMap["CreatedAt"]; ok {
		tmpMap["CreatedAt"] = job.DateTimeToTimestamp(template.CreatedAt)
	}
	var jobTemplate model.JobTemplate
	err = mapstructure.Decode(tmpMap, &jobTemplate)

	if err != nil {
		log.Printf("转换错误:%v\n", err)
		return nil, errors.New("系统错误")
	}

	//获取合同信息
	contract := gormModel.JobContractJournal{}
	err = db.Get().Unscoped().Where("contract_no = ? and is_handled = ? and remark = '' ", jobTemplate.ContractNo, constant.JobJournalHandled).Last(&contract).Error

	if err != nil {
		return nil, errors.New("查找岗位失败:" + err.Error())
	}
	jobTemplate.ContractStartDate = &contract.BeginTimestamp
	jobTemplate.ContractEndDate = &contract.EndTimestamp


	var specify model.Job
	var base model.Work
	quota := 0

	var media_urls []*string

	job_detail, err := r.Prisma.Jobs(&prisma.JobsParams{
		Where: &prisma.JobWhereInput{
			TemplateId: &template.ID,
		},
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	if job_detail != nil && len(job_detail) > 0 {

		if job_detail[0].Quota == 0 {
			quota = int(job_detail[0].Quota)
		}

		if job_detail[0].SingleRewardMin == nil {
			f := float64(0)
			job_detail[0].SingleRewardMin = &f
		}
		if job_detail[0].SingleRewardMax == nil {
			f := float64(0)
			job_detail[0].SingleRewardMax = &f
		}
		jSingleRewardMin := job_detail[0].SingleRewardMin
		jSingleRewardMax := job_detail[0].SingleRewardMax
		proof_type := int(*job_detail[0].ProofType)

		specify.WorkID = int(job_detail[0].WorkId)
		specify.Quota = &quota
		specify.SingleRewardMin = *jSingleRewardMin
		specify.SingleRewardMax = *jSingleRewardMax
		specify.IsNeedProof = int(job_detail[0].IsNeedProof)
		specify.ProofDescription = job_detail[0].ProofDescription
		specify.ProofType = &proof_type
		specify.Remark = job_detail[0].Remark

		work, err := r.Prisma.Work(prisma.WorkWhereUniqueInput{
			ID: &job_detail[0].WorkId,
		}).Exec(ctx)

		if err != nil {
			return nil, err
		}

		c := []byte(*work.MediaUrls)
		_ = json.Unmarshal(c, &media_urls)

		end_at := int(*work.EndAt)

		var mediaUrls []*string
		if work.MediaUrls != nil {
			_ = json.Unmarshal([]byte(*work.MediaUrls), &mediaUrls)
		}

		base.EndAt = &end_at
		base.Requirement = work.Requirement
		base.MediaCoverURL = work.MediaCoverUrl
		base.MediaUrls = mediaUrls
		base.Type = utils.Int322PointInt(work.Type)
		base.Duration = utils.Int322PointInt(*work.Duration)

	}

	return &model.JobTemplateInfo{
		Base: &base,
		JobDetail: &specify,
		Template:  &jobTemplate,
		MediaUrls: media_urls,
	}, nil
}
