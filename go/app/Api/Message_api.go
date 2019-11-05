package Api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"iQuest/app/eureka"
	"iQuest/config"
	"iQuest/library/response"
	"log"
)

type MsgResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type GetSendMsgContentData struct {
	WorkId         int32  `json:"WorkId"`
	SendId         string `json:"SendId"`
	ReceiverId     string `json:"ReceiverId"`
	SendType       int    `json:"SendType"`
	WorkProgressId int32  `json:"WorkProgressId"`
	TaskMemberId int32  `json:"TaskMemberId"`
	CompanyName    string `json:"CompanyName"`
	WorkName       string `json:"WorkName"`
	ParamField          string  `json:"ParamField"`
}

var SEND_MESSAGE_ARR = [][]string{
	{
		"岗位申请已被通过",
		"您申请%s的岗位：%s 已被通过，点击消息查看详情",
	},
	{
		"岗位申请未能通过",
		"很遗憾，您申请%s的岗位：%s 未能通过，点击消息查看详情",
	},
	{
		"收到新的岗位邀请",
		"您收到一条来自%s的岗位邀请，点击消息前往处理",
	},
	{
		"收到新的结算单待确认",
		"您参加的岗位 %s 有一条新的结算单待确认，请点击本消息查看详情并确认，延期未结算可能会对您的薪酬发放造成影响哦。",
	},
}

//发送消息
func SendMessage(Data GetSendMsgContentData) (MsgResp, error) {
	var t MsgResp
	urls := config.Viper.GetString("MESSAGE_URL") + "/v1/message/sendMessage"
	log.Println(urls)
	//生成发送内容
	var title string
	var content string
	switch Data.SendType {
	case 1:
		title = SEND_MESSAGE_ARR[Data.SendType-1][0]
		content = fmt.Sprintf(SEND_MESSAGE_ARR[Data.SendType-1][1], Data.CompanyName, Data.WorkName)
	case 2:
		title = SEND_MESSAGE_ARR[Data.SendType-1][0]
		content = fmt.Sprintf(SEND_MESSAGE_ARR[Data.SendType-1][1], Data.CompanyName, Data.WorkName)
	case 3:
		title = SEND_MESSAGE_ARR[Data.SendType-1][0]
		content = fmt.Sprintf(SEND_MESSAGE_ARR[Data.SendType-1][1], Data.CompanyName)
	case 4:
		title = SEND_MESSAGE_ARR[Data.SendType-1][0]
		content = fmt.Sprintf(SEND_MESSAGE_ARR[Data.SendType-1][1], Data.WorkName)
	}

	//构造参数
	post_data := map[string]interface{}{
		"work_id":          Data.WorkId,
		"sender_id":        Data.SendId,
		"receiver_id":      Data.ReceiverId,
		"title":            title,
		"content":          content,
		"message_type":     Data.SendType,
		"work_progress_id": Data.WorkProgressId,
		"task_member_id": Data.TaskMemberId,
		"type":             0,
		"param_field":     Data.ParamField,
	}
	bytesData, err := json.Marshal(post_data)

	resp, err := eureka.Post(urls, bytes.NewReader(bytesData))
	log.Println(urls)
	if err != nil {
		t = MsgResp{
			Code: response.Error,
		}
		log.Println(err)
		return t, err
	}

	err = json.NewDecoder(resp.Body).Decode(&t)
	log.Println("%v",resp.Body)
	log.Println(t)
	if t.Code != 0 {
		t = MsgResp{
			Code: response.Error,
		}
		return t, nil
	}

	return t, nil

}
