package resolver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"iQuest/app/Api"
	"iQuest/app/constant"
	"iQuest/app/graphql/model"
	"iQuest/app/graphql/prisma"
	gormModel "iQuest/app/model"
	session "iQuest/app/model/user"
	JobService "iQuest/app/service"
	"iQuest/db"
	"iQuest/library/utils"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/fatih/structs"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"github.com/speps/go-hashids"
	"github.com/xlstudio/wxbizdatacrypt"
)

//邀请参加
func (r *mutationResolver) Invite(ctx context.Context, data model.InviteJoinInput) (*model.JobMember, error) {
	user := ctx.Value("user")
	if user == nil {
		return nil, errors.New("无登录信息")
	}
	sendUserId := strconv.Itoa(int(user.(session.SessionUser).UserID))
	userId := string(data.UserID)

	workId := int32(data.WorkID)

	work, err := r.Prisma.Work(prisma.WorkWhereUniqueInput{ID: &workId}).Exec(ctx)
	if work == nil {
		return nil, errors.New("数据出错")
	}

	workTitle := utils.GetWorkTitle(work.CompanyId)

	//权限判断
	companyID := int32(user.(session.SessionUser).CompanyID)
	if work.CompanyId != companyID {
		return nil, errors.New("权限不足")
	}
	if work.EndAt != nil {
		//招募截止为空则无限招募时间 不为空判断招募时间
		if int32(time.Now().Unix()) > int32(*work.EndAt) {
			//招募时间已截止
			return nil, errors.New(workTitle + "不允许邀请")
		}
	}

	if *work.Status != int32(constant.WORK_STATUS_NORMAL) {
		//任务状态不正常
		return nil, errors.New(workTitle + "不允许邀请")
	}
	//1: 邀请中，2已申请，3已同意参加，4申请完成，5已同意完成，6已评分，7发布者已评分，8互评, 9拒绝申请, 10拒绝完成, 20踢出任务
	//判断用户是否已在流程中
	jobData, err := r.Prisma.JobMembers(&prisma.JobMembersParams{
		Where: &prisma.JobMemberWhereInput{
			ProgressIn:    constant.JOB_STATUS_CAN_APPLY,
			ParticipantId: prisma.Str(userId),
			WorkId:        &workId,
		},
	}).Exec(ctx)

	if jobData != nil {
		//发送消息
		type T struct {
			CompanyName string `json:"companyName"`
		}
		var t1 T
		err = json.Unmarshal([]byte(*work.Extend), &t1)
		newSendParam := Api.GetSendMsgContentData{
			WorkId:       workId,
			SendId:       sendUserId,
			ReceiverId:   userId,
			SendType:     3,
			CompanyName:  t1.CompanyName,
			TaskMemberId: jobData[0].ID,
			WorkName:     work.Name,
		}
		_, _ = Api.SendMessage(newSendParam)
		t0 := model.JobMember{ID: int(jobData[0].ID)}
		return &t0, err
	}

	sourceInvite := int32(constant.WORK_SOURCE_INVITE)
	member, err := r.Prisma.CreateJobMember(prisma.JobMemberCreateInput{
		ParticipantId: &userId,
		PublisherId:   work.UserId,
		WorkId:        workId,
		CompanyId:     companyID,
		Source:        &sourceInvite,
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	t := model.JobMember{ID: int(member.ID)}

	//入库任务进程表
	statusApplyProcess := string(constant.PROGRESS_STATUS_INVITE)
	_, _ = r.Prisma.CreateWorkProgress(prisma.WorkProgressCreateInput{
		AppId:         work.AppId,
		PublisherId:   work.UserId,
		ParticipantId: &userId,
		WorkId:        workId,
		Type:          statusApplyProcess,
	}).Exec(ctx)

	//发送消息
	type T struct {
		CompanyName string `json:"companyName"`
	}
	var t1 T
	err = json.Unmarshal([]byte(*work.Extend), &t1)
	sendParam := Api.GetSendMsgContentData{
		WorkId:       workId,
		SendId:       sendUserId,
		ReceiverId:   userId,
		SendType:     3,
		CompanyName:  t1.CompanyName,
		TaskMemberId: member.ID,
		WorkName:     work.Name,
	}
	_, _ = Api.SendMessage(sendParam)
	return &t, err
}

//申请参加任务
func (r *mutationResolver) Apply(ctx context.Context, workID int) (*model.JobMember, error) {
	user := ctx.Value("user")
	if user == nil {
		return nil, errors.New("无登录信息")
	}
	userId := user.(session.SessionUser).UniqueUserId

	workId := int32(workID)
	work, err := r.Prisma.Work(prisma.WorkWhereUniqueInput{ID: &workId}).Exec(ctx)
	if work == nil || err != nil {
		return nil, errors.New("数据出错")
	}
	workTitle := utils.GetWorkTitle(work.CompanyId)

	companyId := work.CompanyId
	//申请任务间隔时间判断
	key := "apply:" + fmt.Sprint("%d", userId) + "-" + fmt.Sprint("%d", workId)
	value, _ := db.Redis().Get(key).Result()
	if value != "" {
		return nil, errors.New("您已成功申请" + workTitle + "，请等待对方确认")
	}

	//不能申请自己的任务
	if userId == work.UserId {
		return nil, errors.New("不能申请自己的" + workTitle)
	}

	//判断是否已经拉黑
	blackData, err := r.Prisma.CompanyUserBlacklistsConnection(&prisma.CompanyUserBlacklistsConnectionParams{
		Where: &prisma.CompanyUserBlacklistWhereInput{
			WorkId:        &workId,
			ParticipantId: &userId,
		},
	}).Aggregate(ctx)

	if err != nil || blackData.Count > 0 {
		return nil, errors.New("您已被禁止参与该" + workTitle)
	}

	//是否还有进行中的任务
	jobData, err := r.Prisma.JobMembers(&prisma.JobMembersParams{
		Where: &prisma.JobMemberWhereInput{
			ProgressIn:    []int32{2, 3, 4, 10},
			ParticipantId: &userId,
			WorkId:        &workId,
		},
	}).Exec(ctx)

	if jobData != nil {
		return nil, errors.New("已申请或有进行中" + workTitle + ",不能再次申请")
	}

	//任务招募时间是否有效
	if work.EndAt != nil {
		//招募截止为空则无限招募时间 不为空判断招募时间
		if int32(time.Now().Unix()) > int32(*work.EndAt) {
			//招募时间已截止
			return nil, errors.New(workTitle + "不允许参加")
		}
	}

	//任务状态是否正确
	if *work.Status != int32(constant.WORK_STATUS_NORMAL) ||
		work.IsPublic != int32(constant.WORK_STATUS_PUBLIC) {
		//任务状态不正常
		return nil, errors.New(workTitle + "不允许参加")
	}

	//TODO 事务

	//查询关注人数
	oldExtend, _ := r.Prisma.WorkExtends(&prisma.WorkExtendsParams{
		Where: &prisma.WorkExtendWhereInput{
			WorkId: &workId,
		},
	}).Exec(ctx)
	//设置关注人数
	focusCount := JobService.GetFocusCount(constant.TASK_FOCUS_NUM_APPLY)
	newCount := oldExtend[0].FocusCount + focusCount
	_, _ = r.Prisma.UpdateWorkExtend(prisma.WorkExtendUpdateParams{
		Data: prisma.WorkExtendUpdateInput{
			FocusCount: &newCount,
		},
		Where: prisma.WorkExtendWhereUniqueInput{
			ID: &workId,
		},
	}).Exec(ctx)

	//有邀请就参加任务
	inviteStatus := int32(1)
	pageLimit := int32(1)
	orderBy := prisma.JobMemberOrderByInputIDDesc
	isInviteData, err := r.Prisma.JobMembers(&prisma.JobMembersParams{
		Where: &prisma.JobMemberWhereInput{
			Progress:      &inviteStatus,
			ParticipantId: &userId,
			WorkId:        &workId,
		},
		OrderBy: &orderBy,
		First:   &pageLimit,
	}).Exec(ctx)

	if isInviteData != nil {
		statusApprove := int32(constant.STATUS_APPROVE)
		jobMember, err := r.Prisma.UpdateJobMember(prisma.JobMemberUpdateParams{
			Data: prisma.JobMemberUpdateInput{
				Progress: &statusApprove,
			},
			Where: prisma.JobMemberWhereUniqueInput{
				ID: &isInviteData[0].ID,
			},
		}).Exec(ctx)
		if jobMember == nil || err != nil {
			return nil, errors.New(workTitle + "参加失败")
		}
		t := model.JobMember{ID: int(jobMember.ID)}
		return &t, err
	}

	//参加任务
	statusApply := int32(constant.STATUS_APPLY)
	member, err := r.Prisma.CreateJobMember(prisma.JobMemberCreateInput{
		ParticipantId: &userId,
		PublisherId:   work.UserId,
		WorkId:        workId,
		CompanyId:     companyId,
		Progress:      &statusApply,
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	//入库任务进程表
	statusApplyProcess := string(constant.PROGRESS_STATUS_APPLY)
	_, _ = r.Prisma.CreateWorkProgress(prisma.WorkProgressCreateInput{
		AppId:         work.AppId,
		PublisherId:   work.UserId,
		ParticipantId: &userId,
		WorkId:        workId,
		Type:          statusApplyProcess,
	}).Exec(ctx)

	//设置redis
	db.Redis().Set(key, 1, time.Hour)

	//查询常用人员表
	user_id, _ := r.Prisma.CommonlyUsedPersonnels(&prisma.CommonlyUsedPersonnelsParams{
		Where: &prisma.CommonlyUsedPersonnelWhereInput{
			UserId:    &userId,
			CompanyId: &companyId,
			DeletedAt: prisma.Int32(0),
		},
	}).Exec(ctx)

	if user_id == nil {
		//查询用户中心用户信息
		intUserId, _ := strconv.Atoi(userId)
		ver, err := Api.FindRealnameInfoByUserid(int64(intUserId))

		if err == nil && ver.Data.State == 1 {
			//判断手机号,为空则查询用户微信授权表
			mobilePhone := ver.Data.MobilePhone
			if mobilePhone == "" {
				var userWeChat gormModel.UserWeChatAuthorize
				db.Get().Unscoped().Where("user_id = ?", userId).First(&userWeChat)
				if userWeChat.Mobile != "" {
					mobilePhone = userWeChat.Mobile
				}
			}
			cardNo := ""
			if ver.Data.CredentialsType == "idcard" {
				cardNo = ver.Data.CredentialsNo
			}

			//入库常用人员表
			_, _ = r.Prisma.CreateCommonlyUsedPersonnel(prisma.CommonlyUsedPersonnelCreateInput{
				AppId:     work.AppId,
				CompanyId: companyId,
				UserId:    userId,
				Name:      ver.Data.RealName,
				Mobile:    mobilePhone,
				CardNo:    cardNo,
				DeletedAt: prisma.Int32(0),
			}).Exec(ctx)
		}
	}

	v := model.JobMember{ID: int(member.ID)}
	return &v, err
}

//接受邀请
func (r *mutationResolver) ApproveInvite(ctx context.Context, memberID int) (*model.JobMember, error) {
	user := ctx.Value("user")
	if user == nil {
		return nil, errors.New("无登录信息")
	}
	userID := user.(session.SessionUser).UniqueUserId

	memberId := int32(memberID)

	//判断任务是否有效
	jobData, err := r.Prisma.JobMember(prisma.JobMemberWhereUniqueInput{
		ID: &memberId,
	}).Exec(ctx)

	if err != nil || jobData == nil || jobData.Progress != constant.STATUS_INVITE {
		return nil, errors.New("数据错误")
	}

	//判断当前用户是否是参与者
	participantId := *jobData.ParticipantId
	if participantId != userID {
		return nil, errors.New("权限不足")
	}

	workId := int32(jobData.WorkId)
	work, err := r.Prisma.Work(prisma.WorkWhereUniqueInput{ID: &workId}).Exec(ctx)
	if work == nil {
		return nil, errors.New("数据出错")
	}

	workTitle := utils.GetWorkTitle(work.CompanyId)

	if work.EndAt != nil {
		//招募截止为空则无限招募时间 不为空判断招募时间
		if int32(time.Now().Unix()) > int32(*work.EndAt) {
			//招募时间已截止
			return nil, errors.New(workTitle + "不允许参加")
		}
	}

	if *work.Status != int32(constant.WORK_STATUS_NORMAL) {
		//任务状态不正常
		return nil, errors.New(workTitle + "不允许参加")
	}

	//可以参加任务
	statusApprove := int32(constant.STATUS_APPROVE)
	nowTimestamp := int32(time.Now().Unix())
	jobMember, err := r.Prisma.UpdateJobMember(prisma.JobMemberUpdateParams{
		Data: prisma.JobMemberUpdateInput{
			Progress:      &statusApprove,
			ParticipateAt: &nowTimestamp,
		},
		Where: prisma.JobMemberWhereUniqueInput{
			ID: &jobData.ID,
		},
	}).Exec(ctx)
	if jobMember == nil || err != nil {
		//return nil, errors.New("任务参加失败")
		return nil, err
	}

	//入库任务进程表
	statusApplyProcess := string(constant.PROGRESS_STATUS_APPROVE_INVITE)
	_, _ = r.Prisma.CreateWorkProgress(prisma.WorkProgressCreateInput{
		AppId:         work.AppId,
		PublisherId:   work.UserId,
		ParticipantId: &userID,
		WorkId:        workId,
		Type:          statusApplyProcess,
	}).Exec(ctx)

	//任务参加回调
	//sUserID, _ := strconv.Atoi(userID)
	//_, _ = Api.UserJoinJobCallBack(int64(sUserID), workId)
	t := model.JobMember{ID: int(jobMember.ID)}
	return &t, err
}

//同意参加
func (r *mutationResolver) Approve(ctx context.Context, data model.JobInput) (*model.JobMember, error) {
	user := ctx.Value("user")
	if user == nil {
		return nil, errors.New("无登录信息")
	}
	sendUserId := strconv.Itoa(int(user.(session.SessionUser).UserID))
	userID := string(data.UserID)
	memberId := int32(data.MemberID)

	//判断任务是否有效
	jobData, err := r.Prisma.JobMember(prisma.JobMemberWhereUniqueInput{
		ID: &memberId,
	}).Exec(ctx)

	if err != nil || jobData == nil || jobData.Progress != constant.STATUS_APPLY {
		return nil, errors.New("数据错误")
	}

	workId := int32(jobData.WorkId)
	work, err := r.Prisma.Work(prisma.WorkWhereUniqueInput{ID: &workId}).Exec(ctx)
	if work == nil {
		return nil, errors.New("数据出错")
	}

	//权限判断
	companyID := int32(user.(session.SessionUser).CompanyID)
	if work.CompanyId != companyID {
		return nil, errors.New("权限不足")
	}

	workTitle := utils.GetWorkTitle(work.CompanyId)

	if work.EndAt != nil {
		//招募截止为空则无限招募时间 不为空判断招募时间
		if int32(time.Now().Unix()) > int32(*work.EndAt) {
			//招募时间已截止
			return nil, errors.New(workTitle + "不允许参加")
		}
	}

	if *work.Status != int32(constant.WORK_STATUS_NORMAL) {
		//任务状态不正常
		return nil, errors.New(workTitle + "不允许参加")
	}

	//可以参加任务
	statusApprove := int32(constant.STATUS_APPROVE)
	nowTimeFormat := int32(time.Now().Unix())
	jobMember, err := r.Prisma.UpdateJobMember(prisma.JobMemberUpdateParams{
		Data: prisma.JobMemberUpdateInput{
			Progress:      &statusApprove,
			ParticipateAt: &nowTimeFormat,
		},
		Where: prisma.JobMemberWhereUniqueInput{
			ID: &jobData.ID,
		},
	}).Exec(ctx)
	if jobMember == nil || err != nil {
		return nil, errors.New(workTitle + "参加失败")
	}

	//入库任务进程表
	statusApplyProcess := string(constant.PROGRESS_STATUS_APPROVE)
	_, _ = r.Prisma.CreateWorkProgress(prisma.WorkProgressCreateInput{
		AppId:         work.AppId,
		PublisherId:   work.UserId,
		ParticipantId: &userID,
		WorkId:        workId,
		Type:          statusApplyProcess,
	}).Exec(ctx)

	//任务参加回调
	//sUserID, _ := strconv.Atoi(userID)
	//_, _ = Api.UserJoinJobCallBack(int64(sUserID), workId)

	//发送消息
	type T struct {
		CompanyName string `json:"companyName"`
	}
	var t1 T
	err = json.Unmarshal([]byte(*work.Extend), &t1)
	sendParam := Api.GetSendMsgContentData{
		WorkId:       workId,
		SendId:       sendUserId,
		ReceiverId:   userID,
		SendType:     1,
		CompanyName:  t1.CompanyName,
		TaskMemberId: jobMember.ID,
		WorkName:     work.Name,
	}
	_, _ = Api.SendMessage(sendParam)

	t := model.JobMember{ID: int(jobMember.ID)}
	return &t, err
}

//拒绝参加
func (r *mutationResolver) Refuse(ctx context.Context, data model.JobInput) (*model.JobMember, error) {
	user := ctx.Value("user")
	if user == nil {
		return nil, errors.New("无登录信息")
	}
	sendUserId := strconv.Itoa(int(user.(session.SessionUser).UserID))
	userID := string(data.UserID)
	memberId := int32(data.MemberID)

	//判断任务是否有效
	jobData, err := r.Prisma.JobMember(prisma.JobMemberWhereUniqueInput{
		ID: &memberId,
	}).Exec(ctx)

	if err != nil || jobData == nil || jobData.Progress != constant.STATUS_APPLY {
		return nil, errors.New("数据错误")
	}

	workId := int32(jobData.WorkId)
	work, err := r.Prisma.Work(prisma.WorkWhereUniqueInput{ID: &workId}).Exec(ctx)
	if work == nil {
		return nil, errors.New("数据出错")
	}

	//权限判断
	companyID := int32(user.(session.SessionUser).CompanyID)
	if work.CompanyId != companyID {
		return nil, errors.New("权限不足")
	}

	//拒绝参加任务
	statusApprove := int32(constant.STATUS_REFUSE_APPLY)
	jobMember, err := r.Prisma.UpdateJobMember(prisma.JobMemberUpdateParams{
		Data: prisma.JobMemberUpdateInput{
			Progress: &statusApprove,
		},
		Where: prisma.JobMemberWhereUniqueInput{
			ID: &jobData.ID,
		},
	}).Exec(ctx)
	if jobMember == nil || err != nil {
		return nil, errors.New("拒绝参加失败")
	}

	//入库任务进程表
	statusApplyProcess := string(constant.PROGRESS_STATUS_REFUSE)
	_, _ = r.Prisma.CreateWorkProgress(prisma.WorkProgressCreateInput{
		AppId:         work.AppId,
		PublisherId:   work.UserId,
		ParticipantId: &userID,
		WorkId:        workId,
		Type:          statusApplyProcess,
	}).Exec(ctx)

	//发送消息
	type T struct {
		CompanyName string `json:"companyName"`
	}
	var t1 T
	err = json.Unmarshal([]byte(*work.Extend), &t1)
	sendParam := Api.GetSendMsgContentData{
		WorkId:       workId,
		SendId:       sendUserId,
		ReceiverId:   userID,
		SendType:     2,
		CompanyName:  t1.CompanyName,
		TaskMemberId: jobMember.ID,
		WorkName:     work.Name,
	}
	_, _ = Api.SendMessage(sendParam)

	t := model.JobMember{ID: int(jobMember.ID)}
	return &t, err
}

func (r *queryResolver) JobDetail(ctx context.Context, workID int, taskMemberID *int) (*model.JobInfo, error) {
	wId := int32(workID)

	user := ctx.Value("user")
	if user == nil {
		return nil, errors.New("无登录信息")
	}
	logInId := strconv.Itoa(int(user.(session.SessionUser).UserID))
	userId := user.(session.SessionUser).UniqueUserId
	var template model.JobTemplate
	//任务数据
	jobArr, err := r.Prisma.Jobs(&prisma.JobsParams{
		Where: &prisma.JobWhereInput{
			WorkId: &wId,
		},
	}).Exec(ctx)

	if jobArr == nil || err != nil {
		return nil, errors.New("数据出错")
	}

	template_id := jobArr[0].TemplateId
	template_info, err := r.Prisma.JobTemplate(prisma.JobTemplateWhereUniqueInput{
		ID: &template_id,
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	_ = mapstructure.Decode(template_info, &template)

	//work数据
	work, err := r.Prisma.Work(prisma.WorkWhereUniqueInput{
		ID: &wId,
	}).Exec(ctx)

	if work == nil || err != nil {
		return nil, errors.New("数据出错")
	}

	workTitle := utils.GetWorkTitle(work.CompanyId)

	if int(*(work.Status)) != constant.WORK_STATUS_NORMAL &&
		int(*(work.Status)) != constant.WORK_STATUS_CLOSED {
		return nil, errors.New("该" + workTitle + "审核中或已下架")
	}
	//progress数据
	/*progressData,err := r.Prisma.WorkProgresses(&prisma.WorkProgressesParams{
		Where:&prisma.WorkProgressWhereInput{
			WorkId:&jobData.WorkId,
		},
	}).Exec(ctx)

	if err != nil {
		return nil, errors.New("数据出错")
	}*/

	//查询总参与人数(进行中)
	progressData, err := r.Prisma.JobMembersConnection(&prisma.JobMembersConnectionParams{
		Where: &prisma.JobMemberWhereInput{
			WorkId:   &jobArr[0].WorkId,
			Progress: prisma.Int32(constant.STATUS_APPROVE),
		},
	}).Aggregate(ctx)

	if err != nil {
		return nil, err
	}

	//查询关注人数
	oldExtend, _ := r.Prisma.WorkExtends(&prisma.WorkExtendsParams{
		Where: &prisma.WorkExtendWhereInput{
			WorkId: prisma.Int32(wId),
		},
	}).Exec(ctx)
	//设置关注人数
	focusCount := JobService.GetFocusCount(constant.TASK_FOCUS_NUM_APPLY)
	newCount := oldExtend[0].FocusCount + focusCount
	_, _ = r.Prisma.UpdateWorkExtend(prisma.WorkExtendUpdateParams{
		Data: prisma.WorkExtendUpdateInput{
			FocusCount: &newCount,
		},
		Where: prisma.WorkExtendWhereUniqueInput{
			ID: &wId,
		},
	}).Exec(ctx)

	//判断当前用户是发布者还是参与者(未参与者)
	//查询jobmember
	var jobMember []prisma.JobMember
	var achievementCount int32

	if logInId == work.UserId {
		//发布者
		//查询绩效总条数
		res, _ := Api.GetJobAchievementCount(work.CompanyId, wId)
		achievementCount = res.Data
	} else if userId != "" {
		//参与者
		whereField := prisma.JobMemberWhereInput{
			WorkId:        prisma.Int32(wId),
			ParticipantId: &userId,
		}

		if taskMemberID != nil {
			taskMemberId := int32(*taskMemberID)
			whereField.ID = &taskMemberId
		}

		orderBy := prisma.JobMemberOrderByInputIDDesc
		jobMember, err = r.Prisma.JobMembers(&prisma.JobMembersParams{
			Where:   &whereField,
			OrderBy: &orderBy,
			First:   prisma.Int32(1),
		}).Exec(ctx)

		if err != nil && jobMember == nil {
			return nil, errors.New("数据出错")
		}

	}

	if jobMember == nil {
		jobMember = append(jobMember, prisma.JobMember{})
	}

	if work.Duration == nil {
		work.Duration = utils.Int2PointInt32(0)
	}
	duration := int(*work.Duration)
	if work.EndAt == nil {
		work.EndAt = utils.Int2PointInt32(0)
	}
	endAt := int(*work.EndAt)
	if work.Status == nil {
		work.Status = utils.Int2PointInt32(0)
	}
	status := int(*work.Status)

	workType := int(work.WorkType)
	payType := int(work.PayType)

	source := int(work.Source)
	types := int(work.Type)

	var mediaUrls []*string
	if work.MediaUrls != nil {
		_ = json.Unmarshal([]byte(*work.MediaUrls), &mediaUrls)
	}

	base := model.Work{
		ID:            int(work.ID),
		Appid:         work.AppId,
		CompanyID:     int(work.CompanyId),
		UserID:        work.UserId,
		ServiceTypeID: int(work.ServiceTypeId),
		WorkType:      &workType,
		Name:          work.Name,
		Requirement:   work.Requirement,
		PayType:       &payType,
		Duration:      &duration,
		EndAt:         &endAt,
		Source:        &source,
		Status:        &status,
		Type:          &types,
		IsPublic:      int(work.IsPublic),
		MediaCoverURL: work.MediaCoverUrl,
		MediaUrls:     mediaUrls,
		Resume:        work.Resume,
		Extend:        work.Extend,
		CreatedAt:     JobService.DateTimeToTimestamp(work.CreatedAt),
	}
	jobData := jobArr[0]
	//wIdString := string(wId)
	if jobData.PayStatus == nil {
		jobData.PayStatus = utils.Int2PointInt32(0)
	}
	payStatus := int(*jobData.PayStatus)
	quota := int(jobData.Quota)
	isCanComment := int(jobData.IsCanComment)
	proofType := int(*jobData.ProofType)
	var isShowProgress int
	isShowProgress = 1
	if jobMember[0].ID == 0 {
		isShowProgress = 0
	}

	Specify := model.Job{
		WorkID:           workID,
		Category:         int(jobData.Category),
		PayStatus:        &payStatus,
		Progress:         int(jobData.Progress),
		Quota:            &quota,
		SingleRewardMin:  *jobData.SingleRewardMin,
		SingleRewardMax:  *jobData.SingleRewardMax,
		IsCanComment:     &isCanComment,
		IsNeedProof:      int(jobData.IsNeedProof),
		ProofDescription: jobData.ProofDescription,
		ProofType:        &proofType,
		Extend:           jobData.Extend,
		UpdatedAt:        JobService.DateTimeToTimestamp(jobData.UpdatedAt),
		MemberCount:      (*int)(unsafe.Pointer(&progressData.Count)),
		AchievementCount: utils.Int322PointInt(achievementCount),
		IsShowProgress:   &isShowProgress,
	}

	memberData := model.JobMember{
		ID:            int(jobMember[0].ID),
		WorkID:        utils.Int322PointInt(jobMember[0].WorkId),
		PublisherID:   &(jobMember[0].PublisherId),
		ParticipantID: jobMember[0].ParticipantId,
		Source:        utils.Int322PointInt(jobMember[0].Source),
		Progress:      utils.Int322PointInt(jobMember[0].Progress),
		ProofFileURL:  jobMember[0].ProofFileUrl,
		ParticipateAt: (*int)(unsafe.Pointer(&(jobMember[0].ParticipateAt))),
		FinishAt:      (*int)(unsafe.Pointer(&(jobMember[0].FinishAt))),
		Extend:        jobMember[0].Extend,
		CreatedAt:     JobService.DateTimeToTimestamp(jobMember[0].CreatedAt),
	}

	t := model.JobInfo{
		Base:     &base,
		Specify:  &Specify,
		Member:   &memberData,
		Template: &template,
	}
	return &t, err
}

//参与记录详情
func (r *queryResolver) JobMember(ctx context.Context, workID int, status *int, pageItem *int, pageNumber int) (*model.JobMemberPagination, error) {
	wId := int32(workID)
	user := ctx.Value("user")
	if user == nil {
		return nil, errors.New("无登录信息")
	}

	//work数据
	work, err := r.Prisma.Work(prisma.WorkWhereUniqueInput{
		ID: &wId,
	}).Exec(ctx)

	if work == nil || err != nil {
		return nil, errors.New("数据出错")
	}

	//权限判断
	companyID := int32(user.(session.SessionUser).CompanyID)
	if work.CompanyId != companyID {
		//return nil, errors.New("权限不足")
	}

	whereInput := prisma.JobMemberWhereInput{
		WorkId: &wId,
	}

	if status != nil {
		progress := int32(*status)
		whereInput.Progress = &progress
	}

	//获取总页数
	pageData, err := r.Prisma.JobMembersConnection(&prisma.JobMembersConnectionParams{
		Where: &whereInput,
	}).Aggregate(ctx)
	if err != nil {
		return nil, err
	}

	size, skip := JobService.GetPrismaPageParam(pageNumber, pageItem)

	jobMember, err := r.Prisma.JobMembers(&prisma.JobMembersParams{
		Where: &whereInput,
		First: prisma.Int32(int32(size)),
		Skip:  prisma.Int32(int32(skip)),
	}).Exec(ctx)

	/*jobMember,err := r.Prisma.JobMembersConnection(&prisma.JobMembersConnectionParams{
		Where:&prisma.JobMemberWhereInput{
			WorkId:&wId,
		},
	}).Edges().Exec(ctx)*/

	if err != nil {
		return nil, err
	}
	jobMemberAll := []*model.JobMember{}

	for i := 0; i < len(jobMember); i++ {

		userId := *jobMember[i].ParticipantId

		//查询用户数据
		userData, _ := r.Prisma.CommonlyUsedPersonnels(&prisma.CommonlyUsedPersonnelsParams{
			Where: &prisma.CommonlyUsedPersonnelWhereInput{
				UserId:    &userId,
				CompanyId: &companyID,
			},
		}).Exec(ctx)
		if userData == nil {
			userData = append(userData, prisma.CommonlyUsedPersonnel{})
		}
		//查询任务历程
		progressData, _ := r.Prisma.WorkProgresses(&prisma.WorkProgressesParams{
			Where: &prisma.WorkProgressWhereInput{
				ParticipantId: jobMember[i].ParticipantId,
				WorkId:        &wId,
				Type:          prisma.Str(constant.PROGRESS_STATUS_APPROVE_INVITE),
			},
		}).Exec(ctx)

		progressAll := []*model.WorkProgress{}
		for j := 0; j < len(progressData); j++ {
			pWorkId := int(progressData[j].WorkId)
			pPublisherId := progressData[j].PublisherId
			pParticipantId := *progressData[j].ParticipantId
			pType := string(progressData[j].Type)
			pCreatedAt := string(progressData[j].CreatedAt)
			progressAll = append(progressAll, &model.WorkProgress{
				ID:            int(progressData[j].ID),
				ParticipantID: &pParticipantId,
				PublisherID:   &pPublisherId,
				WorkID:        &pWorkId,
				Type:          &pType,
				CreatedAt:     JobService.DateTimeToTimestamp(pCreatedAt),
			})
		}
		publishId := jobMember[i].PublisherId
		participantId := *jobMember[i].ParticipantId
		workId := int(jobMember[i].WorkId)
		progress := int(jobMember[i].Progress)
		updatedAt := string(jobMember[i].UpdatedAt)
		createdAt := string(jobMember[i].CreatedAt)
		jobMemberAll = append(jobMemberAll, &model.JobMember{
			ID:            int(jobMember[i].ID),
			WorkID:        &workId,
			PublisherID:   &publishId,
			ParticipantID: &participantId,
			Progress:      &progress,
			Remark:        jobMember[i].Remark,
			CreatedAt:     JobService.DateTimeToTimestamp(createdAt),
			UpdatedAt:     JobService.DateTimeToTimestamp(updatedAt),
			ParticipantUser: &model.CommonlyUsedPersonnel{
				UserID: userData[0].UserId,
				Name:   userData[0].Name,
			},
			WorkProgress: progressAll,
		})
	}

	l := model.JobMemberPagination{
		TotalItem: int(pageData.Count),
		TotalPage: int(math.Ceil(float64(pageData.Count) / float64(size))),
		Items:     jobMemberAll,
	}
	return &l, nil
}

//修改任务参与人员备注
func (r *mutationResolver) ChangeMemberRemark(ctx context.Context, memberID int, remark *string) (*bool, error) {
	user := ctx.Value("user")
	if user == nil {
		return nil, errors.New("无登录信息")
	}

	mId := int32(memberID)

	//判断user_id是否有权限
	jobMember, err := r.Prisma.JobMember(prisma.JobMemberWhereUniqueInput{
		ID: &mId,
	}).Exec(ctx)

	if jobMember == nil && err != nil {
		return nil, errors.New("数据出错")
	}

	//更改备注
	res, err := r.Prisma.UpdateJobMember(prisma.JobMemberUpdateParams{
		Where: prisma.JobMemberWhereUniqueInput{
			ID: &mId,
		},
		Data: prisma.JobMemberUpdateInput{
			Remark: remark,
		},
	}).Exec(ctx)
	if res == nil || err != nil {
		return nil, errors.New("入库失败,请重试")
	}

	status := true
	return &status, err
}

//获取首页统计数据
func (r *queryResolver) DataStatistics(ctx context.Context, companyID int) (*model.Statistics, error) {
	user := ctx.Value("user")
	if user == nil {
		return nil, errors.New("无登录信息")
	}

	companyId := int32(companyID)
	//任务
	job, err := r.Prisma.WorksConnection(&prisma.WorksConnectionParams{
		Where: &prisma.WorkWhereInput{
			CompanyId: prisma.Int32(companyId),
			Status:    prisma.Int32(constant.WORK_STATUS_NORMAL),
		},
	}).Aggregate(ctx)

	if err != nil {
		return nil, errors.New("查询失败")
	}

	//任务人员
	jobMember, err := r.Prisma.JobMembersConnection(&prisma.JobMembersConnectionParams{
		Where: &prisma.JobMemberWhereInput{
			CompanyId: prisma.Int32(companyId),
			Progress:  prisma.Int32(constant.STATUS_APPROVE),
		},
	}).Aggregate(ctx)

	if err != nil {
		return nil, errors.New("查询失败")
	}

	jobCount := int(job.Count)
	jobMemberCount := int(jobMember.Count)
	jobStatistics := model.JobStatistics{
		JobCount:       &jobCount,
		JobMemberCount: &jobMemberCount,
	}

	//人员
	//本月新参与任务人数(这月参加了多次只算一次)
	mouthBeginTime := JobService.GetFirstDateOfMonth(time.Now()).Format(constant.DateTimeLayoutWithTimeZone)

	jobMemberNew, err := r.Prisma.JobMembers(&prisma.JobMembersParams{
		Where: &prisma.JobMemberWhereInput{
			CreatedAtGte: &mouthBeginTime,
			CompanyId:    prisma.Int32(companyId),
			Progress:     prisma.Int32(constant.STATUS_APPROVE),
		},
	}).Exec(ctx)

	if err != nil {
		return nil, errors.New("查询失败")
	}
	//去重
	strMap := make(map[string]string)
	for i := 0; i < len(jobMemberNew); i++ {
		strMap[*jobMemberNew[i].ParticipantId] = *jobMemberNew[i].ParticipantId
	}

	//查询待处理数据(申请参加)
	toDoData, err := r.Prisma.JobMembersConnection(&prisma.JobMembersConnectionParams{
		Where: &prisma.JobMemberWhereInput{
			Progress:  prisma.Int32(constant.STATUS_APPLY),
			CompanyId: prisma.Int32(companyId),
		},
	}).Aggregate(ctx)

	if err != nil {
		return nil, err
	}

	joinMember := len(strMap)
	toDoCount := int(toDoData.Count)
	memberStatistics := model.MemberStatistics{
		JoinMember: &joinMember,
		ToDoCount:  &toDoCount,
	}

	//企业签约数据查询
	signData, _ := Api.GetCompanySignData(companyId)
	signStatistics := []*model.SignStatistics{}
	if signData.Code == http.StatusOK {
		signArrData := signData.Data
		for o := 0; o < len(signData.Data); o++ {
			serviceCompanyID, _ := strconv.Atoi(signArrData[o].ServiceCompanyId)
			signStatistics = append(signStatistics, &model.SignStatistics{
				ServiceCompanyID:   &serviceCompanyID,
				ServiceCompanyName: &(signArrData[o].ServiceCompanyName),
				ServiceTypeName:    &(signArrData[o].ServiceTypeName),
			})
		}
	}
	Statistics := model.Statistics{
		Job:    &jobStatistics,
		Member: &memberStatistics,
		Sign:   signStatistics,
	}
	return &Statistics, nil
}

func (r *queryResolver) CompanyProvideAmount(ctx context.Context, companyID int) (*model.CompanyStatistics, error) {
	user := ctx.Value("user")
	if user == nil {
		return nil, errors.New("无登录信息")
	}

	companyId := int32(companyID)
	//数据统计
	//获取最近7天时间
	bTime := time.Now().AddDate(0, 0, -7).Format(constant.YmdLayoutNoLine)
	eTime := time.Now().Format(constant.YmdLayoutNoLine)
	companyData, _ := Api.GetCompanyProvideAmount(companyId, bTime, eTime)

	weekStatistics := []*model.WeekStatistics{}
	if companyData.Code == http.StatusOK {
		companyArrData := companyData.Data
		for j := 0; j < len(companyData.Data); j++ {
			weekStatistics = append(weekStatistics, &model.WeekStatistics{
				Day:   &(companyArrData[j].Billdate),
				Money: &(companyArrData[j].Amount),
				Count: utils.Int322PointInt(companyArrData[j].Count),
			})
		}
	}
	Statistics := model.CompanyStatistics{
		Statistics: weekStatistics,
	}
	return &Statistics, nil
}

//邀请列表
func (r *queryResolver) InviteList(ctx context.Context, pageItem *int, pageNumber int, appID *string) (*model.ListPagination, error) {
	user := ctx.Value("user")
	if user == nil {
		return nil, errors.New("无登录信息")
	}
	userId := user.(session.SessionUser).UniqueUserId

	//获取总页数
	pageData, err := r.Prisma.JobMembersConnection(&prisma.JobMembersConnectionParams{
		Where: &prisma.JobMemberWhereInput{
			ParticipantId: prisma.Str(userId),
			Progress:      prisma.Int32(constant.STATUS_INVITE),
		},
	}).Aggregate(ctx)
	if err != nil {
		return nil, err
	}

	size, skip := JobService.GetPrismaPageParam(pageNumber, pageItem)

	orderBy := prisma.JobMemberOrderByInputCreatedAtDesc
	jobMember, err := r.Prisma.JobMembers(&prisma.JobMembersParams{
		Where: &prisma.JobMemberWhereInput{
			ParticipantId: prisma.Str(userId),
			Progress:      prisma.Int32(constant.STATUS_INVITE),
		},
		OrderBy: &orderBy,
		First:   prisma.Int32(int32(size)),
		Skip:    prisma.Int32(int32(skip)),
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	jobInfo := []*model.JobInfo{}
	for i := 0; i < len(jobMember); i++ {
		//查询work
		workData, _ := r.Prisma.Work(prisma.WorkWhereUniqueInput{
			ID: prisma.Int32(jobMember[i].WorkId),
		}).Exec(ctx)

		//查询job
		jobData, _ := r.Prisma.Jobs(&prisma.JobsParams{
			Where: &prisma.JobWhereInput{
				WorkId: prisma.Int32(jobMember[i].WorkId),
			},
		}).Exec(ctx)

		//查询用户数据
		userData, _ := r.Prisma.CommonlyUsedPersonnels(&prisma.CommonlyUsedPersonnelsParams{
			Where: &prisma.CommonlyUsedPersonnelWhereInput{
				UserId:    prisma.Str(jobMember[i].PublisherId),
				CompanyId: prisma.Int32(jobMember[i].CompanyId),
			},
		}).Exec(ctx)

		//查询模板数据
		templateData, _ := r.Prisma.JobTemplate(prisma.JobTemplateWhereUniqueInput{
			ID: prisma.Int32(jobData[0].TemplateId),
		}).Exec(ctx)

		if jobData == nil {
			jobData = append(jobData, prisma.Job{})
		}
		if userData == nil {
			userData = append(userData, prisma.CommonlyUsedPersonnel{})
		}

		if templateData == nil {
			templateData = &prisma.JobTemplate{}
		}
		//拼数据
		if workData.Duration == nil {
			workData.Duration = utils.Int2PointInt32(0)
		}
		pDuration := int(*workData.Duration)
		if workData.EndAt == nil {
			workData.EndAt = utils.Int2PointInt32(0)
		}
		pEndAt := int(*workData.EndAt)
		if workData.Status == nil {
			workData.Status = utils.Int2PointInt32(0)
		}
		pStatus := int(*workData.Status)
		pCreatedAt := string(workData.CreatedAt)

		if jobData[0].PayStatus == nil {
			jobData[0].PayStatus = utils.Int2PointInt32(0)
		}
		jPayStatus := int(*jobData[0].PayStatus)
		if jobData[0].SingleRewardMin == nil {
			f := float64(0)
			jobData[0].SingleRewardMin = &f
		}
		if jobData[0].SingleRewardMax == nil {
			f := float64(0)
			jobData[0].SingleRewardMax = &f
		}
		jSingleRewardMin := jobData[0].SingleRewardMin
		jSingleRewardMax := jobData[0].SingleRewardMax
		if jobData[0].ProofType == nil {
			jobData[0].ProofType = utils.Int2PointInt32(0)
		}
		jProofType := int(*jobData[0].ProofType)

		mWorkId := int(jobMember[i].WorkId)
		mParticipantId := *jobMember[i].ParticipantId
		if jobMember[i].ParticipateAt == nil {
			jobMember[i].ParticipateAt = utils.Int2PointInt32(0)
		}
		mParticipateAt := int(*jobMember[i].ParticipateAt)
		if jobMember[i].FinishAt == nil {
			jobMember[i].FinishAt = utils.Int2PointInt32(0)
		}
		mFinishAt := int(*jobMember[i].FinishAt)

		var mediaUrls []*string
		if workData.MediaUrls != nil {
			_ = json.Unmarshal([]byte(*workData.MediaUrls), &mediaUrls)
		}

		jobInfo = append(jobInfo, &model.JobInfo{
			Base: &model.Work{
				ID:            int(workData.ID),
				Appid:         workData.AppId,
				CompanyID:     int(workData.CompanyId),
				UserID:        string(workData.UserId),
				ServiceTypeID: int(workData.ServiceTypeId),
				WorkType:      utils.Int322PointInt(workData.WorkType),
				Name:          workData.Name,
				Requirement:   workData.Requirement,
				PayType:       utils.Int322PointInt(workData.PayType),
				Duration:      &pDuration,
				EndAt:         &pEndAt,
				Source:        utils.Int322PointInt(workData.Source),
				Status:        &pStatus,
				Type:          utils.Int322PointInt(workData.Type),
				IsPublic:      int(workData.IsPublic),
				MediaCoverURL: workData.MediaCoverUrl,
				MediaUrls:     mediaUrls,
				Extend:        workData.Extend,
				CreatedAt:     JobService.DateTimeToTimestamp(pCreatedAt),
			},
			Specify: &model.Job{
				WorkID:           int(jobData[0].WorkId),
				Category:         int(jobData[0].Category),
				PayStatus:        &jPayStatus,
				Progress:         int(jobData[0].Progress),
				Quota:            utils.Int322PointInt(jobData[0].Quota),
				SingleRewardMin:  *jSingleRewardMin,
				SingleRewardMax:  *jSingleRewardMax,
				IsCanComment:     utils.Int322PointInt(jobData[0].IsCanComment),
				IsNeedProof:      int(jobData[0].IsNeedProof),
				ProofDescription: jobData[0].ProofDescription,
				ProofType:        &jProofType,
				Remark:           jobData[0].Remark,
				Extend:           jobData[0].Extend,
			},
			Member: &model.JobMember{
				ID:            int(jobMember[i].ID),
				WorkID:        &mWorkId,
				PublisherID:   &(jobMember[i].PublisherId),
				ParticipantID: &mParticipantId,
				Source:        utils.Int322PointInt(jobMember[i].Source),
				Progress:      utils.Int322PointInt(jobMember[i].Progress),
				ProofFileURL:  jobMember[i].ProofFileUrl,
				ParticipateAt: &mParticipateAt,
				FinishAt:      &mFinishAt,
				Extend:        jobMember[i].Extend,
				CreatedAt:     JobService.DateTimeToTimestamp(jobMember[i].CreatedAt),
				PublishUser: &model.CommonlyUsedPersonnel{
					UserID: userData[0].UserId,
					Name:   userData[0].Name,
				},
			},
			Template: &model.JobTemplate{
				ServiceTypeName: templateData.ServiceTypeName,
				CompanyName:     templateData.CompanyName,
			},
		})

	}

	jobAll := model.ListPagination{
		TotalItem: int(pageData.Count),
		TotalPage: int(math.Ceil(float64(pageData.Count) / float64(size))),
		Items:     jobInfo,
	}

	return &jobAll, err
}

//可参加列表
func (r *queryResolver) List(ctx context.Context, pageNumber int, pageItem int, appID *string) (*model.JobPagination, error) {
	userValue := ctx.Value("user")
	user := userValue.(session.SessionUser)
	var err error
	var works *[]gormModel.Work
	total, totalPage := 0, 0

	paginator := &model.JobPagination{
		PageInfo: &model.PageInfo{
			TotalItem: total,
			TotalPage: totalPage,
		},
	}

	work := &gormModel.Work{}

	offset := (pageNumber - 1) * pageItem
	total, works, err = work.CanJoins(user.UniqueUserId, offset, pageItem, appID)

	if err != nil {
		return nil, errors.New("系统异常")

	}
	log.Printf("%v", works)

	if total == 0 || works == nil || len(*works) == 0 {
		return paginator, nil
	}

	var jobs []*model.JobInfo

	for _, work := range *works {
		var base model.Work
		var specify model.Job

		var mediaUrls []*string
		if work.MediaUrls != "" {
			_ = json.Unmarshal([]byte(work.MediaUrls), &mediaUrls)
		}

		tmpMap := structs.Map(work)
		if _, ok := tmpMap["CreatedAt"]; ok {
			tmpMap["CreatedAt"] = work.CreatedAt.Local().Unix()
		}
		_ = mapstructure.Decode(tmpMap, &base)
		base.MediaUrls = mediaUrls

		var j gormModel.Job
		err = db.Get().Unscoped().Where("work_id = ?", work.ID).First(&j).Error
		if err != nil {
			log.Printf("找不到任相关信息:" + err.Error())
			continue
		}
		_ = mapstructure.Decode(j, &specify)

		//查询template信息
		var t gormModel.JobTemplate
		var template model.JobTemplate
		err = db.Get().Unscoped().Where("id = ?", specify.TemplateID).First(&t).Error
		if err != nil {
			log.Printf("找不到模板信息:" + err.Error())
			continue
		}

		_ = mapstructure.Decode(t, &template)

		//members

		//查询用户数据
		userData, _ := r.Prisma.CommonlyUsedPersonnels(&prisma.CommonlyUsedPersonnelsParams{
			Where: &prisma.CommonlyUsedPersonnelWhereInput{
				UserId:    &work.UserId,
				CompanyId: &work.CompanyId,
			},
		}).Exec(ctx)
		if userData == nil {
			userData = append(userData, prisma.CommonlyUsedPersonnel{})
		}

		job := model.JobInfo{
			Base:    &base,
			Specify: &specify,
			Member: &model.JobMember{
				PublishUser: &model.CommonlyUsedPersonnel{
					UserID: userData[0].UserId,
					Name:   userData[0].Name,
				},
			},
			Template: &template,
		}
		jobs = append(jobs, &job)
	}

	totalPage = int(math.Ceil(float64(total) / float64(pageItem)))

	return &model.JobPagination{
		PageInfo: &model.PageInfo{
			TotalItem: total,
			TotalPage: totalPage,
		},
		Items: jobs,
	}, nil
}

//已参加列表
func (r *queryResolver) JoinList(ctx context.Context, pageItem *int, pageNumber int, appID *string) (*model.ListPagination, error) {
	user := ctx.Value("user")
	if user == nil {
		return nil, errors.New("无登录信息")
	}
	userId := user.(session.SessionUser).UniqueUserId
	//获取总页数
	pageData, err := r.Prisma.JobMembersConnection(&prisma.JobMembersConnectionParams{
		Where: &prisma.JobMemberWhereInput{
			ParticipantId: prisma.Str(userId),
			ProgressIn:    constant.JOB_STATUS_JOINED,
		},
	}).Aggregate(ctx)
	if err != nil {
		return nil, err
	}

	size, skip := JobService.GetPrismaPageParam(pageNumber, pageItem)

	orderBy := prisma.JobMemberOrderByInputParticipateAtDesc
	jobMember, err := r.Prisma.JobMembers(&prisma.JobMembersParams{
		Where: &prisma.JobMemberWhereInput{
			ParticipantId: prisma.Str(userId),
			ProgressIn:    constant.JOB_STATUS_JOINED,
		},
		OrderBy: &orderBy,
		First:   prisma.Int32(int32(size)),
		Skip:    prisma.Int32(int32(skip)),
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}
	jobInfo := []*model.JobInfo{}
	for i := 0; i < len(jobMember); i++ {
		//JobService.GetJobInfoById(jobMember[i],jobInfo)
		//查询work
		workData, _ := r.Prisma.Work(prisma.WorkWhereUniqueInput{
			ID: prisma.Int32(jobMember[i].WorkId),
		}).Exec(ctx)

		//状态不正确的任务跳过
		if *workData.Status != constant.JobTemplateEnable {
			continue
		}

		//查询job
		jobData, _ := r.Prisma.Jobs(&prisma.JobsParams{
			Where: &prisma.JobWhereInput{
				WorkId: prisma.Int32(jobMember[i].WorkId),
			},
		}).Exec(ctx)

		//查询用户数据
		userData, _ := r.Prisma.CommonlyUsedPersonnels(&prisma.CommonlyUsedPersonnelsParams{
			Where: &prisma.CommonlyUsedPersonnelWhereInput{
				UserId:    prisma.Str(jobMember[i].PublisherId),
				CompanyId: prisma.Int32(jobMember[i].CompanyId),
			},
		}).Exec(ctx)

		//查询模板数据
		templateData, _ := r.Prisma.JobTemplate(prisma.JobTemplateWhereUniqueInput{
			ID: prisma.Int32(jobData[0].TemplateId),
		}).Exec(ctx)

		if jobData == nil {
			jobData = append(jobData, prisma.Job{})
		}
		if userData == nil {
			userData = append(userData, prisma.CommonlyUsedPersonnel{})
		}
		if templateData == nil {
			templateData = &prisma.JobTemplate{}
		}
		//拼数据
		if workData.Duration == nil {
			workData.Duration = utils.Int2PointInt32(0)
		}
		pDuration := int(*workData.Duration)
		if workData.EndAt == nil {
			workData.EndAt = utils.Int2PointInt32(0)
		}
		pEndAt := int(*workData.EndAt)
		if workData.Status == nil {
			workData.Status = utils.Int2PointInt32(0)
		}
		pStatus := int(*workData.Status)
		pCreatedAt := string(workData.CreatedAt)

		if jobData[0].PayStatus == nil {
			jobData[0].PayStatus = utils.Int2PointInt32(0)
		}
		jPayStatus := int(*jobData[0].PayStatus)
		if jobData[0].SingleRewardMin == nil {
			f := float64(0)
			jobData[0].SingleRewardMin = &f
		}
		if jobData[0].SingleRewardMax == nil {
			f := float64(0)
			jobData[0].SingleRewardMax = &f
		}
		jSingleRewardMin := jobData[0].SingleRewardMin
		jSingleRewardMax := jobData[0].SingleRewardMax
		if jobData[0].ProofType == nil {
			jobData[0].ProofType = utils.Int2PointInt32(0)
		}
		jProofType := int(*jobData[0].ProofType)

		mWorkId := int(jobMember[i].WorkId)
		mParticipantId := *jobMember[i].ParticipantId
		if jobMember[i].ParticipateAt == nil {
			jobMember[i].ParticipateAt = utils.Int2PointInt32(0)
		}
		mParticipateAt := int(*jobMember[i].ParticipateAt)
		if jobMember[i].FinishAt == nil {
			jobMember[i].FinishAt = utils.Int2PointInt32(0)
		}
		mFinishAt := int(*jobMember[i].FinishAt)

		var mediaUrls []*string
		if workData.MediaUrls != nil {
			_ = json.Unmarshal([]byte(*workData.MediaUrls), &mediaUrls)
		}

		jobInfo = append(jobInfo, &model.JobInfo{
			Base: &model.Work{
				ID:            int(workData.ID),
				Appid:         workData.AppId,
				CompanyID:     int(workData.CompanyId),
				UserID:        string(workData.UserId),
				ServiceTypeID: int(workData.ServiceTypeId),
				WorkType:      utils.Int322PointInt(workData.WorkType),
				Name:          workData.Name,
				Requirement:   workData.Requirement,
				PayType:       utils.Int322PointInt(workData.PayType),
				Duration:      &pDuration,
				EndAt:         &pEndAt,
				Source:        utils.Int322PointInt(workData.Source),
				Status:        &pStatus,
				Type:          utils.Int322PointInt(workData.Type),
				IsPublic:      int(workData.IsPublic),
				MediaCoverURL: workData.MediaCoverUrl,
				MediaUrls:     mediaUrls,
				Extend:        workData.Extend,
				CreatedAt:     JobService.DateTimeToTimestamp(pCreatedAt),
			},
			Specify: &model.Job{
				WorkID:           int(jobData[0].WorkId),
				Category:         int(jobData[0].Category),
				PayStatus:        &jPayStatus,
				Progress:         int(jobData[0].Progress),
				Quota:            utils.Int322PointInt(jobData[0].Quota),
				SingleRewardMin:  *jSingleRewardMin,
				SingleRewardMax:  *jSingleRewardMax,
				IsCanComment:     utils.Int322PointInt(jobData[0].IsCanComment),
				IsNeedProof:      int(jobData[0].IsNeedProof),
				ProofDescription: jobData[0].ProofDescription,
				ProofType:        &jProofType,
				Remark:           jobData[0].Remark,
				Extend:           jobData[0].Extend,
			},
			Member: &model.JobMember{
				ID:            int(jobMember[i].ID),
				WorkID:        &mWorkId,
				PublisherID:   &(jobMember[i].PublisherId),
				ParticipantID: &mParticipantId,
				Source:        utils.Int322PointInt(jobMember[i].Source),
				Progress:      utils.Int322PointInt(jobMember[i].Progress),
				ProofFileURL:  jobMember[i].ProofFileUrl,
				ParticipateAt: &mParticipateAt,
				FinishAt:      &mFinishAt,
				Extend:        jobMember[i].Extend,
				CreatedAt:     JobService.DateTimeToTimestamp(jobMember[i].CreatedAt),
				PublishUser: &model.CommonlyUsedPersonnel{
					UserID: userData[0].UserId,
					Name:   userData[0].Name,
				},
			},
			Template: &model.JobTemplate{
				ServiceTypeName: templateData.ServiceTypeName,
				CompanyName:     templateData.CompanyName,
			},
		})

	}
	jobAll := model.ListPagination{
		TotalItem: int(pageData.Count),
		TotalPage: int(math.Ceil(float64(pageData.Count) / float64(size))),
		Items:     jobInfo,
	}
	return &jobAll, err
}

//申请完成/提交凭证
func (r *mutationResolver) UploadAchievement(ctx context.Context, data model.ApplyCompleteInput) (*bool, error) {
	memberId := int32(data.MemberID)
	user := ctx.Value("user")
	if user == nil {
		return nil, errors.New("无登录信息")
	}
	userId := user.(session.SessionUser).UniqueUserId

	//判断参与记录是否正确
	jobMember, err := r.Prisma.JobMember(prisma.JobMemberWhereUniqueInput{
		ID: &memberId,
	}).Exec(ctx)
	if err != nil {
		return nil, errors.New("数据错误")
	}
	//判断角色
	if userId != *(jobMember.ParticipantId) {
		return nil, errors.New("权限不足")
	}

	//work 判断状态
	workData, err := r.Prisma.Work(prisma.WorkWhereUniqueInput{
		ID: prisma.Int32(jobMember.WorkId),
	}).Exec(ctx)

	workTitle := utils.GetWorkTitle(workData.CompanyId)

	if err != nil {
		return nil, errors.New("数据错误")
	}

	status := utils.Int322PointInt(*workData.Status)
	if *status != constant.WORK_STATUS_NORMAL {
		return nil, errors.New("该" + workTitle + "审核中或已下架")
	}
	//入库任务进程

	//上传的凭证入扩展字段
	var extendField = map[string][]*string{}
	extendField["proof_file_url"] = data.ProofFileURL
	extend, _ := json.Marshal(extendField)
	extendString := string(extend)

	statusProcess := string(constant.PROGRESS_STATUS_JOB_UPLOAD)
	_, _ = r.Prisma.CreateWorkProgress(prisma.WorkProgressCreateInput{
		AppId:         workData.AppId,
		PublisherId:   workData.UserId,
		ParticipantId: &userId,
		WorkId:        jobMember.WorkId,
		Type:          statusProcess,
		Extend:        &extendString,
	}).Exec(ctx)

	v := true
	return &v, nil
}

//任务进程
func (r *queryResolver) Process(ctx context.Context, workID int, pageNumber int, pageItem *int) ([]*model.WorkProgress, error) {
	wId := int32(workID)
	user := ctx.Value("user")
	if user == nil {
		return nil, errors.New("无登录信息")
	}
	userId := user.(session.SessionUser).UniqueUserId

	//判断权限(该用户是否参加该任务)
	jobMember, err := r.Prisma.JobMembersConnection(&prisma.JobMembersConnectionParams{
		Where: &prisma.JobMemberWhereInput{
			WorkId:        &wId,
			ParticipantId: &userId,
		},
	}).Aggregate(ctx)
	if err != nil || jobMember.Count <= 0 {
		return nil, errors.New("权限不足")
	}

	//查询数据
	size, skip := JobService.GetPrismaPageParam(pageNumber, pageItem)

	orderBy := prisma.WorkProgressOrderByInputIDDesc
	workProgress, err := r.Prisma.WorkProgresses(&prisma.WorkProgressesParams{
		Where: &prisma.WorkProgressWhereInput{
			WorkId: &wId,
			Or: []prisma.WorkProgressWhereInput{
				prisma.WorkProgressWhereInput{
					ParticipantId: &userId,
				},
				prisma.WorkProgressWhereInput{
					Type: prisma.Str(constant.PROGRESS_STATUS_CREATE),
				},
			},
		},
		OrderBy: &orderBy,
		First:   prisma.Int32(int32(size)),
		Skip:    prisma.Int32(int32(skip)),
	}).Exec(ctx)

	if err != nil {
		return nil, errors.New("数据错误")
	}

	work, err := r.Prisma.Work(prisma.WorkWhereUniqueInput{ID: &wId}).Exec(ctx)
	if work == nil || err != nil {
		return nil, errors.New("数据出错")
	}

	type T struct {
		CompanyName string `json:"companyName"`
	}
	var t1 T
	err = json.Unmarshal([]byte(*work.Extend), &t1)
	companyName := t1.CompanyName

	//拼接数据
	progressAll := []*model.WorkProgress{}
	for j := 0; j < len(workProgress); j++ {
		//查询用户数据
		/*userData, _ := r.Prisma.CommonlyUsedPersonnels(&prisma.CommonlyUsedPersonnelsParams{
			Where: &prisma.CommonlyUsedPersonnelWhereInput{
				UserId: prisma.Str(workProgress[j].PublisherId),
				CompanyId:prisma.Int32(work.CompanyId),
			},
		}).Exec(ctx)

		if userData == nil {
			userData = append(userData, prisma.CommonlyUsedPersonnel{})
		}*/

		//查询用户数据
		if workProgress[j].Type == constant.PROGRESS_STATUS_CREATE {
			defaultParticipantId := ""
			workProgress[j].ParticipantId = &defaultParticipantId
		}
		ParticipantUserData, _ := r.Prisma.CommonlyUsedPersonnels(&prisma.CommonlyUsedPersonnelsParams{
			Where: &prisma.CommonlyUsedPersonnelWhereInput{
				UserId:    prisma.Str(*(workProgress[j].ParticipantId)),
				CompanyId: prisma.Int32(work.CompanyId),
			},
		}).Exec(ctx)
		if ParticipantUserData == nil {
			ParticipantUserData = append(ParticipantUserData, prisma.CommonlyUsedPersonnel{})
		}

		pWorkId := int(workProgress[j].WorkId)
		pPublisherId := workProgress[j].PublisherId
		pParticipantId := *workProgress[j].ParticipantId
		pType := string(workProgress[j].Type)
		pCreatedAt := string(workProgress[j].CreatedAt)
		progressAll = append(progressAll, &model.WorkProgress{
			ID:            int(workProgress[j].ID),
			ParticipantID: &pParticipantId,
			PublisherID:   &pPublisherId,
			WorkID:        &pWorkId,
			Type:          &pType,
			Extend:        workProgress[j].Extend,
			CreatedAt:     JobService.DateTimeToTimestamp(pCreatedAt),
			PublishUser: &model.CommonlyUsedPersonnel{
				Name: companyName,
			},
			ParticipantUser: &model.CommonlyUsedPersonnel{
				UserID: ParticipantUserData[0].UserId,
				Name:   ParticipantUserData[0].Name,
			},
		})
	}
	return progressAll, err
}

func (r *queryResolver) UserFlowPage(ctx context.Context, userID string, companyID int, workID int, pageNumber int, pageItem *int) (*model.UserFlowPagination, error) {
	wId := int32(workID)
	cID := int32(companyID)
	user := ctx.Value("user")
	if user == nil {
		return nil, errors.New("无登录信息")
	}

	//权限控制(userid对应的workid)
	workData, err := r.Prisma.Work(prisma.WorkWhereUniqueInput{
		ID: &wId,
	}).Exec(ctx)

	if err != nil || workData == nil {
		return nil, errors.New("数据错误")
	}

	if workData.CompanyId != cID {
		return nil, errors.New("权限不足")
	}

	//查询
	item := utils.Int2PointInt32(*pageItem)
	res, _ := Api.GetUserFlowPage(userID, cID, wId, int32(pageNumber), *item)

	if res.Code == http.StatusOK {
		resAll := []*model.UserFlow{}
		resData := res.Data.List
		for i := 0; i < len(resData); i++ {
			resAll = append(resAll, &model.UserFlow{
				Amount:         &(resData[i].Amount),
				PaymentResTime: &(resData[i].PaymentResTime),
				PayOrderItemID: &(resData[i].PayOrderItemId),
			})
		}
		r := model.UserFlowPagination{
			TotalPage: int(res.Data.Pages),
			TotalItem: int(res.Data.Total),
			Items:     resAll,
		}
		return &r, err
	}

	return nil, nil
}

//拉黑
func (r *mutationResolver) PullOnBlackList(ctx context.Context, companyID int, userID string, workID int) (*bool, error) {
	wId := int32(workID)
	cID := int32(companyID)
	user := ctx.Value("user")
	if user == nil {
		return nil, errors.New("无登录信息")
	}
	uId := strconv.Itoa(int(user.(session.SessionUser).UserID))

	//权限控制(userid对应的workid)
	workData, err := r.Prisma.Work(prisma.WorkWhereUniqueInput{
		ID: &wId,
	}).Exec(ctx)

	if err != nil || workData == nil {
		return nil, errors.New("数据错误")
	}

	if workData.CompanyId != cID {
		return nil, errors.New("权限不足")
	}

	//判断userid
	jobData, err := r.Prisma.JobMembers(&prisma.JobMembersParams{
		Where: &prisma.JobMemberWhereInput{
			WorkId:        &wId,
			ParticipantId: prisma.Str(userID),
		},
	}).Exec(ctx)

	if err != nil || jobData == nil {
		return nil, errors.New("数据错误")
	}

	//判断是否已经拉黑
	blackData, err := r.Prisma.CompanyUserBlacklistsConnection(&prisma.CompanyUserBlacklistsConnectionParams{
		Where: &prisma.CompanyUserBlacklistWhereInput{
			WorkId:        &wId,
			ParticipantId: prisma.Str(userID),
		},
	}).Aggregate(ctx)

	if err != nil || blackData.Count > 0 {
		return nil, errors.New("重复操作")
	}
	//拉黑
	insert, err := r.Prisma.CreateCompanyUserBlacklist(prisma.CompanyUserBlacklistCreateInput{
		CompanyId:     cID,
		ParticipantId: userID,
		PublisherId:   uId,
		WorkId:        &wId,
		Type:          prisma.Int32(2),
	}).Exec(ctx)

	if insert == nil || err != nil {
		return nil, errors.New("拉黑失败")
	}
	v := true
	return &v, err
}

//小程序红点数量
func (r *queryResolver) RedDotCount(ctx context.Context) (*model.RedDotCount, error) {
	user := ctx.Value("user")
	if user == nil {
		return nil, errors.New("无登录信息")
	}
	uId := user.(session.SessionUser).UniqueUserId
	//获取当前用户的邀请接岗数量
	var count int
	db.Get().Model(&gormModel.Member{}).Unscoped().
		Where("participant_id = ?  AND progress = ? ", uId, constant.STATUS_INVITE).
		Count(&count)

	//获取当前用户的待结算数量
	var settlementCount int
	db.Get().Model(&gormModel.JobSettlementLog{}).Unscoped().
		Where("user_id = ?  AND status = ? ", uId, constant.UN_CONFIRM_STATUS).
		Count(&settlementCount)

	res := model.RedDotCount{
		Invite:     count,
		Settlement: settlementCount,
		Job:        count,
		All:        count + settlementCount,
	}
	return &res, nil

}

//获取上传地址
func (r *queryResolver) GetUploadURL(ctx context.Context, taskMemberID int) (string, error) {

	var taskMember gormModel.Member
	err := db.Get().Model(&gormModel.Member{}).Unscoped().
		Where("id = ? ", taskMemberID).First(&taskMember).Error
	if err != nil {
		return "", errors.New("查询失败:" + err.Error())
	}

	if taskMember.Progress != constant.STATUS_APPROVE {
		return "", errors.New("状态不正确")
	}

	shortUrl, err := Api.Dwz(int64(taskMemberID))
	return shortUrl, err

}

//上传文件
func (r *mutationResolver) SetUploadURL(ctx context.Context, taskMemberID string, fileURL []*string) (bool, error) {
	hd := hashids.NewData()
	hd.MinLength = 30
	h, _ := hashids.NewWithData(hd)
	d, _ := h.DecodeWithError(taskMemberID)
	var taskMember gormModel.Member
	err := db.Get().Model(&gormModel.Member{}).Unscoped().
		Where("id = ? ", d).First(&taskMember).Error

	if taskMember.Progress != constant.STATUS_APPROVE {
		return false, errors.New("状态不正确")
	}

	//上传的凭证入扩展字段
	var extendField = map[string][]*string{}
	extendField["proof_file_url"] = fileURL
	extend, _ := json.Marshal(extendField)
	extendString := string(extend)

	user := &gormModel.WorkProgress{
		ParticipantId: taskMember.ParticipantId,
		PublisherId:   taskMember.PublisherId,
		WorkId:        taskMember.WorkId,
		Type:          string(constant.PROGRESS_STATUS_JOB_UPLOAD),
		Extend:        extendString,
	}

	if err := db.Get().Create(&user).Error; err != nil {
		return false, err
	}
	return true, err
}

//绑定手机号
func (r *mutationResolver) BindPhone(ctx context.Context, userID string, encryptedData string, sessionKey string, iv string, group string) (string, error) {
	//查询手机号
	var user gormModel.UserWeChatAuthorize
	err := db.Get().Model(&gormModel.UserWeChatAuthorize{}).Unscoped().
		Where("user_id = ? ", userID).First(&user).Error

	if user.Mobile != "" {
		return user.Mobile, nil
	}

	appID := "wxe813fbeb8daa55df"
	if group == "ishouru-xn" {
		appID = "wx4514ec5212c407b7"
	}

	pc := wxbizdatacrypt.WxBizDataCrypt{AppID: appID, SessionKey: sessionKey}
	result, err := pc.Decrypt(encryptedData, iv, true) //第三个参数解释： 需要返回 JSON 数据类型时 使用 true, 需要返回 map 数据类型时 使用 false
	if err != nil {
		return "", err
	}
	type PhoneNumber struct {
		PhoneNumber string `json:"phoneNumber"`
	}
	var s PhoneNumber
	_ = json.Unmarshal([]byte(result.(string)), &s)
	mobile := s.PhoneNumber

	if mobile == "" {
		return "", errors.New("手机号为空")
	}
	weChat := &gormModel.UserWeChatAuthorize{
		UserId: userID,
		Mobile: mobile,
	}

	err = db.Get().Create(&weChat).Error

	if err != nil {
		return "", err
	}
	return mobile, err
}

//判断是否绑定手机号
func (r *queryResolver) IsBindPhone(ctx context.Context, userID string) (string, error) {
	//查询手机号
	var user gormModel.UserWeChatAuthorize
	err := db.Get().Model(&gormModel.UserWeChatAuthorize{}).Unscoped().
		Where("user_id = ? ", userID).First(&user).Error

	if strings.Replace(user.Mobile, " ", "", -1) != "" {
		return user.Mobile, nil
	}
	//查询接口
	userIdInt, _ := strconv.Atoi(userID)
	ver, err := Api.FindRealnameInfoByUserid(int64(userIdInt))
	if err != nil || ver.Code != 0 {
		return "", errors.New("查询用户信息接口出错")
	}
	if strings.Replace(ver.Data.MobilePhone, " ", "", -1) == "" {
		return "", err
	}
	log.Printf("IsBindPhone:userId:%v,MobilePhone:%v", userIdInt, ver.Data.MobilePhone)

	//存入表
	weChat := &gormModel.UserWeChatAuthorize{
		UserId: userID,
		Mobile: user.Mobile,
	}

	err = db.Get().Create(&weChat).Error
	return ver.Data.MobilePhone, err
}

//结算列表,已结算/待结算
func (r *queryResolver) SettlementList(ctx context.Context, pageItem int, pageNumber int, settlementType int) (*model.ListPagination, error) {
	user := ctx.Value("user")
	uId := user.(session.SessionUser).UniqueUserId

	//count
	var count int
	if err := db.Get().Unscoped().Model(&model.JobSettlementLog{}).
		Where("status = ? and user_id = ?", settlementType, uId).
		Count(&count).Error; err != nil {
		return nil, err
	}

	var err error
	offset := (pageNumber - 1) * pageItem
	settlement := []*gormModel.JobSettlementLog{}
	err = db.Get().Unscoped().Model(&model.JobSettlementLog{}).
		Where("status = ? and user_id = ?", settlementType, uId).
		Offset(offset).Limit(pageItem).
		Find(&settlement).Error
	if err != nil {
		return nil, err
	}
	jobs := []*model.JobInfo{}
	for _, settlementLog := range settlement {
		var base model.Work
		var specify model.Job

		var w gormModel.Work
		err = db.Get().Unscoped().Where("id = ?", settlementLog.WorkID).First(&w).Error
		if err != nil {
			log.Printf("找不到信息:" + err.Error())
			continue
		}
		tmpMap := structs.Map(w)
		if _, ok := tmpMap["CreatedAt"]; ok {
			tmpMap["CreatedAt"] = w.CreatedAt.Local().Unix()
		}
		_ = mapstructure.Decode(tmpMap, &base)

		var j gormModel.Job
		err = db.Get().Unscoped().Where("work_id = ?", settlementLog.WorkID).First(&j).Error
		if err != nil {
			log.Printf("找不到信息:" + err.Error())
			continue
		}
		_ = mapstructure.Decode(j, &specify)

		//查询template信息
		var t gormModel.JobTemplate
		var template model.JobTemplate
		err = db.Get().Unscoped().Where("id = ?", specify.TemplateID).First(&t).Error
		if err != nil {
			log.Printf("找不到模板信息:" + err.Error())
			continue
		}

		_ = mapstructure.Decode(t, &template)

		//members

		//查询用户数据
		userData, _ := r.Prisma.CommonlyUsedPersonnels(&prisma.CommonlyUsedPersonnelsParams{
			Where: &prisma.CommonlyUsedPersonnelWhereInput{
				UserId:    &settlementLog.UserID,
				CompanyId: prisma.Int32(int32(base.CompanyID)),
			},
		}).Exec(ctx)
		if userData == nil {
			userData = append(userData, prisma.CommonlyUsedPersonnel{})
		}

		createdAt := int(settlementLog.CreatedAt.Unix())
		job := model.JobInfo{
			Base:    &base,
			Specify: &specify,
			Member: &model.JobMember{
				PublishUser: &model.CommonlyUsedPersonnel{
					UserID: userData[0].UserId,
					Name:   userData[0].Name,
				},
			},
			Template: &template,
			Settlement: &model.Settlement{
				ID:        settlementLog.ID,
				Amount:    &settlementLog.Amount,
				CreatedAt: &createdAt,
			},
		}
		jobs = append(jobs, &job)
	}

	pagination := model.ListPagination{
		TotalItem: int(count),
		TotalPage: int(math.Ceil(float64(count) / float64(pageItem))),
		Items:     jobs,
	}

	return &pagination, nil
}

//获取待结算数据详情
func (r *queryResolver) PendingDetail(ctx context.Context, settlementID int) (*model.Settlement, error) {
	user := ctx.Value("user")
	uId := user.(session.SessionUser).UniqueUserId
	var j gormModel.JobSettlementLog
	err := db.Get().Unscoped().Where("id = ? and user_id = ?", settlementID, uId).First(&j).Error
	if err != nil {
		return nil, err
	}
	createdAt := int(j.CreatedAt.Unix())
	settlement := &model.Settlement{
		ID:        j.ID,
		File:      &j.File,
		Status:    &j.Status,
		Amount:    &j.Amount,
		CreatedAt: &createdAt,
	}

	return settlement, err
}

//批量同意
func (r *mutationResolver) BatchApprove(ctx context.Context, workID int) (bool, error) {
	user := ctx.Value("user")
	if user == nil {
		return false, errors.New("无登录信息")
	}
	sendUserId := strconv.Itoa(int(user.(session.SessionUser).UserID))
	companyID := int32(user.(session.SessionUser).CompanyID)

	workId := int32(workID)
	work, err := r.Prisma.Work(prisma.WorkWhereUniqueInput{ID: &workId}).Exec(ctx)
	if work == nil {
		return false, errors.New("数据出错")
	}
	workTitle := utils.GetWorkTitle(work.CompanyId)

	//权限判断
	if work.CompanyId != companyID {
		return false, errors.New("权限不足")
	}

	if work.EndAt != nil {
		//招募截止为空则无限招募时间 不为空判断招募时间
		if int32(time.Now().Unix()) > int32(*work.EndAt) {
			//招募时间已截止
			return false, errors.New(workTitle + "不允许参加")
		}
	}

	if *work.Status != int32(constant.WORK_STATUS_NORMAL) {
		//任务状态不正常
		return false, errors.New(workTitle + "不允许参加")
	}

	//可以参加任务 查询邀请中任务
	var jobMember []*gormModel.Member
	err = db.Get().Unscoped().Where("work_id = ? and progress = ?", workID, constant.STATUS_APPLY).Find(&jobMember).Error
	if err != nil || err == gorm.ErrRecordNotFound {
		return false, errors.New("没有申请中数据")
	}

	//更新状态
	err = db.Get().Unscoped().Model(gormModel.Member{}).Where("work_id = ? and progress = ?", workID, constant.STATUS_APPLY).UpdateColumn("progress", constant.STATUS_APPROVE).Error
	if err != nil {
		return false, errors.New("更新状态失败")
	}

	//创建
	for i := 0; i < len(jobMember); i++ {
		//入库任务进程表
		statusApplyProcess := string(constant.PROGRESS_STATUS_APPROVE)
		_, _ = r.Prisma.CreateWorkProgress(prisma.WorkProgressCreateInput{
			AppId:         work.AppId,
			PublisherId:   work.UserId,
			ParticipantId: &jobMember[i].ParticipantId,
			WorkId:        workId,
			Type:          statusApplyProcess,
		}).Exec(ctx)

		//发送消息
		type T struct {
			CompanyName string `json:"companyName"`
		}
		var t1 T
		err = json.Unmarshal([]byte(*work.Extend), &t1)
		sendParam := Api.GetSendMsgContentData{
			WorkId:       workId,
			SendId:       sendUserId,
			ReceiverId:   jobMember[i].ParticipantId,
			SendType:     1,
			CompanyName:  t1.CompanyName,
			TaskMemberId: int32(jobMember[i].ID),
			WorkName:     work.Name,
		}
		_, _ = Api.SendMessage(sendParam)
	}
	return true, err
}

//查看上传凭证
func (r *queryResolver) UploadRecord(ctx context.Context, workID int, userID string, pageNumber int, pageItem *int) (*model.UploadRecordPagination, error) {
	user := ctx.Value("user")
	if user == nil {
		return nil, errors.New("无登录信息")
	}

	var count int
	err := db.Get().Model(gormModel.WorkProgress{}).Where("work_id = ? and participant_id = ? and type = ?", workID, userID, constant.PROGRESS_STATUS_JOB_UPLOAD).
		Count(&count).Error
	if err != nil {
		return nil, err
	}

	var workProgress []gormModel.WorkProgress
	offset := (pageNumber - 1) * *pageItem
	err = db.Get().Where("work_id = ? and participant_id = ? and type = ?", workID, userID, constant.PROGRESS_STATUS_JOB_UPLOAD).
		Offset(offset).Limit(pageItem).Find(&workProgress).Error
	if err != nil {
		return nil, err
	}

	works := []*model.UploadRecord{}
	for i := 0; i < len(workProgress); i++ {

		type T struct {
			File []string `json:"proof_file_url"`
		}
		var t1 T
		err = json.Unmarshal([]byte(workProgress[i].Extend), &t1)
		work := model.UploadRecord{
			ID:        workProgress[i].ID,
			File:      t1.File,
			CreatedAt: int(workProgress[i].CreatedAt.Unix()),
		}
		works = append(works, &work)
	}
	total_page := 0
	if int(count) > 0 {
		total_page = int(math.Ceil(float64(count) / float64(*pageItem)))
	}
	res := model.UploadRecordPagination{
		TotalItem: count,
		TotalPage: total_page,
		Items:     works,
	}

	return &res, err
}
