package resolver

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gogf/gf/util/gconv"
	"iQuest/app/graphql/model"
	m "iQuest/app/model"
	"iQuest/app/model/user"
	"iQuest/app/service/job"
	"iQuest/app/service/work"
	"iQuest/db"
	"math"
	"strconv"
	"time"
)

// CreateJobSettlementLogs 创建结算记录
func (r *mutationResolver) CreateJobSettlementLogs(ctx context.Context, data []*model.CreateJobSettlementLogInput) ([]*model.JobSettlementLog, error) {
	u := ctx.Value("user").(user.SessionUser)

	//获取岗位信息
	w, err := work.GetByID(data[0].WorkID)
	if err != nil {
		return nil, err
	}

	if u.CompanyID != w.CompanyId {
		return nil, errors.New("权限有误：公司不匹配")
	}

	var logs []*m.JobSettlementLog
	for _, item := range data {
		if item.Amount == 0 {
			return nil, errors.New("数据有误：金额不能为0")
		}
		var log m.JobSettlementLog
		gconv.Struct(item, &log)
		log.OperatorUserID = gconv.String(u.UserID)
		logs = append(logs, &log)
	}
	logs, err = job.CreateSettlementLogs(logs)
	if err != nil {
		return nil, err
	}
	var settlementLogs []*model.JobSettlementLog
	for _, item := range logs {
		var slog model.JobSettlementLog
		gconv.Struct(item, &slog)
		slog.ID = item.ID
		slog.CreatedAt = int(item.CreatedAt.Unix())
		slog.UpdatedAt = int(item.UpdatedAt.Unix())
		settlementLogs = append(settlementLogs, &slog)

	}
	return settlementLogs, nil
}

// JobSettlementLogs 获取结算记录
func (r *queryResolver) JobSettlementLogs(ctx context.Context, pageNumber int, pageItem int, workID int, userID string) (*model.JobSettlementLogsPagination, error) {
	u := ctx.Value("user").(user.SessionUser)

	//获取岗位信息
	w, err := work.GetByID(workID)
	if err != nil {
		return nil, err
	}

	if u.CompanyID != w.CompanyId {
		return nil, errors.New("权限有误：公司不匹配")
	}

	condition := job.SettlementLogsCondition{
		PageNum:  pageNumber,
		PageSize: pageItem,
		WorkID:   workID,
		UserID:   userID,
	}

	logs, total, err := job.GetSettlementLogs(condition)
	if err != nil {
		return nil, err
	}
	pageInfo := model.PageInfo{
		TotalItem: int(total),
		TotalPage: int(math.Ceil(float64(total) / float64(pageItem))),
	}
	var settlementLogs []*model.JobSettlementLog

	for _, item := range logs {
		var slog model.JobSettlementLog
		gconv.Struct(item, &slog)
		slog.ID = item.ID
		slog.CreatedAt = int(item.CreatedAt.Unix())
		slog.UpdatedAt = int(item.UpdatedAt.Unix())
		settlementLogs = append(settlementLogs, &slog)

	}
	data := model.JobSettlementLogsPagination{
		PageInfo: &pageInfo,
		Items:    settlementLogs,
	}
	return &data, nil
}

// ConfirmSettlement 确认结算
func (r *mutationResolver) ConfirmSettlement(ctx context.Context, settlementId int) (bool, error) {
	user := ctx.Value("user").(user.SessionUser)
	uId := user.UniqueUserId
	bool, err := job.ConfirmSettlementLogs(uId, int64(settlementId))

	return bool, err
}

// UpdateJobSettlementLog 更新结算记录
func (r *mutationResolver) UpdateJobSettlementLog(ctx context.Context, data model.UpdateJobSettlementLogInput) (*model.JobSettlementLog, error) {

	u := ctx.Value("user").(user.SessionUser)

	//获取岗位信息
	w, err := work.GetByID(data.WorkID)
	if err != nil {
		return nil, err
	}

	if u.CompanyID != w.CompanyId {
		return nil, errors.New("权限有误：公司不匹配")
	}

	param := m.JobSettlementLog{
		OperatorUserID: gconv.String(u.UserID),
		File:           data.File,
		Amount:         data.Amount,
	}
	param.ID = data.ID

	log, err := job.UpdateSettlementLog(param)
	if err != nil {
		return nil, err
	}

	var slog model.JobSettlementLog
	gconv.Struct(log, &slog)
	slog.ID = log.ID
	slog.CreatedAt = int(log.CreatedAt.Unix())
	slog.UpdatedAt = int(log.UpdatedAt.Unix())
	return &slog, nil
}

//结算总表列表
func (r *queryResolver) JobSettlements(ctx context.Context, pageNumber int, pageItem int, workID *int, batchID *string, name *string, createdBeginAt *string, createdEndAt *string, isToBeConfirm *bool) (*model.JobSettlementsPagination, error) {

	w := ctx.Value("user").(user.SessionUser)
	where := "work.company_id = " + strconv.Itoa(int(w.CompanyID))

	//组装where条件
	if workID != nil && *workID != 0 {
		where += " and work_id = " + strconv.Itoa(*workID)
	}
	if batchID != nil && *batchID != "" {
		where += " and batch_id = " + *batchID
	}
	if name != nil && *name != "" {
		where += " and work.name like '%" + *name + "%'"
	}
	if createdBeginAt != nil && *createdBeginAt != "" && createdEndAt != nil && *createdEndAt != "" {
		where += " and job_settlements.created_at between '" + *createdBeginAt + "' and '" + *createdEndAt + "'"
	}

	//待确认
	having := "1 = 1"
	if isToBeConfirm != nil && *isToBeConfirm == true {
		having += " and un_finish > 0"
	}
	type TempJobSettlements struct {
		ID          int       `json:"id"`           //id
		BatchID     string    `json:"batch_id"`     //批次id
		WorkID      int       `json:"work_id"`      //岗位id
		Amount      float64   `json:"amount"`       //绩效总金额
		SettleCount int       `json:"settle_count"` //绩效总条数
		CreatedAt   time.Time `json:"createdAt"`    //
		Name        string    `json:"name"`         //岗位名
		Extend      string    `json:"extend"`       //扩展字段
		UnFinish    int       `json:"un_finish"`    //未完成
	}

	//获取总数
	var cc []*TempJobSettlements
	err := db.Get().Table("job_settlements").
		Select("count(if(job_settlement_log.status='0',true,null )) as un_finish").
		Where(where).
		Joins("left join work on job_settlements.work_id = work.id").
		Joins("left join job_settlement_log ON job_settlements.batch_id = job_settlement_log.settle_id").
		Group("job_settlement_log.settle_id").
		Having(having).
		Scan(&cc).Error //此处用count会导致数据不准确,因此使用scan
	if err != nil {
		return nil, err
	}

	//组装一个sql
	var settlements []*TempJobSettlements
	offset := (pageNumber - 1) * pageItem
	err = db.Get().Table("job_settlements").
		Select("job_settlements.id,job_settlements.work_id,job_settlements.batch_id,job_settlements.amount,job_settlements.created_at,job_settlements.settle_count, work.name as name,work.extend as extend,count(if(job_settlement_log.status='0',true,null )) as un_finish").
		Where(where).
		Joins("left join work on job_settlements.work_id = work.id").
		Joins("left join job_settlement_log ON job_settlements.batch_id = job_settlement_log.settle_id").
		Group("job_settlement_log.settle_id").
		Having(having).
		Order("un_finish desc,job_settlements.created_at desc").
		Offset(offset).Limit(pageItem).
		Scan(&settlements).Error

	if err != nil {
		return nil, err
	}

	//转成prisma model
	type T struct {
		CompanyName string `json:"companyName"`
	}
	var items []*model.JobSettlements
	for i, _ := range settlements {
		var temp *model.JobSettlements
		gconv.Struct(settlements[i], &temp)
		temp.CreatedAt = int(settlements[i].CreatedAt.Unix())

		var t1 T
		err = json.Unmarshal([]byte(settlements[i].Extend), &t1)
		temp.CompanyName = t1.CompanyName
		items = append(items, temp)
	}
	count := len(cc)
	pagination := model.JobSettlementsPagination{
		PageInfo: &model.PageInfo{
			TotalPage: int(math.Ceil(float64(count) / float64(pageItem))),
			TotalItem: count,
		},
		Items: items,
	}

	return &pagination, err
}
