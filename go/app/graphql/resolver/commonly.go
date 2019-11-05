package resolver

import (
	"context"
	"crypto/md5"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"iQuest/app/Api"
	"iQuest/app/constant"
	"iQuest/app/graphql/model"
	"iQuest/app/graphql/prisma"
	"iQuest/app/graphql/subscription"
	gormModel "iQuest/app/model"
	m_user "iQuest/app/model/user"
	service "iQuest/app/service/job"
	"iQuest/config"
	"iQuest/db"
	"iQuest/library/response"
	"io"
	"math"
	"math/rand"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

type Data struct {
	Name        string `json:"name"`
	Mobile      string `json:"mobile"`
	IdCardNo    string `json:"id_card_no"`
	CompanyId   int    `json:"company_id"`
	CompanyName string `json:"company_name"`
	Id          string `json:"id"`
	FileHash    string `json:"file_hash"`
}

/**
创建常用人员
*/
func (r *mutationResolver) CreateCommonlyUsedPersonnel(ctx context.Context, newCommonlyUsedPersonnelData model.NewCommonlyUsedPersonnelInput) (*model.CommonlyUsedPersonnel, error) {
	user := ctx.Value("user").(m_user.SessionUser)

	//用户认证与创建接口
	service_user, err := Api.UserVerifiedAndCreate(m_user.UserServiceVerfiedInput{
		CompanyId:   int(user.CompanyID),
		CompanyName: newCommonlyUsedPersonnelData.CompanyName,
		RealName:    newCommonlyUsedPersonnelData.Name,
		IdCardNo:    newCommonlyUsedPersonnelData.CardNo,
		MobilePhone: newCommonlyUsedPersonnelData.Mobile,
		BankCardNo:  *newCommonlyUsedPersonnelData.BankNo,
	})

	if err != nil {
		return nil, err
	}

	service_user_id := strconv.Itoa(int(service_user.UserId))
	if service_user_id == "0" {
		return nil, errors.New(service_user.Msg)
	}

	ver, err := Api.FindRealnameInfoByUserid(int64(service_user.UserId))

	if err != nil {

		return nil, err
	}

	if ver.Data.State != 1 {
		return nil, errors.New("用户认证失败")
	}

	user_id, _ := r.Prisma.CommonlyUsedPersonnels(&prisma.CommonlyUsedPersonnelsParams{
		Where: &prisma.CommonlyUsedPersonnelWhereInput{
			UserId:    &service_user_id,
			CompanyId: &user.CompanyID,
			DeletedAt: prisma.Int32(0),
		},
	}).Exec(ctx)
	if user_id != nil {
		return nil, errors.New("用户已存在")
	}

	commonly, err := r.Prisma.CreateCommonlyUsedPersonnel(prisma.CommonlyUsedPersonnelCreateInput{
		CompanyId: user.CompanyID,
		UserId:    service_user_id,
		Name:      newCommonlyUsedPersonnelData.Name,
		Avatar:    newCommonlyUsedPersonnelData.Avatar,
		Mobile:    newCommonlyUsedPersonnelData.Mobile,
		AppId:     newCommonlyUsedPersonnelData.AppID,
		CardNo:    newCommonlyUsedPersonnelData.CardNo,
		BankNo:    *newCommonlyUsedPersonnelData.BankNo,
		Education: newCommonlyUsedPersonnelData.Education,
		Address:   newCommonlyUsedPersonnelData.Address,
		Remark:    newCommonlyUsedPersonnelData.Remark,
		DeletedAt: prisma.Int32(0),
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	var user_info gormModel.UserWeChatAuthorize
	if err := db.Get().Model(&gormModel.UserWeChatAuthorize{}).Unscoped().Where("user_id = ? ", service_user_id).First(&user_info).Error; err != nil {
		weChat := &gormModel.UserWeChatAuthorize{
			UserId: service_user_id,
			Mobile: newCommonlyUsedPersonnelData.Mobile,
		}

		if err := db.Get().Create(&weChat).Error; err != nil {
			return nil, err
		}
	}

	t := model.CommonlyUsedPersonnel{
		ID: int(commonly.ID),
	}
	return &t, err
}

/**
删除常用人员
*/
func (r *mutationResolver) DeleteCommonlyUsedPersonnel(ctx context.Context, id int) (bool, error) {
	c_id := int32(id)
	timeUnix := int32(time.Now().Unix())
	_, err := r.Prisma.UpdateCommonlyUsedPersonnel(prisma.CommonlyUsedPersonnelUpdateParams{
		prisma.CommonlyUsedPersonnelUpdateInput{
			DeletedAt: &timeUnix,
		},
		prisma.CommonlyUsedPersonnelWhereUniqueInput{
			ID: &c_id,
		},
	}).Exec(ctx)

	if err != nil {
		return false, err
	}
	return true, nil
}

/**
更新常用人员
*/
func (r *mutationResolver) UpdateCommonlyUsedPersonnel(ctx context.Context, id *int, updateCommonlyUsedPersonnelData model.UpdateCommonlyUsedPersonnelInput) (*model.CommonlyUsedPersonnel, error) {
	c_id := int32(*id)

	_, err := r.Prisma.UpdateCommonlyUsedPersonnel(prisma.CommonlyUsedPersonnelUpdateParams{
		prisma.CommonlyUsedPersonnelUpdateInput{
			Avatar:    updateCommonlyUsedPersonnelData.Avatar,
			Mobile:    updateCommonlyUsedPersonnelData.Mobile,
			Education: updateCommonlyUsedPersonnelData.Education,
			Address:   updateCommonlyUsedPersonnelData.Address,
			Remark:    updateCommonlyUsedPersonnelData.Remark,
			BankNo:    updateCommonlyUsedPersonnelData.BankNo,
		},
		prisma.CommonlyUsedPersonnelWhereUniqueInput{
			ID: &c_id,
		},
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	t := model.CommonlyUsedPersonnel{
		ID: *id,
	}

	return &t, nil
}

/**
常用人员详情
*/
func (r *queryResolver) CommonlyUsedPersonnelDetail(ctx context.Context, id int) (*model.CommonlyUsedPersonnelInfo, error) {

	user := ctx.Value("user").(m_user.SessionUser)
	commonly_id := int32(id)
	detail, err := r.Prisma.CommonlyUsedPersonnel(prisma.CommonlyUsedPersonnelWhereUniqueInput{
		ID: &commonly_id,
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	//获取签约信息接口
	sign_info, err := Api.GetSignInfo(detail.CardNo)

	if err != nil {
		return nil, err
	}

	var party_list []*model.PartyA

	for i := 0; i < len(sign_info.Data); i++ {
		//过滤非本公司的签约信息
		if sign_info.Data[i].CompanyId == strconv.Itoa(int(user.CompanyID)) {
			api_sign_time := strconv.Itoa(int(sign_info.Data[i].SignTime))
			var partyA model.PartyA
			if sign_info.Data[i].IsGroup == true {
				server_name := ""
				for _, v := range sign_info.Data[i].GroupInfo {
					server_name += v.ServerName + "/"
				}
				content := server_name[0 : len(server_name)-1]
				partyA = model.PartyA{
					CompanyName: &content,
					SignTime:    &api_sign_time,
				}
			} else {
				partyA = model.PartyA{
					CompanyName: &sign_info.Data[i].ServerName,
					SignTime:    &api_sign_time,
				}
			}

			party_list = append(party_list, &partyA)
		}
	}

	sign_time := (*int)(unsafe.Pointer(detail.SigningTime))

	t := model.CommonlyUsedPersonnelInfo{
		ID:        id,
		Avatar:    detail.Avatar,
		Name:      detail.Name,
		Mobile:    detail.Mobile,
		CardNo:    detail.CardNo,
		BankNo:    &detail.BankNo,
		Education: detail.Education,
		Address:   detail.Address,
		Remark:    detail.Remark,
		SignTime:  sign_time,
		PartyA:    party_list,
		PartyB:    &detail.Name,
	}

	return &t, nil
}

/**
常用人员列表
*/

func (r *queryResolver) CommonlyUsedPersonnelLists(ctx context.Context, page *int, pageSize *int, id *int, name *string, mobile *string, createdAtBegin *int, createdAtEnd *int) (*model.CommonlyUsedPersonnelPagination, error) {

	user := ctx.Value("user").(m_user.SessionUser)
	companyId := user.CompanyID

	if *page < 0 || *pageSize < 0 {
		return nil, errors.New(response.ParamErrorMsg)
	}

	skip := int32((*page - 1) * (*pageSize))
	page_size := int32(*pageSize)
	order_by := prisma.CommonlyUsedPersonnelOrderByInputIDDesc

	c_id := (*int32)(unsafe.Pointer(id))

	time_begin := (*int64)(unsafe.Pointer(createdAtBegin))
	time_end := (*int64)(unsafe.Pointer(createdAtEnd))

	var created_at_begin *string
	if time_begin != nil {
		time_begin_utc := time.Unix(*time_begin, 0).UTC().Format(time.RFC3339)
		created_at_begin = &time_begin_utc
	}
	var created_at_end *string
	if time_end != nil {
		time_end_utc := time.Unix(*time_end, 0).UTC().Format(time.RFC3339)
		created_at_end = &time_end_utc
	}

	total, err := r.Prisma.CommonlyUsedPersonnelsConnection(&prisma.CommonlyUsedPersonnelsConnectionParams{
		Where: &prisma.CommonlyUsedPersonnelWhereInput{
			ID:           c_id,
			CompanyId:    &companyId,
			Name:         name,
			Mobile:       mobile,
			CreatedAtGte: created_at_begin,
			CreatedAtLte: created_at_end,
			DeletedAt:    prisma.Int32(0),
		},
	}).Aggregate(ctx)

	commonly_list, err := r.Prisma.CommonlyUsedPersonnels(&prisma.CommonlyUsedPersonnelsParams{
		Where: &prisma.CommonlyUsedPersonnelWhereInput{
			ID:           c_id,
			CompanyId:    &companyId,
			Name:         name,
			Mobile:       mobile,
			CreatedAtGte: created_at_begin,
			CreatedAtLte: created_at_end,
			DeletedAt:    prisma.Int32(0),
		},
		OrderBy: &order_by,
		Skip:    &skip,
		First:   &page_size,
	}).Exec(ctx)
	if err != nil {
		return nil, err
	}

	var commonly_all []*model.CommonlyUsedPersonnelList
	var user_ids []string
	var work_ids []int32

	//用户id数组
	for i := 0; i < len(commonly_list); i++ {
		user_ids = append(user_ids, commonly_list[i].UserId)
	}

	progress_list := []int32{1, 2, 9}
	//批量查询用户参加岗位的id
	positions, err := r.Prisma.JobMembers(&prisma.JobMembersParams{
		Where: &prisma.JobMemberWhereInput{
			ParticipantIdIn: user_ids,
			CompanyId:       &companyId,
			ProgressNotIn:   progress_list,
		},
	}).Exec(ctx)

	if err != nil {
		return nil, err
	}

	//用户与岗位id关系集合
	user_work_relation := make(map[string][]int32)
	for i := 0; i < len(positions); i++ {
		if !InArray(positions[i].WorkId, work_ids) {
			work_ids = append(work_ids, positions[i].WorkId)
		}
		position_user_id := positions[i].ParticipantId

		for j := 0; j < len(user_ids); j++ {
			if *position_user_id == user_ids[j] {
				user_work_relation[*position_user_id] = append(user_work_relation[*position_user_id], positions[i].WorkId)
			}
		}

	}

	status_list := []int32{1, 3}
	//通过岗位id批量查询岗位信息
	works, err := r.Prisma.Works(&prisma.WorksParams{
		Where: &prisma.WorkWhereInput{
			IDIn:     work_ids,
			StatusIn: status_list,
		},
	}).Exec(ctx)

	for i := 0; i < len(commonly_list); i++ {
		//岗位数量统计
		commonly_work_name := 0
		for _, v := range user_work_relation[commonly_list[i].UserId] {
			for g := 0; g < len(works); g++ {

				if v == works[g].ID {
					commonly_work_name += 1
				}
			}
		}

		//发放流水统计接口
		res, err := Api.UserFlowCount(commonly_list[i].UserId, int(companyId))
		if err != nil {
			return nil, err
		}

		commonly_sign_time := (*int)(unsafe.Pointer(commonly_list[i].SigningTime))

		loc, _ := time.LoadLocation("Local")
		//time.Unix(*time_begin, 0).UTC().Format(time.RFC3339)
		theTime, _ := time.ParseInLocation(time.RFC3339, commonly_list[i].CreatedAt, loc)
		sr := int(theTime.Unix())

		commonly_all = append(commonly_all, &model.CommonlyUsedPersonnelList{
			ID:          int(commonly_list[i].ID),
			UserID:      &commonly_list[i].UserId,
			Name:        &commonly_list[i].Name,
			Mobile:      &commonly_list[i].Mobile,
			Address:     commonly_list[i].Address,
			Position:    &commonly_work_name,
			Achievement: &res.Data,
			Remark:      commonly_list[i].Remark,
			SignTime:    commonly_sign_time,
			CreatedAt:   &sr,
		})
	}
	total_page := 0
	if int(total.Count) > 0 {
		total_page = int(math.Ceil(float64(total.Count) / float64(*pageSize)))
	}

	t := model.CommonlyUsedPersonnelPagination{
		TotalItem: int(total.Count),
		TotalPage: total_page,
		Items:     commonly_all,
	}

	return &t, nil

}

/**
常用人员详情绩效接口
*/
func (r *queryResolver) AchievementsDetail(ctx context.Context, id string, page *int, pageSize *int) (*model.AchievementsPagination, error) {

	user := ctx.Value("user").(m_user.SessionUser)
	companyId := user.CompanyID
	if *page < 0 || *pageSize < 0 {
		return nil, errors.New(response.ParamErrorMsg)
	}

	skip := int32(*page)
	page_size := int32(*pageSize)

	job_member_list, err := service.GetJobMemberJoinWork(service.JobCondition{
		PageNum:   skip,
		PageSize:  page_size,
		UserId:    id,
		CompanyId: int(companyId),
	})

	if err != nil {
		return nil, err
	}

	var post_ids []string
	var work_ids []int
	var ach_list []*model.AchievementsInfo
	job_arr := job_member_list["item"].([]service.JobMemberList)

	total_item := int(job_member_list["total"].(int32))
	total_page := int(job_member_list["last_page"].(int32))
	if total_item < 0 {
		t := model.AchievementsPagination{
			TotalItem: 0,
			TotalPage: total_page,
			Items:     ach_list,
		}
		return &t, nil
	}

	for i := 0; i < len(job_arr); i++ {
		work_id := strconv.Itoa(int(job_arr[i].WorkId))
		post_ids = append(post_ids, work_id)
		work_ids = append(work_ids, job_arr[i].WorkId)
	}

	res, err := Api.UserPositionFlowCount(id, int64(companyId), post_ids)

	if err != nil {
		return nil, err
	}

	//获取服务类型接口
	service_data, err := Api.ServiceName()
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(job_arr); i++ {
		//服务类型过滤
		position_type := ""

		for j := 0; j < len(service_data.Data); j++ {
			if job_arr[i].ServiceTypeId == service_data.Data[j].ServiceId {
				position_type = service_data.Data[j].ServiceName
			}
		}

		var postId []string
		work_id := strconv.Itoa(job_arr[i].WorkId)
		postId = append(postId, work_id)

		//发放次数统计
		count := 0
		for j := 0; j < len(res.Data); j++ {
			if res.Data[j].PostId == work_id {
				count = res.Data[j].Count
			}
		}

		created_at := int(job_arr[i].CreatedAt.Unix())
		ach_list = append(ach_list, &model.AchievementsInfo{
			ID:            job_arr[i].Id,
			WorkID:        &job_arr[i].WorkId,
			PositionTitle: &job_arr[i].Name,
			PositionType:  &position_type,
			Status:        &job_arr[i].Progress,
			GiveTimes:     &count,
			JoinTimes:     &created_at,
		})
	}

	//total, err := r.Prisma.JobMembersConnection(&prisma.JobMembersConnectionParams{
	//	Where: &prisma.JobMemberWhereInput{
	//		ParticipantId: &id,
	//	},
	//}).Aggregate(ctx)
	//
	//
	//
	//if err != nil {
	//	return nil, err
	//}
	//
	//achieves, err := r.Prisma.JobMembers(&prisma.JobMembersParams{
	//	Where: &prisma.JobMemberWhereInput{
	//		ParticipantId: &id,
	//	},
	//	Skip:  &skip,
	//	First: &page_size,
	//}).Exec(ctx)
	//
	//if err != nil {
	//	return nil, err
	//}
	//
	//var ach_list []*model.AchievementsInfo
	//for i := 0; i < len(achieves); i++ {
	//	work_id := strconv.Itoa(int(achieves[i].WorkId))
	//	post_ids = append(post_ids, work_id)
	//	work_ids = append(work_ids, achieves[i].WorkId)
	//}
	//
	////根据岗位以及用户id获取发放信息接口
	//res, err := Api.UserPositionFlowCount(id, int64(companyId), post_ids)
	//
	//if err != nil {
	//	return nil, err
	//}
	//
	//status := []int32{1, 3}
	////获取岗位信息
	//works_info, err := r.Prisma.Works(&prisma.WorksParams{
	//	Where: &prisma.WorkWhereInput{
	//		IDIn:     work_ids,
	//		StatusIn: status,
	//	},
	//}).Exec(ctx)
	//
	//if err != nil {
	//	return nil, err
	//}
	//
	////获取服务类型接口
	//service_data, err := Api.ServiceName()
	//
	//if err != nil {
	//	return nil, err
	//}
	//
	//for i := 0; i < len(achieves); i++ {
	//	//服务类型过滤
	//	position_type := ""
	//	for h := 0; h < len(works_info); h++ {
	//		for j := 0; j < len(service_data.Data); j++ {
	//			if (int(works_info[h].ServiceTypeId) == service_data.Data[j].ServiceId) && works_info[h].ID == achieves[i].WorkId {
	//				position_type = service_data.Data[j].ServiceName
	//			}
	//		}
	//	}
	//
	//	var postId []string
	//	work_id := strconv.Itoa(int(achieves[i].WorkId))
	//	postId = append(postId, work_id)
	//
	//	//发放次数统计
	//	count := 0
	//	for j := 0; j < len(res.Data); j++ {
	//		if res.Data[j].PostId == work_id {
	//			count = res.Data[j].Count
	//		}
	//	}
	//
	//	work, err := r.Prisma.Works(&prisma.WorksParams{
	//		Where: &prisma.WorkWhereInput{
	//			ID: &achieves[i].WorkId,
	//		},
	//	}).Exec(ctx)
	//
	//	if err != nil {
	//		return nil, err
	//	}
	//	output_work_id := int(achieves[i].WorkId)
	//	progress := int(achieves[i].Progress)
	//	created_at := job.DateTimeToTimestamp(achieves[i].CreatedAt)
	//	ach_list = append(ach_list, &model.AchievementsInfo{
	//		ID:            int(achieves[i].ID),
	//		WorkID:        &output_work_id,
	//		PositionTitle: &work[0].Name,
	//		PositionType:  &position_type,
	//		Status:        &progress,
	//		GiveTimes:     &count,
	//		JoinTimes:     created_at,
	//	})
	//}

	t := model.AchievementsPagination{
		TotalItem: total_item,
		TotalPage: total_page,
		Items:     ach_list,
	}

	return &t, nil
}

func InArray(s interface{}, d []int32) bool {
	for _, v := range d {
		if s == v {
			return true
		}
	}
	return false
}

/**
导入用户
*/
func (r *mutationResolver) ImportUser(ctx context.Context, req []*model.UploadFile, companyName string) (*model.ImportStatus, error) {

	var s model.ImportStatus
	user := ctx.Value("user").(m_user.SessionUser)

	file_name := req[0].File.Filename
	upload_path := config.Viper.GetString("UPLOAD_PATH") + strconv.FormatInt(time.Now().Unix(), 10) + file_name
	ext := path.Ext(file_name)

	if ext != ".xlsx" {
		s = model.ImportStatus{
			ID:     "",
			Status: false,
		}
		return &s, nil
	}

	fW, err := os.Create(upload_path)
	if err != nil {
		return nil, err
	}
	defer fW.Close()

	_, err = io.Copy(fW, req[0].File.File)
	if err != nil {
		return nil, err
	}

	f, err := excelize.OpenFile(upload_path)
	if err != nil {
		return nil, err
	}

	file_hash, err := Md5File(f.Path)
	if err != nil {
		return nil, err
	}
	id := randString(8)
	file_hash = file_hash + id + strconv.Itoa(int(time.Now().Unix()))
	queue := service.QueueConn.OpenQueue(constant.ImportUser)

	rows, err := f.GetRows(f.GetSheetName(1))

	total_data := 0
	for i, row := range rows {
		// 去掉第一行，第一行是表头
		if i == 0 {
			continue
		}
		total_data += 1
		var data Data

		for j, colCell := range row {

			colCell = strings.Replace(colCell, " ", "", -1)
			if j == 0 && colCell == "Null" {
				continue
			}
			if j == 0 && colCell != "Null" {
				if IsChineseChar(colCell) != true {
					return nil, errors.New("无法导入，存在不合法参数" + colCell)
				}
				data.Name = colCell

			}
			if j == 1 {
				data.Mobile = colCell
			}
			if j == 2 {
				user_id, _ := r.Prisma.CommonlyUsedPersonnels(&prisma.CommonlyUsedPersonnelsParams{
					Where: &prisma.CommonlyUsedPersonnelWhereInput{
						CardNo:    &colCell,
						CompanyId: &user.CompanyID,
						DeletedAt: prisma.Int32(0),
					},
				}).Exec(context.Background())

				if user_id != nil {
					return nil, errors.New("无法导入，身份证" + colCell + "已存在")
				}

				data.IdCardNo = colCell
			}
			data.CompanyId = int(user.CompanyID)
			data.CompanyName = companyName
			data.Id = id
			data.FileHash = file_hash
		}

		payloadBytes, err := json.Marshal(data)
		if err != nil {
			return nil, errors.New("导入用户队列入队出错")
		}

		//队列
		import_user := queue.PublishBytes(payloadBytes)
		if !import_user {
			//TODO 不成功处理放弃治疗
		}
	}

	s = model.ImportStatus{
		ID:     file_hash,
		Status: true,
		Total:  total_data,
	}
	return &s, nil
}

func SHA256File(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	h := sha256.New()
	_, err = io.Copy(h, file)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func Md5File(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}

	defer file.Close()

	h := md5.New()
	_, err = io.Copy(h, file)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func IsChineseChar(str string) bool {
	for _, r := range str {
		if !(regexp.MustCompile("^[\u4e00-\u9fa5]+").MatchString(string(r))) {
			return false
		}
	}
	return true
}

func (r *subscriptionResolver) ImportError(ctx context.Context, roomName string) (<-chan *model.Message, error) {

	subscription.Server.MU.Lock()
	room := subscription.Server.Rooms[roomName]
	if room == nil {
		room = &subscription.ChatRoom{
			Name: roomName,
		}
		subscription.Server.Rooms[roomName] = room
	}
	subscription.Server.MU.Unlock()

	//id := randString(8)
	events := make(chan *model.Message, 1)

	go func() {
		<-ctx.Done()
		subscription.Server.MU.Lock()
		delete(subscription.Server.Rooms, roomName)
		subscription.Server.MU.Unlock()
	}()

	return events, nil
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
