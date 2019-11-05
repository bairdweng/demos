package resolver

import (
	"context"
	"encoding/json"
	"errors"
	gormModel "iQuest/app/model"
	"iQuest/app/request/job"
	"iQuest/config"
	"iQuest/db"
	"net/http"
	"strconv"

	"github.com/gogf/gf/encoding/ghtml"
	"github.com/gogf/gf/util/gconv"

	"github.com/99designs/gqlgen/graphql"
	"github.com/fatih/structs"

	"iQuest/app/Api"
	"iQuest/app/constant"
	"iQuest/app/graphql/model"
	"iQuest/app/graphql/prisma"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"

	// "iQuest/app/model/job"
	"iQuest/app/model/user"
	session "iQuest/app/model/user"
	convertTime "iQuest/app/service"
	service "iQuest/app/service/job"
	workService "iQuest/app/service/work"

	//	"iQuest/config"
	"iQuest/library/utils"
	"log"
	"math"

	//	"net/http"

	"time"
	"unsafe"
)

//创建岗位
func (r *mutationResolver) CreateJob(ctx context.Context, data model.NewJobInput) (*model.JobInfo, error) {

	if data.Type != constant.TypeText && data.MediaUrls == nil {
		return nil, errors.New("MediaUrls不能为空")
	}

	var payType int32
	if data.PayType != nil {
		payType = int32(*data.PayType)
	}

	var isNeeProof = constant.ProofTypeNone
	if constant.ProofTypeNone != data.ProofType {
		isNeeProof = 1
		if data.ProofDescription == nil {
			return nil, errors.New("proofDescription不能为空")
		}
	}

	//var	endAt string
	if data.EndAt != nil {
		/*t, err := timeUtils.ParseToPRCTime(*data.EndAt)
		if err != nil {
			return nil, err
		}*/

		if int64(*data.EndAt) < time.Now().Unix() {
			return nil, errors.New("任务结束时间必须大于当前时间")
		}
		//endAt = t.Format(constant.YmdLayout)
	}
	/*
		if data.ContractStartDate > data.ContractEndDate {
			return nil, errors.New("合同有效期不正确")
		}*/

	var mediaUrls string
	if data.Type != constant.TypeText {
		if data.MediaCoverURL == nil || data.MediaUrls == nil || len(data.MediaUrls) < 1 {
			return nil, errors.New("mediaCoverURL和mediaUrls都不能为空不能为空")
		}
		b, _ := json.Marshal(data.MediaUrls)
		mediaUrls = string(b)
	}

	status := constant.WORK_STATUS_UNAUDIT

	user := ctx.Value("user")
	if user == nil {
		return nil, errors.New("无登录信息")
	}

	userId := user.(session.SessionUser).UniqueUserId

	name := data.Name
	requirement := data.Requirement
	settlementRule := data.SettlementRule

	//入库job_contract_Journal表
	contractJournal := gormModel.JobContractJournal{
		ContractNo:       data.ContractNo,
		CompanyId:        int32(data.CompanyID),
		ServiceCompanyId: int32(data.ServiceCompanyID),
		BeginTimestamp:   data.ContractStartDate,
		EndTimestamp:     data.ContractEndDate,
		ActiveTimestamp:  data.ContractStartDate / 1000,
		IsHandled:        1,
	}
	err := db.Get().Unscoped().Where("contract_no = ? and is_handled=?", data.ContractNo, constant.JobJournalHandled).FirstOrCreate(&contractJournal).Error
	if err != nil {
		return nil, errors.New("系统异常:" + err.Error())
	}

	//判断appid
	appIdValue := data.Appid
	if utils.IsInXinNiaoCompanyIdArr(int32(data.CompanyID)) {
		appIdValue = constant.XINNIAO_APPID
	}

	//获取合同信息
	contract := gormModel.JobContractJournal{}
	err = db.Get().Unscoped().Where("contract_no = ? and is_handled = ? and remark = '' ", data.ContractNo, constant.JobJournalHandled).Last(&contract).Error

	if !(int64(contract.BeginTimestamp) <= time.Now().Unix() && int64(contract.EndTimestamp) > time.Now().Unix()) {
		return nil, errors.New("合同不在有效期,不能发岗")
	}

	templateId := 0
	if data.TemplateID != nil {
		//必须要有serviceCompanyId
		status = constant.WORK_STATUS_NORMAL
		templateId = *data.TemplateID
		//TODO 合同生效期才能发岗
		template, err := r.Prisma.JobTemplate(prisma.JobTemplateWhereUniqueInput{
			ID: utils.Int2PointInt32(templateId),
		}).Exec(ctx)
		if err != nil {
			return nil, errors.New("查找岗位模板出错")
		}
		name = template.Name
		requirement = template.Requirement
		settlementRule = template.SettlementRule

	} else {
		extendByte, _ := json.Marshal(data)
		extend := string(extendByte)

		//创建模板
		// t, err := r.Prisma.CreateJobTemplate(prisma.JobTemplateCreateInput{
		// 	AppId:              &appIdValue,
		// 	CompanyId:          utils.Int2PointInt32(data.CompanyID),
		// 	CompanyName:        &data.CompanyName,
		// 	UserId:             &userId,
		// 	WorkType:           constant.WorkTypeJob,
		// 	ServiceTypeId:      int32(data.ServiceTypeID),
		// 	ServiceTypeName:    &data.ServiceTypeNmae,
		// 	ServiceCompanyId:   utils.Int2PointInt32(data.ServiceCompanyID),
		// 	ServiceCompanyName: &data.ServiceCompanyName,
		// 	SettlementRule:     settlementRule,
		// 	Name:               name,
		// 	Requirement:        requirement,
		// 	IsEnable:           int32(constant.JobTemplateUnaided),
		// 	Source:             int32(constant.JobTemplateSourceCompany),
		// 	Extend:             &extend,
		// 	DeletedAt:          prisma.Int32(0),
		// 	ContractNo:         &data.ContractNo,
		// }).Exec(ctx)

		t, err := service.CreateJobTemplate(&gormModel.JobTemplate{
			AppId:              appIdValue,
			CompanyId:          int32(data.CompanyID),
			CompanyName:        data.CompanyName,
			UserId:             userId,
			WorkType:           constant.WorkTypeJob,
			ServiceTypeId:      int32(data.ServiceTypeID),
			ServiceTypeName:    data.ServiceTypeNmae,
			ServiceCompanyId:   int32(data.ServiceCompanyID),
			ServiceCompanyName: data.ServiceCompanyName,
			SettlementRule:     settlementRule,
			Name:               name,
			Requirement:        requirement,
			IsEnable:           constant.JobTemplateUnaided,
			Source:             constant.JobTemplateSourceCompany,
			Extend:             extend,
			DeletedAt:          0,
			ContractNo:         data.ContractNo,
		})
		if err != nil {
			return nil, errors.New("创建模板失败:" + err.Error())
		}

		templateId = int(t.ID)

		//TODO 模板审核请求, 税筹
		ginCtx := ctx.Value("ginContext")
		token := ginCtx.(*gin.Context).GetHeader(config.Viper.GetString("HEADER_AUTH"))

		bizContent := job.BizContent{
			Appid:           appIdValue,
			Name:            data.Name,
			Requirement:     ghtml.StripTags(data.Requirement),
			SettlementRule:  ghtml.StripTags(data.SettlementRule),
			ServiceTypeID:   int32(data.ServiceTypeID),
			ServiceTypeName: data.ServiceTypeNmae,
		}

		bizContentBytes, err := json.Marshal(bizContent)
		if err != nil {
			return nil, err
		}

		resp, err := Api.CreateJobTemplate(job.CreateJobTemplateInput{
			Attach:               appIdValue,
			BizExtendData:        string(bizContentBytes),
			BusinessID:           strconv.Itoa(templateId),
			BusinessType:         Api.BusinessType,
			CallBackUrl:          "http://" + config.Viper.GetString("EUREKA_SERVICE_NAME") + config.Viper.GetString("JOB_TEMPLATE_AUDIT_CALLBACK_URL"),
			CustomerCompanyId:    int32(data.CompanyID),
			CustomerCompanyName:  data.CompanyName,
			ProcessDefinitionKey: Api.ProcessDefinitionKey,
			ProfileId:            int32(data.ProfileID),
			ServiceCompanyId:     int32(data.ServiceCompanyID),
			ServiceCompanyName:   data.ServiceCompanyName,
			UserId:               utils.String2PointInt32(userId),
			UserName:             user.(session.SessionUser).UserName,
			Appid:                appIdValue,
		}, token)

		if err != nil {
			id := int32(t.ID)
			_, _ = r.Prisma.DeleteJobTemplate(prisma.JobTemplateWhereUniqueInput{
				ID: &id,
			}).Exec(ctx)

			return nil, errors.New("启动审核流程失败:" + err.Error())
		}

		if resp != nil && resp.Code != http.StatusOK {
			return nil, errors.New("送审失败")
		}

	}

	var extendStruct prisma.JobTemplateCreateInput
	_ = mapstructure.Decode(data, &extendStruct)
	extendBytes, _ := json.Marshal(extendStruct)
	extendString := string(extendBytes)

	// work, err := r.Prisma.CreateWork(prisma.WorkCreateInput{
	// 	AppId:         appIdValue,
	// 	CompanyId:     int32(data.CompanyID),
	// 	UserId:        user.(session.SessionUser).UniqueUserId,
	// 	ServiceTypeId: int32(data.ServiceTypeID),
	// 	WorkType:      constant.WorkTypeJob,
	// 	Name:          name,
	// 	Requirement:   requirement,
	// 	PayType:       &payType,
	// 	EndAt:         (*int32)(unsafe.Pointer(data.EndAt)),
	// 	Status:        utils.Int2PointInt32(status),
	// 	Extend:        &extendString,
	// 	MediaUrls:     &mediaUrls,
	// 	MediaCoverUrl: data.MediaCoverURL,
	// 	Resume:        data.Resume,
	// 	IsPublic:      utils.Int2PointInt32(data.IsPublic),
	// 	Type:          utils.Int2PointInt32(data.Type), //TODO
	// 	DeletedAt:     prisma.Int32(0),
	// }).Exec(ctx)

	work, err := workService.CreateWork(&gormModel.Work{
		AppId:         appIdValue,
		CompanyId:     int32(data.CompanyID),
		UserId:        user.(session.SessionUser).UniqueUserId,
		ServiceTypeId: int32(data.ServiceTypeID),
		WorkType:      constant.WorkTypeJob,
		Name:          name,
		Requirement:   requirement,
		PayType:       int(payType),
		EndAt:         *data.EndAt,
		Status:        status,
		Extend:        extendString,
		MediaUrls:     mediaUrls,
		MediaCoverUrl: *data.MediaCoverURL,
		Resume:        *data.Resume,
		IsPublic:      data.IsPublic,
		Type:          data.Type, //TODO
		DeletedAt:     0,
	})
	if err != nil {
		return nil, err
	}

	//work progress
	_, err = r.Prisma.CreateWorkProgress(prisma.WorkProgressCreateInput{
		AppId:       appIdValue,
		PublisherId: userId,
		WorkId:      int32(work.ID),
		Type:        constant.PROGRESS_STATUS_CREATE,
	}).Exec(ctx)
	if err != nil {
		log.Printf("任务进程失败:" + err.Error())
	}

	//插入独有数据 job
	job, err := r.Prisma.CreateJob(prisma.JobCreateInput{
		WorkId:           int32(work.ID),
		Category:         constant.WorkTypeJob,
		Progress:         utils.Int2PointInt32(constant.ProgressOnGoing), //进行中,从常量拿
		SingleRewardMin:  &data.SingleRewardMin,
		SingleRewardMax:  &data.SingleRewardMax,
		IsNeedProof:      utils.Int2PointInt32(isNeeProof),
		ProofType:        utils.Int2PointInt32(data.ProofType),
		ProofDescription: data.ProofDescription,
		TemplateId:       int32(templateId),
		DeletedAt:        prisma.Int32(0),
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	//插入用户被邀请的记录 job_member
	if data.InviteIds != nil {
		publisherId := user.(session.SessionUser).UniqueUserId
		source := constant.WORK_SOURCE_INVITE
		progress := constant.STATUS_INVITE
		//now := int32(time.Now().Unix())

		queue := service.QueueConn.OpenQueue(constant.InviteQueue)
		//defer queue.Close()

		for _, id := range data.InviteIds {
			participantId := id
			member, err := r.Prisma.CreateJobMember(prisma.JobMemberCreateInput{
				WorkId:        int32(work.ID),
				PublisherId:   publisherId,
				ParticipantId: &participantId,
				//CommonUseId: &publisherId, //TODO 未知ID
				Source:    utils.Int2PointInt32(source),
				Progress:  utils.Int2PointInt32(progress),
				DeletedAt: prisma.Int32(0),
				//ParticipateAt: &now,
			}).Exec(ctx)

			if err != nil {
				graphql.AddErrorf(ctx, "邀请失败,用户ID: %d", participantId)
			}

			//work progress
			_, err = r.Prisma.CreateWorkProgress(prisma.WorkProgressCreateInput{
				AppId:         appIdValue,
				ParticipantId: &participantId,
				PublisherId:   publisherId,
				WorkId:        int32(work.ID),
				Type:          constant.PROGRESS_STATUS_INVITE,
			}).Exec(ctx)
			if err != nil {
				log.Printf("任务进程失败:" + err.Error())
			}

			payload := map[string]interface{}{
				"userId": id,
				"workId": work.ID,
			}

			payloadBytes, err := json.Marshal(payload)
			if nil != err {
				log.Printf("队列入列消息序列化出错")
			}

			//发送站内信
			sendParam := Api.GetSendMsgContentData{
				WorkId:      int32(work.ID),
				SendId:      userId,
				ReceiverId:  participantId,
				SendType:    3,
				CompanyName: data.CompanyName,
				TaskMemberId: member.ID,
				WorkName:    work.Name,
			}
			go Api.SendMessage(sendParam)

			//发短信邀请 队列
			isSend := queue.PublishBytes(payloadBytes)
			if !isSend {
				//TODO 发送不成功处理?
			}

		}
	}

	//插入统计数据 work_extend
	focusCount := int32(len(data.InviteIds))
	_, err = r.Prisma.CreateWorkExtend(prisma.WorkExtendCreateInput{
		WorkId:     int32(work.ID),
		AppId:      appIdValue, //TODO 从哪拿
		FocusCount: &focusCount,
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	var specify model.Job
	var base model.Work

	if work != nil {
		var mediaUrls []*string
		if work.MediaUrls != "" {
			_ = json.Unmarshal([]byte(work.MediaUrls), &mediaUrls)
		}

		tmpMap := structs.Map(*work)
		tmpMap["CreatedAt"] = int(work.CreatedAt.Unix())
		tmpMap["MediaUrls"] = mediaUrls
		err = mapstructure.Decode(tmpMap, &base)
		if err != nil {
			return nil, err
		}
	}

	if job != nil {
		tmpMap := structs.Map(*job)
		tmpMap["UpdatedAt"] = convertTime.DateTimeToTimestamp(job.UpdatedAt)
		err = mapstructure.Decode(tmpMap, &specify)
		if err != nil {
			return nil, err
		}
	}

	return &model.JobInfo{
		Base:    &base,
		Specify: &specify,
	}, nil
}

func (r *queryResolver) Jobs(ctx context.Context, pageNumber int, pageItem int, search *model.SearchJobInput) (*model.JobPagination, error) {

	user := ctx.Value("user")
	if user == nil {
		return nil, errors.New("无登录信息")
	}

	//userId := user.(session.SessionUser).UserID
	companyId := user.(session.SessionUser).CompanyID

	condition := prisma.WorkWhereInput{
		CompanyId: &companyId,
		DeletedAt: prisma.Int32(0),
	}
	if search != nil {
		if search.ID != nil {
			condition.ID = (*int32)(unsafe.Pointer(search.ID))
		}

		if search.Name != nil {
			condition.NameContains = search.Name
		}

		if search.Status != nil {
			condition.Status = (*int32)(unsafe.Pointer(search.Status))
		}

		if search.ServiceTypeID != nil {
			condition.ServiceTypeId = (*int32)(unsafe.Pointer(search.ServiceTypeID))
		}

		if search.Status != nil {
			condition.Status = (*int32)(unsafe.Pointer(search.Status))
		}

		//TODO 日期验证
		var begin, end string
		var beginTime, endTime time.Time
		//var err error

		if search.Beign != nil {
			//beginTime, err = timeUtils.ParseToPRCTime(*search.Beign)
			beginTime = time.Unix(int64(*search.Beign), 0)
			begin = beginTime.UTC().Format("2006-01-02T15:04:05Z")
		}

		if search.End != nil {
			//endTime, err = timeUtils.ParseToPRCTime(*search.End)
			endTime = time.Unix(int64(*search.End), 0)
			end = endTime.UTC().Format("2006-01-02T15:04:05Z")
		}

		if begin != "" && end != "" {
			if beginTime.After(endTime) {
				return nil, errors.New("结束时间必须大于开始时间")
			}

			timeCondition := []prisma.WorkWhereInput{
				{
					CreatedAtGte: &begin,
				},
				{
					CreatedAtLte: &end,
				},
			}
			condition.And = timeCondition
		} else if begin != "" {
			condition.CreatedAtGte = &begin
		} else if end != "" {
			condition.CreatedAtLte = &end
		}

	}

	skip := (pageNumber - 1) * pageItem
	orderByIdDesc := prisma.WorkOrderByInputIDDesc
	works, err := r.Prisma.Works(&prisma.WorksParams{
		Where:   &condition,
		First:   utils.Int2PointInt32(pageItem),
		Skip:    utils.Int2PointInt32(skip),
		OrderBy: &orderByIdDesc,
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	var jobs []*model.JobInfo
	if works != nil {
		for _, work := range works {
			var job model.JobInfo
			var base model.Work
			var specify model.Job
			var template model.JobTemplate

			tmpMap := structs.Map(work)
			tmpMap["CreatedAt"] = convertTime.DateTimeToTimestamp(work.CreatedAt)
			_ = mapstructure.Decode(tmpMap, &base) //TODO job取得是work id,是否需要取job id
			var mediaUrls []*string
			if work.MediaUrls != nil {
				_ = json.Unmarshal([]byte(*work.MediaUrls), &mediaUrls)
			}

			base.MediaUrls = mediaUrls

			/*
				ret := r.Prisma.Client.GetOne(
					nil,
					prisma.JobWhereUniqueInput{WorkId:&work.ID},
					[2]string{"JobWhereUniqueInput!", "Job"},
					"job",
					[]string{"id","workId","category","settlementRule","payStatus","progress","quota","singleRewardMax","singleRewardMin","isCanComment","isNeedProof","proofDescription","proofType","remark","templateId","extend","createdAt","updatedAt","deletedAt"})
				var s interface{}
				a, err := ret.Exec(ctx, s)

				if nil != err {
					fmt.Println(err)
				}

				fmt.Println(a)
				fmt.Println(ret)
			*/

			jobArr, err := r.Prisma.Jobs(&prisma.JobsParams{
				Where: &prisma.JobWhereInput{
					WorkId:    &work.ID,
					DeletedAt: prisma.Int32(0),
				},
			}).Exec(ctx)

			if err != nil {
				continue
			}

			if len(jobArr) == 0 {
				continue //TODO 理论上1-1对应,不应该出现这种情况,事实上出现了
			}
			_ = mapstructure.Decode(jobArr[0], &specify)

			template_id := jobArr[0].TemplateId
			template_info, err := r.Prisma.JobTemplate(prisma.JobTemplateWhereUniqueInput{
				ID: &template_id,
			}).Exec(ctx)

			if err != nil {
				return nil, err
			}

			_ = mapstructure.Decode(template_info, &template)

			//获取合同信息
			contract := gormModel.JobContractJournal{}
			err = db.Get().Unscoped().Where("contract_no = ? and is_handled = ? and remark = '' ", template.ContractNo, constant.JobJournalHandled).Last(&contract).Error

			if err != nil {
				return nil, errors.New("查找合同信息失败:" + err.Error())
			}
			template.ContractStartDate = &contract.BeginTimestamp
			template.ContractEndDate = &contract.EndTimestamp

			//获取参与人数与未处理人数
			var joinCount, unprocessCount int
			joinCount, unprocessCount = 0, 0

			joinMembers, _ := r.Prisma.JobMembersConnection(&prisma.JobMembersConnectionParams{
				Where: &prisma.JobMemberWhereInput{
					WorkId: &work.ID,
					ProgressIn: constant.JOB_STATUS_ING,
				},
			}).Aggregate(ctx)

			if joinMembers != nil {
				joinCount = int(joinMembers.Count)
			}

			//fmt.Printf(" pointer %p:\n", &joinCount)

			specify.MemberCount = &joinCount

			inviteMembers, err := r.Prisma.JobMembersConnection(&prisma.JobMembersConnectionParams{
				Where: &prisma.JobMemberWhereInput{
					WorkId:   &work.ID,
					Source:   utils.Int2PointInt32(constant.WORK_SOURCE_APPLY),
					Progress: utils.Int2PointInt32(constant.WorkProgressApplying),
				},
			}).Aggregate(ctx)

			if err != nil {
				joinCount = 0
			}

			//上传凭证(申请完成)
			var workProgress []gormModel.WorkProgress
			var applyCompleteCount int
			err = db.Get().Unscoped().Where("id in (?) ", db.Get().Table("work_progress").Unscoped().Select("max(id)").Where("work_id = ? and type = ? ", work.ID, constant.PROGRESS_STATUS_JOB_UPLOAD).Group("participant_id, work_id").QueryExpr()).Find(&workProgress).Count(&applyCompleteCount).Error
			if err != nil {
				log.Printf("查询上传凭证错误:%s\n", err.Error())
			}

			if inviteMembers == nil {
				unprocessCount = 0
			} else {
				unprocessCount = int(inviteMembers.Count)
			}
			unprocessCount += applyCompleteCount

			specify.UnprocessCount = &unprocessCount
			job.Base = &base
			job.Specify = &specify
			job.Template = &template
			jobs = append(jobs, &job)
		}
	}

	var totalItem, totalPage int64
	totalItem, totalPage = 0, 0
	jobAggregate, _ := r.Prisma.WorksConnection(&prisma.WorksConnectionParams{
		Where: &condition,
	}).Aggregate(ctx)

	if jobAggregate != nil {
		totalItem = jobAggregate.Count
	}

	totalPage = int64(math.Ceil(float64(totalItem) / float64(pageItem)))

	return &model.JobPagination{
		PageInfo: &model.PageInfo{
			TotalItem: int(totalItem),
			TotalPage: int(totalPage),
		},
		Items: jobs,
	}, nil
}

//参与列表
func (r *queryResolver) JobMembers(ctx context.Context, ids []int, workID int) ([]*model.JobMember, error) {
	u := ctx.Value("user").(user.SessionUser)

	//获取岗位信息
	w, err := workService.GetByID(workID)
	if err != nil {
		return nil, err
	}

	if u.CompanyID != w.CompanyId {
		return nil, errors.New("权限有误：公司不匹配")
	}

	array, err := service.GetJobMembersByIds(workID, ids)

	members := []*model.JobMember{}
	for _, item := range array {
		createdAt := int(item.CreatedAt.Unix())
		updatedAt := int(item.UpdatedAt.Unix())
		reward := int(item.Reward)
		member := model.JobMember{
			ID:            item.ID,
			WorkID:        &item.WorkId,
			PublisherID:   &item.PublisherId,
			ParticipantID: &item.ParticipantId,
			Source:        &item.Source,
			Progress:      &item.Progress,
			ProofFileURL:  &item.ProofFileUrl,
			Reward:        &reward,
			ParticipateAt: &item.ParticipateAt,
			FinishAt:      &item.FinishAt,
			Extend:        &item.Extend,
			CreatedAt:     &createdAt,
			UpdatedAt:     &updatedAt,
			Remark:        &item.Remark,
		}
		if item.ParticipantUser.ID != 0 {
			gconv.Struct(item.ParticipantUser, &member.ParticipantUser)
		}
		if len(item.JobSettlementLogs) > 0 {
			gconv.Struct(item.JobSettlementLogs[0], &member.LastJobSettlementLog)
			member.LastJobSettlementLog.ID = item.JobSettlementLogs[0].ID
			member.LastJobSettlementLog.CreatedAt = int(item.JobSettlementLogs[0].CreatedAt.Unix())
			member.LastJobSettlementLog.UpdatedAt = int(item.JobSettlementLogs[0].UpdatedAt.Unix())
		}
		members = append(members, &member)
	}
	return members, err
}

/*
func (r *mutationResolver) AddTest(ctx context.Context, data model.AddTInput) (*model.JobT, error) {
 	job, err := r.Prisma.CreateJob(prisma.JobCreateInput{
		WorkId: 1122334455,
		Category: constant.WorkTypeJob,
		Progress: utils.Int2PointInt32(1) , //TODO 进行中,从常量拿
		SingleRewardMin: &data.SingleRewordMin,
		SingleRewardMax: &data.SingleRewordMax,
		IsNeedProof: utils.Int2PointInt32(0),
		ProofType: utils.Int2PointInt32(data.ProofType),
		ProofDescription: data.ProofDescription,
		Work:prisma.WorkCreateOneInput{
			CreateJobTemplateRequest:&prisma.WorkCreateInput{
				AppId:         data.Common.Appid,
				CompanyId:     int32(data.Common.CompanyID),
				//UserId:        user.(session.SessionUser).UserID, //TODO 接口从哪里拿userid
				ServiceTypeId: int32(data.Common.ServiceTypeID),
				WorkType:      constant.WorkTypeJob,
				Name:          data.Common.Name,
				Requirement:   data.Common.Requirement,
				PayType: utils.PointInt2PointInt32(data.Common.PayType, 0),
			},
			Connect:&prisma.WorkWhereUniqueInput{
				ID:utils.Int2PointInt32(22),
			},
		},
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

 	fmt.Println(job)
 	fmt.Println(&job)
 	//fmt.Println(*job)

 	return &model.JobT{
 		ID: int(job.ID),
 		SingleRewordMax: data.SingleRewordMax,
 		SingleRewordMin: data.SingleRewordMin,
 		ProofType: data.ProofType,
 		ProofDescription: data.ProofDescription,
 		Common: &model.Common{
 			Appid: data.Common.Appid,
 			CompanyID: data.Common.CompanyID,
 			ServiceTypeID: 1,
 			Type: data.Common.Type,
 			Name: data.Common.Name,
 			Requirement: data.Common.Requirement,
 			SettlementRule: data.SettlementRule,
		},
	}, nil
}

func (r *queryResolver) JobT(ctx context.Context, id int) (*model.JobT, error) {

	job, err := r.Prisma.Job(prisma.JobWhereUniqueInput{
		ID: utils.Int2PointInt32(id),
	}).Exec(ctx)

	r.Prisma.JobsConnection()

	return &model.JobT{
		ID: int(job.ID),
		SingleRewordMax: *job.SingleRewardMax,
		SingleRewordMin: *job.SingleRewardMin,
		ProofDescription: job.ProofDescription,
		Common: &model.Common{
			Appid: "issdd",
			CompanyID: job.CompanyID,
			ServiceTypeID: 1,
			Type: job.Type,
			Name: job.Name,
			Requirement: job.Requirement,
			SettlementRule: job.SettlementRule,
		},
	}, err
}*/

//更新岗位信息,目前只开放了remark
func (r *mutationResolver) UpdateJob(ctx context.Context, data *model.UpdateJobInput) (*model.Job, error) {

	_, err := r.Prisma.Work(prisma.WorkWhereUniqueInput{
		ID: utils.Int2PointInt32(data.WorkID),
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	user := ctx.Value("user")
	if user == nil {
		return nil, errors.New("无登录信息")
	}

	/*userId := user.(session.SessionUser).UserID

	if work.UserId != utils.Int322String(userId) {
		return nil, errors.New("无权限")
	}*/

	jobArr, err := r.Prisma.Jobs(&prisma.JobsParams{
		Where: &prisma.JobWhereInput{
			WorkId: utils.Int2PointInt32(data.WorkID),
		},
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	if len(jobArr) == 0 {
		return nil, errors.New("找不到对应岗位")
	}

	id := jobArr[0].ID

	resp, err := r.Prisma.UpdateJob(prisma.JobUpdateParams{
		Data: prisma.JobUpdateInput{
			Remark: &data.Remark,
		},
		Where: prisma.JobWhereUniqueInput{
			ID: &id,
		},
	}).Exec(ctx)

	var job model.Job
	tmpMap := structs.Map(resp)
	if _, ok := tmpMap["UpdatedAt"]; ok {
		tmpMap["UpdatedAt"] = convertTime.DateTimeToTimestamp(resp.UpdatedAt)
	}
	err = mapstructure.Decode(tmpMap, &job)
	if err != nil {
		return nil, err
	}

	return &job, nil
}
