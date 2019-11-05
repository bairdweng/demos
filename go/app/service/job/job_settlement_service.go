package job

import (
	"encoding/json"
	"errors"
	"fmt"
	"iQuest/app/Api"
	"iQuest/app/constant"
	"iQuest/app/model"
	"iQuest/config"
	"iQuest/db"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
)

// SettlementLogsCondition 查询条件
type SettlementLogsCondition struct {
	PageNum  int
	PageSize int
	WorkID   int
	UserID   string
}

// GetSettlementLogs 获取结算记录
func GetSettlementLogs(condition SettlementLogsCondition) ([]*model.JobSettlementLog, int64, error) {

	logs := []*model.JobSettlementLog{}
	offset := (condition.PageNum - 1) * condition.PageSize
	db.Get().Where(model.JobSettlementLog{WorkID: condition.WorkID, UserID: condition.UserID}).Offset(offset).Limit(condition.PageSize).Find(&logs)

	var total int64
	db.Get().Model(&model.JobSettlementLog{}).Where(model.JobSettlementLog{WorkID: condition.WorkID, UserID: condition.UserID}).Count(&total)

	return logs, total, nil
}

// CreateSettlementLogs 创建结算记录
func CreateSettlementLogs(logs []*model.JobSettlementLog) ([]*model.JobSettlementLog, error) {
	var settlementLogs []*model.JobSettlementLog
	tx := db.Get().Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	allAmount := 0.0
	for _, amount := range logs {
		//获取总发放金额
		allAmount += amount.Amount

	}
	insertData := model.JobSettlements{
		WorkID:logs[0].WorkID,
		Amount:allAmount,
		SettleCount:len(logs),
	}
	if err := tx.Create(&insertData).Error;err != nil {
		tx.Rollback()
		return nil, err
	}
	for _, log := range logs {
		log.SettleID = insertData.BatchID
		if err := tx.Create(&log).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		settlementLogs = append(settlementLogs, log)

		//发送短信
		mobile, err := FindMobile(log.UserID)
		if err != nil {
			println("找不到手机号:" + err.Error())
			continue
		}

		var w model.Work
		err = db.Get().Unscoped().Where("id = ?", log.WorkID).First(&w).Error
		if err != nil {
			println("找不到岗位信息:" + err.Error())
			continue
		}

		if mobile != "" {
			newContent := fmt.Sprintf(constant.SETTLEMENT_SMS_CONTENT, w.Name)
			go Api.SendSms(mobile, newContent,config.Viper.GetString("SMS_XINNIAO_ID"))
		}

		//发送站内信
		var extendField = map[string]string{}
		extendField["settlement_id"] = strconv.Itoa(log.ID)
		extend, _ := json.Marshal(extendField)
		extendString := string(extend)

		newSendParam := Api.GetSendMsgContentData{
			WorkId: int32(log.WorkID),
			//SendId:      sendUserId,
			ReceiverId:   log.UserID,
			SendType:     4,
			TaskMemberId: int32(log.MemberID),
			WorkName:     w.Name,
			ParamField:   extendString,
		}
		go Api.SendMessage(newSendParam)
	}

	tx.Commit()
	return settlementLogs, nil
}

//确认结算
func ConfirmSettlementLogs(userId string, settlementId int64) (bool, error) {
	//查询从属关系
	var settlementLog model.JobSettlementLog
	err := db.Get().Unscoped().Model(&model.JobSettlementLog{}).
		Where("id = ? and user_id = ?", settlementId, userId).
		First(&settlementLog).Error

	if err != nil || err == gorm.ErrRecordNotFound {
		return false, err
	}
	//确认
	if err := db.Get().Unscoped().Model(&model.JobSettlementLog{}).
		Where("id = ?", settlementId).
		Update(&model.JobSettlementLog{Status: constant.CONFIRM_STATUS, ConfirmAt: int(time.Now().Unix())}).Error; err != nil {
		return false, err
	}
	return true, nil
}

//GetSettlementLogsByIds 获取结算记录
func GetSettlementLogsByIds(workID int, ids []int) ([]*model.JobSettlementLog, error) {

	var logs []*model.JobSettlementLog
	db.Get().Where("work_id = ?", workID).Where("id in (?)", ids).Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Where("deleted_at = 0 or deleted_at is null").Unscoped()
	}).Find(&logs)

	return logs, nil
}

func FindMobile(userID string) (string, error) {
	//查询手机号
	var user model.UserWeChatAuthorize
	err := db.Get().Model(&model.UserWeChatAuthorize{}).Unscoped().
		Where("user_id = ? ", userID).First(&user).Error

	if user.Mobile != "" {
		return user.Mobile, nil
	}

	//查询接口
	userIdInt, _ := strconv.Atoi(userID)
	ver, err := Api.FindRealnameInfoByUserid(int64(userIdInt))
	if err != nil || ver.Code != 0 {
		println("接口找不到手机号")
		return "", err
	}

	if ver.Data.MobilePhone == "" {
		return "", err
	}
	return ver.Data.MobilePhone, err
}

// UpdateSettlementLog 更新结算记录
func UpdateSettlementLog(param model.JobSettlementLog) (*model.JobSettlementLog, error) {

	var log model.JobSettlementLog
	if db.Get().First(&log, param.ID).RecordNotFound() {
		return nil, errors.New("数据有误：无结算记录")
	}
	if log.Status == constant.CONFIRM_STATUS {
		return nil, errors.New("操作有误：结算已被确认")
	}
	log.File = param.File
	log.Amount = param.Amount
	log.OperatorUserID = param.OperatorUserID
	if err := db.Get().Save(&log).Error; err != nil {
		return nil, err
	}
	return &log, nil
}
