package contract

import (
	"context"
	"encoding/json"
	"iQuest/app/constant"
	"iQuest/app/graphql"
	"iQuest/app/graphql/prisma"
	"iQuest/app/model"
	gormModel "iQuest/app/model"
	request "iQuest/app/request/job"
	"iQuest/db"
	"iQuest/library/utils"
	"log"
	"math"
	"time"
)

func ProcessContractJournal(journal model.JobContractJournal) bool {

	if int64(journal.ActiveTimestamp) > time.Now().Unix() {
		return true
	}

	//把历史数据丢弃
	err := db.Get().
		Model(&model.JobContractJournal{}).
		Unscoped().
		Where("company_id = ? and service_company_id = ? and id < ? and active_timestamp > ? ", journal.CompanyId, journal.ServiceCompanyId, journal.ID, journal.ActiveTimestamp).
		UpdateColumn("is_handled", constant.JobJournalHandled).
		Error
	log.Printf("处理合同变更旧数据%v\n", err)
	if err != nil {
		return false
	}

	//现有时间更新到 job_template_delay
	isNeedProcessBegin, isNeedProcessEnd := 0, 0
	if int64(journal.BeginTimestamp) > time.Now().Unix() {
		isNeedProcessBegin = 1
	}
	if int64(journal.EndTimestamp) > time.Now().Unix() {
		isNeedProcessEnd = 1
	}
	if !ProcessJobTemplateDelay(model.JobTemplateDelay{
		ContractNo:         journal.ContractNo,
		CompanyId:          journal.CompanyId,
		ServiceCompanyId:   journal.ServiceCompanyId,
		BeginTimestamp:     journal.BeginTimestamp,
		EndTimestamp:       journal.EndTimestamp,
		IsNeedProcessBegin: isNeedProcessBegin,
		IsNeedProcessEnd:   isNeedProcessEnd,
	}) {
		return false
	}

	/*//处理里边的数据
	var data request.BatchCreateTemplateRequest
	err = json.Unmarshal([]byte(journal.Extend), &data)
	if err != nil {
		log.Printf("反序列化出错%v\n" , err)
		return false
	}

	if len(data.ServiceTypes) > 0{
		for _, serviceType := range data.ServiceTypes {
			for _, template := range serviceType.Templates {
				if template.Id != 0 {
					db.Get().Model(model.JobTemplate{}).Unscoped().Where("id = ? ", template.Id).
				}
			}
		}
	}*/

	err = db.Get().Unscoped().Model(&journal).UpdateColumn("is_handled", constant.JobJournalHandled).Error
	log.Printf("合同变更已处理%v\n", err)
	if err != nil {
		return false
	}

	return true
}

//
func ProcessJobTemplateDelay(delay model.JobTemplateDelay) bool {
	err := db.Get().
		Model(&model.JobTemplateDelay{}).
		Unscoped().
		Where("company_id = ? and service_company_id = ? ", delay.CompanyId, delay.ServiceCompanyId).
		Updates(model.JobTemplateDelay{BeginTimestamp: delay.BeginTimestamp, EndTimestamp: delay.EndTimestamp, IsNeedProcessBegin: delay.IsNeedProcessBegin, IsNeedProcessEnd: delay.IsNeedProcessEnd}).
		Error
	log.Printf("合同时间变更(delay)已处理%v\n", err)
	if err != nil {
		return false
	}

	jobTemplateStatus := constant.JobTemplateEnable
	if int64(delay.BeginTimestamp) > time.Now().Unix() { //已经生效
		//模板生效
		jobTemplateStatus = constant.JobTemplateUnactivated
	}

	if int64(delay.EndTimestamp) <= time.Now().Unix() {
		jobTemplateStatus = constant.JobTemplateExpired
	}

	log.Printf("更新合同为%s的模板\n", delay.ContractNo)

	//更新模板
	var templateIds []int64
	//岗位模板失效后不能再启用
	err = db.Get().Model(&model.JobTemplate{}).Unscoped().Where("company_id = ? and service_company_id = ? and ( is_enable <> ? and is_enable <> ? and is_enable <> ?) ", delay.CompanyId, delay.ServiceCompanyId, constant.JobTemplateExpired, constant.JobTemplateUnAudit, constant.JobTemplateReject).Pluck("id", &templateIds).UpdateColumn("is_enable", jobTemplateStatus).Error
	if err != nil {
		log.Printf("更新岗位模板失败:%s\n", err.Error())
		return false
	}

	//更新岗位
	if len(templateIds) > 0 {
		var workIds []int64
		err = db.Get().Model(&model.Job{}).Unscoped().Where("template_id in (?)", templateIds).Pluck("work_id", &workIds).Error
		if err != nil {
			log.Printf("查询岗位出错:%s\n", err.Error())
			return false
		}

		if len(workIds) > 0 {
			err = db.Get().Model(&model.Work{}).Unscoped().Where("id in (?) and (status <> ? and status <> ?)", workIds, constant.WORK_STATUS_UNAUDIT, constant.WORK_STATUS_Reject).UpdateColumn("status", jobTemplateStatus).Error
			if err != nil {
				log.Printf("更新岗位失败:%s\n", err.Error())
				return false
			}
			log.Printf("更新岗位成功\n")

			//todo 岗位失效后更新岗位成员表状态
			/*if jobTemplateStatus == constant.JobTemplateExpired { // ??????????都变成不生效?????
				err = db.Get().Model(&model.Member{}).Unscoped().Where("work_id in (?)", workIds).UpdateColumn("status", constant.STATUS_TASK_EXPIRED).Error
				if err != nil {
					log.Printf("更新岗位成员表失败:%s\n", err.Error())
					return false
				}
				log.Printf("更新岗位成员表成功\n")
			}*/
		}
	}

	err = db.Get().Unscoped().Model(&model.JobTemplateDelay{}).Where("company_id = ? and service_company_id = ?  ", delay.CompanyId, delay.ServiceCompanyId).Update(&delay).Error
	if err != nil {
		log.Printf("更新合同失败:%s\n", err.Error())
		return false
	}
	log.Printf("更新合同:%s\n", delay.ContractNo)

	return true
}

func ProcessContract(journal model.JobContractJournal) bool {
	var datas request.BatchCreateTemplateRequest
	var newJobTemplates []prisma.JobTemplateCreateInput
	var shuiChouALlJobTemplateIds []int32 //税筹持有的全量岗位模板数据
	//判断appid
	appid := "hgt_appid"

	if utils.IsInXinNiaoCompanyIdArr(datas.CompanyId) {
		appid = constant.XINNIAO_APPID
	}

	if journal.ActiveTimestamp > int(time.Now().Unix()) {
		return false
	}

	err := json.Unmarshal([]byte(journal.Extend), &datas)
	if err != nil {
		return false
	}

	for _, serviceType := range datas.ServiceTypes {
		if len(serviceType.Templates) < 1 {
			continue
		}
		//fmt.Printf("serviceTypeName值%v,地址:%p \n", serviceType.ServiceTypeName, &serviceType.ServiceTypeName)
		serviceTypeName := serviceType.ServiceTypeName

		for _, template := range serviceType.Templates {
			if template.Id != 0 {
				shuiChouALlJobTemplateIds = append(shuiChouALlJobTemplateIds, template.Id)
				//只处理变更的
				continue
			}
			//fmt.Printf("downloadCode值%v,地址:%p \n", template.Attachment.DownloadCode, &template.Attachment.DownloadCode)
			downloadCode := template.Attachment.DownloadCode
			displayName := template.Attachment.DisplayName
			var tmpTemplate prisma.JobTemplateCreateInput
			tmpTemplate.CompanyId = &datas.CompanyId
			tmpTemplate.CompanyName = &datas.CompanyName
			tmpTemplate.ServiceCompanyId = &datas.ServiceCompanyId
			tmpTemplate.ServiceCompanyName = &datas.ServiceCompanyName
			tmpTemplate.ServiceTypeId = serviceType.ServiceTypeId
			tmpTemplate.ServiceTypeName = &serviceTypeName
			tmpTemplate.Name = template.Name
			tmpTemplate.Requirement = template.Requirement
			tmpTemplate.SettlementRule = template.SettlementRule
			tmpTemplate.DownloadCode = &downloadCode
			tmpTemplate.DisplayName = &displayName
			tmpTemplate.Source = constant.JobTemplateSourceAudit
			tmpTemplate.AppId = &appid
			tmpTemplate.WorkType = constant.WorkTypeJob
			tmpTemplate.IsEnable = int32(constant.JobTemplateEnable)
			tmpTemplate.Extend = &journal.Extend
			tmpTemplate.ContractNo = &datas.ContractNo
			tmpTemplate.DeletedAt = prisma.Int32(0)
			newJobTemplates = append(newJobTemplates, tmpTemplate)

		}
	}

	//处理已经拿掉的岗位模板
	var expiredIds []int32
	dbInstance := db.Get().
		Model(&gormModel.JobTemplate{}).
		Unscoped().
		Where("company_id = ? and service_company_id = ?", datas.CompanyId, datas.ServiceCompanyId)
	if len(shuiChouALlJobTemplateIds) > 0 {
		dbInstance = dbInstance.Where("id not in (?)", shuiChouALlJobTemplateIds)
	}
	err = dbInstance.Pluck("id", &expiredIds).
		Pluck("id", &expiredIds).
		UpdateColumn("is_enable", constant.JobTemplateExpired).
		Error

	if err != nil {
		log.Printf("岗位模板失效出错%s:%v \n", err, expiredIds) //TODO 记录重试
		return false
	}

	if len(expiredIds) > 0 {
		var workIds []int32
		err = db.Get().
			Model(&gormModel.Job{}).
			Unscoped().
			Where("template_id in (?) ", expiredIds).
			Pluck("work_id", &workIds).Error

		if len(workIds) > 0 {
			err = db.Get().Model(&gormModel.Work{}).
				Unscoped().
				Where("id in (?) ", workIds).
				UpdateColumn("status", constant.JobTemplateExpired).
				Error

			//岗位失效后更新岗位成员表状态
			err = db.Get().Model(&gormModel.Member{}).Unscoped().Where("work_id in (?)", workIds).
				UpdateColumn("progress", constant.STATUS_TASK_EXPIRED).Error
		}

		if err != nil {
			log.Printf("岗位失效出错%s:%v \n", err, workIds) //TODO 记录重试
			return false
		}
	}

	for _, template := range newJobTemplates {
		tmpTemplate, err := graphql.Server.Prisma.CreateJobTemplate(template).Exec(context.Background())
		if err != nil {
			log.Printf("插入岗位模板失败:%s \n 岗位模板信息:%v \n", err, template) //TODO 记录重试
			return false
		}

		//duration := math.Ceil(time.Now().Sub(time.Unix(int64(endTimestamp), 0)).Hours())
		duration := math.Ceil(time.Unix(int64(journal.EndTimestamp), 0).Sub(time.Now()).Hours())
		extend, _ := json.Marshal(template)
		//创建岗位
		work := &gormModel.Work{
			AppId:			appid,
			CompanyId:     datas.CompanyId,
			ServiceTypeId: template.ServiceTypeId,
			WorkType:      constant.WorkTypeJob,
			Name:          template.Name,
			Requirement:   template.Requirement,
			PayType:       constant.PayTypePayNone,
			Duration:      int(duration), //暂时不做
			EndAt:         journal.EndTimestamp,
			Source:        constant.WorkSourceAppShuiChou,
			Status:        int(template.IsEnable),
			Type:          constant.TypeText,
			IsPublic:      constant.WorkTypePrivate,
			Resume: 		"",
			Extend:        string(extend),
			DeletedAt:     0,
			Job: gormModel.Job{
				Category:        constant.WorkTypeJob,
				Progress:        constant.ProgressOnGoing, //进行中,从常量拿
				SingleRewardMin: 0,
				SingleRewardMax: 0,
				IsNeedProof:     constant.ProofTypeNone,
				Quota:           1,
				TemplateId:      tmpTemplate.ID,
				DeletedAt:       0,
			},
		}

		//如果是薪鸟的则强制公开,并自动勾选简历要求
		if utils.IsInXinNiaoCompanyIdArr(work.CompanyId){
			work.IsPublic = constant.WorkTypePublic
			work.Resume = constant.XINNIAO_RESUME_ID
		}

		err = db.Get().Create(&work).Error
		if err != nil {
			log.Printf("自动生成岗位失败:%s \n 岗位模板信息:%v \n", err, template)
			return false
		}

		//统计数据
		err = db.Get().Create(&gormModel.WorkExtend{
			WorkId: work.ID,
			AppId:  work.AppId,
		}).Error

		if err != nil {
			return false
		}

		//进程
		err = db.Get().Create(&gormModel.WorkProgress{
			AppId:       work.AppId,
			PublisherId: work.UserId,
			WorkId:      work.ID,
			Type:        constant.PROGRESS_STATUS_CREATE,
		}).Error

		if err != nil {
			return false
		}
	}

	journal.IsHandled = constant.JobJournalHandled
	err = db.Get().Unscoped().Save(&journal).Error
	if err != nil {
		log.Printf("更新contract_journal失败:%s", err)
		return false
	}

	return true
}
