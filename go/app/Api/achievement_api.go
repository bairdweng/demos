package Api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"iQuest/app/eureka"
	"iQuest/config"
	"log"
	"net/http"
)

type achievementResp struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

func UserJoinJobCallBack(userId int64, workId int32) (achievementResp, error) {
	var t achievementResp
	//urls := config.Viper.GetString("PASS_PORT_URL") + "/user/" + fmt.Sprintf("%d", id)
	urls := config.Viper.GetString("POST_NOTIFY_URL") + "/post-notify/callback"

	post_data := map[string]interface{}{
		"userId": userId,
		"jobId":  workId,
	}
	bytesData, err := json.Marshal(post_data)

	resp, err := eureka.Post(urls, bytes.NewReader(bytesData))

	if err != nil {
		t = achievementResp{
			Code: "1111",
		}
		return t, err
	}

	err = json.NewDecoder(resp.Body).Decode(&t)
	if t.Code != "0000" {
		t = achievementResp{
			Code: "1111",
		}
		return t, nil
	}
	return t, nil

}

//岗位流水条数
type achievementCountResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data int32  `json:"data"`
}

func GetJobAchievementCount(companyId int32, workId int32) (achievementCountResp, error) {
	var t achievementCountResp
	urls := config.Viper.GetString("POST_EVIDENCE_URL") + "/post-evidence/post/flow/count"
	log.Print(urls)
	post_data := map[string]interface{}{
		"custCompanyId": companyId,
		"postId":        workId,
	}
	log.Printf("参数company_id:%s,post_id:%s", post_data["custCompanyId"], post_data["postId"])
	bytesData, err := json.Marshal(post_data)

	resp, err := eureka.Post(urls,
		bytes.NewReader(bytesData))

	if err != nil {
		log.Printf("%v", err)
		return t, err
	}
	err = json.NewDecoder(resp.Body).Decode(&t)
	log.Printf("参数code:%d,msg:%s,data:%v", t.Code, t.Msg, t.Data)
	if t.Code != http.StatusOK {
		return t, err
	}
	return t, err
}

//企业发放金额/次数查询
type companyProvideAmountResp struct {
	Code int                    `json:"code"`
	Msg  string                 `json:"msg"`
	Data []companyProvideAmount `json:"Data"`
}

type companyProvideAmount struct {
	Amount   float64 `json:"amount"`
	Billdate string  `json:"billdate"`
	Count    int32   `json:"count"`
}

func GetCompanyProvideAmount(companyId int32, startAt string, endAt string) (companyProvideAmountResp, error) {
	var t companyProvideAmountResp
	urls := config.Viper.GetString("POST_EVIDENCE_URL") + "/post-evidence/cust-company/pay"
	log.Print(urls)
	post_data := map[string]interface{}{
		"custCompanyId": companyId,
		"startAt":       startAt,
		"endAt":         endAt,
	}
	bytesData, err := json.Marshal(post_data)

	resp, err := eureka.Post(urls, bytes.NewReader(bytesData))

	if err != nil {
		return t, err
	}
	err = json.NewDecoder(resp.Body).Decode(&t)
	log.Print(t)
	if t.Code != http.StatusOK {
		return t, err
	}
	return t, err
}

//企业签约记录查询
type companySignDataResp struct {
	Code int               `json:"code"`
	Msg  string            `json:"msg"`
	Data []companySignData `json:"Data"`
}

type companySignData struct {
	ServiceCompanyId   string `json:"serviceCompanyId"`
	ServiceCompanyName string `json:"serviceCompanyName"`
	ServiceTypeName    string `json:"serviceTypeName"`
}

func GetCompanySignData(companyId int32) (companySignDataResp, error) {
	var t companySignDataResp
	urls := config.Viper.GetString("POST_EVIDENCE_URL") + "/post-evidence/cust-company/sign/" + fmt.Sprintf("%d", companyId)

	resp, err := eureka.Get(urls)

	if err != nil {
		return t, err
	}
	err = json.NewDecoder(resp.Body).Decode(&t)
	if t.Code != http.StatusOK {
		return t, err
	}
	return t, err
}

type UserFlowCountResp struct {
	Code int    `json:"code"`
	Data int    `json:"data"`
	Msg  string `json:"msg"`
}

func UserFlowCount(userId string, companyId int) (*UserFlowCountResp, error) {
	post_data := map[string]interface{}{
		"custCompanyId": companyId,
		"userId":        userId,
	}
	bytesData, err := json.Marshal(post_data)

	resp, err := eureka.Post(config.Viper.GetString("USER_FLOW_COUNT"), bytes.NewReader(bytesData))

	var t UserFlowCountResp
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(resp.Body).Decode(&t)
	return &t, nil

}

type UserPositionResp struct {
	Code int              `json:"code"`
	Msg  string           `json:"msg"`
	Data []UserFlowDetail `json:"data"`
}

type UserFlowDetail struct {
	PostId string `json:"postId"`
	Count  int    `json:"count"`
}

func UserPositionFlowCount(userId string, companyId int64, postId []string) (*UserPositionResp, error) {
	post_data := map[string]interface{}{
		"custCompanyId": companyId,
		"userId":        userId,
		"postIds":       postId,
	}
	bytesData, err := json.Marshal(post_data)

	resp, err := eureka.Post(config.Viper.GetString("USER_POSITION_FLOW_COUNT"), bytes.NewReader(bytesData))

	var t UserPositionResp
	if err != nil {

		return nil, err

	}

	err = json.NewDecoder(resp.Body).Decode(&t)

	return &t, nil

}

type ServiceNameResp struct {
	Code int
	Data []ServiceInfoResp
}

type ServiceInfoResp struct {
	ServiceId      int    `json:"serviceId"`
	ServiceName    string `json:"serviceName"`
	ServiceContent string `json:"serviceContent"`
}

func ServiceName() (*ServiceNameResp, error) {
	var t ServiceNameResp
	url := config.Viper.GetString("SERVICE_NAME")

	resp, err := eureka.Get(url)
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(resp.Body).Decode(&t.Data)

	return &t, nil
}

//岗位发放流水(分页)
type userFlowDataResp struct {
	Code int          `json:"code"`
	Msg  string       `json:"msg"`
	Data userFlowData `json:"Data"`
}

type userFlowData struct {
	Total int32      `json:"total"`
	Pages int32      `json:"pages"`
	List  []userFlow `json:"list"`
}
type userFlow struct {
	Amount         float64 `json:"amount"`
	PayOrderItemId string  `json:"payOrderItemId"`
	PaymentResTime string  `json:"paymentResTime"`
}

func GetUserFlowPage(userId string, companyId int32, workId int32, pageNumber int32, pageItem int32) (userFlowDataResp, error) {
	var t userFlowDataResp
	urls := config.Viper.GetString("POST_EVIDENCE_URL") + "/post-evidence/user-post/flow/page"
	log.Print(urls)
	post_data := map[string]interface{}{
		"userId":        userId,
		"custCompanyId": companyId,
		"postId":        workId,
		"pageNo":        pageNumber,
		"pageSize":      pageItem,
	}
	log.Printf("入参:user_id:%d,companyid:%s,workid:%s", post_data["userId"], post_data["custCompanyId"], post_data["postId"])
	bytesData, err := json.Marshal(post_data)

	resp, err := eureka.Post(urls, bytes.NewReader(bytesData))

	if err != nil {
		return t, err
	}
	err = json.NewDecoder(resp.Body).Decode(&t)
	log.Printf("接受结果,code:%d,msg:%s", t.Code, t.Msg)
	if t.Code != http.StatusOK {
		return t, err
	}
	return t, err
}
