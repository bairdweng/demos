package resolver

import (
	"context"
	"errors"
	"iQuest/app/Api"
	"iQuest/app/graphql/model"
	"iQuest/app/graphql/prisma"
	model2 "iQuest/app/model"
	"iQuest/app/model/user"
	"iQuest/config"
	"iQuest/db"
	"strconv"
)

func (r *mutationResolver) SigningAndCreate(ctx context.Context, signData model.SignInput) (*model.SignRspData, error) {
	user := ctx.Value("user").(user.SessionUser)
	user_id := user.UniqueUserId
	//user_id := 1158322647772168192
	//now_time := int32(time.Now().Unix())
	user_id_int, _ := strconv.ParseInt(user_id, 10, 64)
	user_comp, err := Api.FindRealnameInfoByUserid(user_id_int)

	if err != nil {
		return nil, err
	}

	if user_comp.Code != 0 {
		return nil, errors.New(user_comp.Msg)
	}

	if user_comp.Data.CredentialsType != "idcard" || user_comp.Data.CredentialsNo == "" || user_comp.Data.State == 0 {
		return nil, errors.New("缺少身份证号")
	}

	first := int32(1)

	userWechats,err := r.Prisma.UserWeChatAuthorizes(&prisma.UserWeChatAuthorizesParams{Where:&prisma.UserWeChatAuthorizeWhereInput{UserId:&user_id},First:&first}).Exec(ctx)
	if err != nil {
		return nil, errors.New("缺少手机号")
	}

	if len(userWechats) == 0 {
		return nil, errors.New("缺少手机号")
	}

	mobilePhone := userWechats[0].Mobile
	if mobilePhone == ""{
		return nil, errors.New("缺少手机号")
	}


	//if user_comp.Data.MobilePhone == "" {
	//	return nil, errors.New("缺少手机号")
	//
	//}
	work_id := int32(signData.WorkID)
	//获取job
	jobs, err := r.Prisma.Jobs(&prisma.JobsParams{Where: &prisma.JobWhereInput{WorkId: &work_id}, First: &first}).Exec(ctx)
	if err != nil {
		return nil, errors.New("错误")
	}

	if jobs == nil {
		return nil, errors.New("查不到岗位信息")
	}

	job := jobs[0]
	//获取模板
	template, err := r.Prisma.JobTemplate(prisma.JobTemplateWhereUniqueInput{ID: &job.TemplateId}).Exec(ctx)

	if err != nil {
		return nil, errors.New("错误")
	}

	if template == nil {
		return nil, errors.New("查不到岗位信息")
	}
	//获取app_id
	var app_id *string
	if nil == template.PlatformAppid || "" == *template.PlatformAppid {
		app_id, err = Api.GetAppId(strconv.FormatInt(int64(*template.CompanyId), 10))
		if err != nil {
			return nil, err
		}
	} else {
		app_id = template.PlatformAppid

	}

	sign_data := Api.SignSubmitType{}
	sign_data.Identity = user_comp.Data.CredentialsNo
	sign_data.Name = user_comp.Data.RealName
	sign_data.PersonalMobile = mobilePhone
	sign_data.Sign = "aaa"
	sign_data.TemplateId = "hgt"
	sign_data.ServiceCompanyId = strconv.FormatInt(int64(*template.ServiceCompanyId), 10)
	sign_data.IdentityType = "0"
	sign_data.ExtrSystemId = *app_id
	sign_data.CompanyId = strconv.FormatInt(int64(*template.CompanyId), 10)
	sign_data.UserId = user_id
	sign_data.ExtrOrderId = Api.GetExtrOrderId(sign_data)
	//查询用户是否已经签约
	//quer_data := Api.SignQueryType{
	//	ExtrOrderId:  sign_data.ExtrOrderId,
	//	ExtrSystemId: sign_data.ExtrSystemId,
	//	Sign:         "adfa",
	//}
	//
	//query, err := Api.SignQuery(quer_data)
	//if err != nil {
	//	return nil, err
	//}

	companyType := Api.SignQueryByCompanyType{
		ExtrSystemId:     *app_id,
		UserId:           user_id,
		CompanyId:        strconv.FormatInt(int64(*template.CompanyId), 10),
		ServiceCompanyId: strconv.FormatInt(int64(*template.ServiceCompanyId), 10),
		Sign:             "adfa",
	}


	query, err := Api.SignQueryByCompany(companyType)
	if "ACCEPTED" == *query.ResultCode {
		stateCode := [4]string{"AUTHING", "SIGNING", "CLOSED"}
		for i := 0; i < len(stateCode); i++ {
			if *query.State == stateCode[i] {
				return query, nil
			}
		}
	}

	//签约
	signRes, err := Api.SignSubmit(sign_data)
	if err != nil {
		return nil, err
	}

	if "ACCEPTED" != *signRes.ResultCode {
		return signRes, nil

	}

	if signRes.ExtrSystemID != nil && *signRes.ExtrSystemID != ""{
		_, err = r.Prisma.UpdateJobTemplate(prisma.JobTemplateUpdateParams{
			Data:  prisma.JobTemplateUpdateInput{PlatformAppid: signRes.ExtrSystemID},
			Where: prisma.JobTemplateWhereUniqueInput{ID: &template.ID}}).Exec(ctx)
	}

	//更新到常用人员表

	commonUser, err := r.Prisma.CommonlyUsedPersonnels(&prisma.CommonlyUsedPersonnelsParams{Where: &prisma.CommonlyUsedPersonnelWhereInput{UserId: &user_id, CompanyId: template.CompanyId}}).Exec(ctx)
	if err != nil {
		return signRes, err
	}

	bankCard := ""
	for i := 0; i < len(user_comp.Data.AccountList); i++ {
		if (user_comp.Data.AccountList[i].AccountType == "bankCard") {
			bankCard = user_comp.Data.AccountList[i].AccountType
		}
	}
	if commonUser == nil {
		r.Prisma.CreateCommonlyUsedPersonnel(prisma.CommonlyUsedPersonnelCreateInput{
			CompanyId: *template.CompanyId,
			AppId:     "hgt",
			UserId:    user_id,
			CardNo:    sign_data.Identity,
			Education: nil,
			Remark:    nil,
			Mobile:    *signData.Mobile,
			Name:      user_comp.Data.RealName,
			DeletedAt: prisma.Int32(0),
			BankNo:    bankCard,
		}).Exec(ctx)
	}

	return signRes, err

}

func (r *mutationResolver) SignQuery(ctx context.Context, workID int) (*model.SignRspData, error) {
	user := ctx.Value("user").(user.SessionUser)
	user_id := user.UniqueUserId
	//user_id := 1158322647772168192
	//now_time := int32(time.Now().Unix())
	user_id_int, _ := strconv.ParseInt(user_id, 10, 64)
	user_comp, err := Api.FindRealnameInfoByUserid(user_id_int)

	if err != nil {
		return nil, err
	}
	if user_comp.Code != 0 {
		return nil, errors.New(user_comp.Msg)
	}

	if user_comp.Data.CredentialsType != "idcard" || user_comp.Data.CredentialsNo == "" || user_comp.Data.State == 0 {
		return nil, errors.New("缺少身份证号")
	}

	//work_id := workID
	//first := int32(1)

	//获取job
	var job model2.Job

	if err :=db.Get().Model(model2.Job{}).Unscoped().Where("work_id=?",workID).First(&job).Error; err !=nil{
		return nil, errors.New("查不到岗位信息")

	}

	//获取模板
	var template model2.JobTemplate
	if err :=db.Get().Model(model2.JobTemplate{}).Unscoped().Where("id=?",job.TemplateId).First(&template).Error; err !=nil{
		return nil, errors.New("查不到岗位信息")

	}

	companyType := Api.SignQueryByCompanyType{
		//ExtrSystemId:     *app_id,
		UserId:           user_id,
		CompanyId:        strconv.FormatInt(int64(template.CompanyId), 10),
		ServiceCompanyId: strconv.FormatInt(int64(template.ServiceCompanyId), 10),
		Sign:             "adfa",
	}


	query, err := Api.SignQueryByCompany(companyType)
	if err != nil {
		return nil, err
	}

	return query, nil

}

func (r *queryResolver) GetTemplateDownload(ctx context.Context, workID int) (string, error) {


	work_id := int32(workID)
	//获取job
	first := int32(1)
	jobs, err := r.Prisma.Jobs(&prisma.JobsParams{Where: &prisma.JobWhereInput{WorkId: &work_id}, First: &first}).Exec(ctx)
	if err != nil {
		return "", errors.New("错误")
	}

	if jobs == nil {
		return "", errors.New("查不到岗位信息")
	}

	job := jobs[0]
	//获取模板
	template, err := r.Prisma.JobTemplate(prisma.JobTemplateWhereUniqueInput{ID: &job.TemplateId}).Exec(ctx)

	if err != nil {
		return "", errors.New("错误")
	}

	if template == nil {
		return "", errors.New("查不到岗位信息")
	}
	//获取app_id
	var app_id *string
	if nil == template.PlatformAppid || "" == *template.PlatformAppid {
		app_id, err = Api.GetAppId(strconv.FormatInt(int64(*template.CompanyId), 10))
		if err != nil {
			return "", err
		}
	} else {
		app_id = template.PlatformAppid

	}
	url := config.Viper.GetString("Template_Download")
	path := url +"/econtract/extr/ishouru/template/download?extrSystemId="+ *app_id + "&templateId=hgt&serviceCompanyId=" + strconv.FormatInt(int64(*template.ServiceCompanyId), 10)

	return path,nil

}
