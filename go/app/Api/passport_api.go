package Api

import (
	"encoding/json"
	"fmt"
	"iQuest/app/constant"
	"iQuest/app/eureka"
	"iQuest/app/model/user"
	"iQuest/config"
	"iQuest/library/response"
	"log"
	"net/http"
)

type Resp struct {
	Code int       `json:"code"`
	Data user.User `json:"data"`
	Msg  string    `json:"message"`
}

func GetUserById(id int64) (Resp, error) {

	var t Resp
	url := config.Viper.GetString("PASS_PORT_URL") + fmt.Sprintf("%d", id)
	log.Print(url)

	resp, err := eureka.Get(url)
	if err != nil {
		t = Resp{
			Code: response.Error,
			Data: user.User{},
			Msg:  constant.INTERFACE_ERROR,
		}
		return t, err
	}

	err = json.NewDecoder(resp.Body).Decode(&t)
	log.Printf("res code:%d,msg:%s,data:%v", t.Code, t.Msg, t.Data)
	if t.Code != http.StatusOK {
		t = Resp{
			Code: response.Error,
			Data: user.User{},
			Msg:  constant.USER_NOT_EXIST,
		}
		return t, nil
	}

	return t, nil

}

func GetUserByToken(token string) (Resp, error) {

	var t Resp
	url := config.Viper.GetString("PASS_PORT_URL") + "info/" + token

	resp, err := eureka.Get(url)
	if err != nil {
		t = Resp{
			Code: response.Error,
			Data: user.User{},
			Msg:  constant.INTERFACE_ERROR,
		}
		return t, err
	}

	err = json.NewDecoder(resp.Body).Decode(&t)
	if t.Code != http.StatusOK {
		t = Resp{
			Code: response.Error,
			Data: user.User{},
			Msg:  constant.USER_NOT_EXIST,
		}
		return t, nil
	}

	return t, nil

}

type RespUserInfo struct {
	Code int64        `json:"code"`
	Data UserInfoData `json:"data"`
	Msg  string       `json:"message"`
}
type UserInfoData struct {
	Code                int64                 `json:"code"`
	Msg                 string                `json:"msg"`
	CredentialsNo       string                `json:"credentialsNo"`   //证件号
	CredentialsType     string                `json:"credentialsType"` //证件类型  credentialsType = idcard 身份证
	RealName            string                `json:"realName"`
	RelationCompanyList []RelationCompanyList `json:"relationCompanyList"`
	State               int                   `json:"verifyState"`
	UserId              int64                 `json:"userId"`
	AccountList         []AccountList         `json:"accountList"`
	MobilePhone         string                `json:"mobilePhone"`
}

type RelationCompanyList struct {
	CompanyId   int    `json:"companyId"`
	CompanyName string `json:"companyName"`
	CompanyType string `json:"companyType"`
}

type AccountList struct {
	AccountNo   string `json:"accountNo"`
	AccountName string `json:"accountName"`
	AccountType string `json:"accountType"`
	State       int    `json:"verifyState"`
}

func FindRealnameInfoByLoginid(Loginid int64) (*RespUserInfo, error) {

	var t RespUserInfo
	var td UserInfoData

	url := config.Viper.GetString("PASS_PORT_URL_USER") + "/user/find-realname-user-by-loginid?loginId=" + fmt.Sprintf("%d", Loginid)

	resp, err := eureka.Get(url)
	if err != nil {

		return nil, err
	}

	err = json.NewDecoder(resp.Body).Decode(&td)
	if err != nil {
		return nil, err
	}

	t = RespUserInfo{
		Code: td.Code,
		Data: td,
		Msg:  "用户服务：" + td.Msg,
	}

	return &t, nil

}

func FindRealnameInfoByUserid(userId int64) (*RespUserInfo, error) {

	var t RespUserInfo
	var td UserInfoData

	url := config.Viper.GetString("PASS_PORT_URL_USER") + "/user/find-realname-user-by-userid?userId=" + fmt.Sprintf("%d", userId)
	log.Print(url)
	resp, err := eureka.Get(url)
	if err != nil {

		return nil, err
	}

	err = json.NewDecoder(resp.Body).Decode(&td)
	log.Printf("%v",td)
	if err != nil {

		return nil, nil
	}

	t = RespUserInfo{
		Code: td.Code,
		Data: td,
		Msg:  "用户服务：" + td.Msg,
	}

	return &t, nil

}

//type VerifiedResp struct {
//	Code string            `json:"code"`
//	Data user.VerifiedResp `json:"data"`
//}
//
//func UserVerified(verified user.VerifiedInput) (VerifiedResp, error) {
//
//	nonce := getRandomString()
//
//	post_data := map[string]interface{}{
//		"extrSystemId":   constant.EXTRA_SYSTEM_ID,
//		"requestId":      constant.REQUEST_ID,
//		"signType":       constant.SIGN_TYPE,
//		"sign":           "sign",
//		"nonce":          nonce,
//		"timestamp":      time.Now().Unix(),
//		"notifyUrl":      "",
//		"name":           verified.Name,
//		"idcard":         verified.IdCard,
//		"validType":      verified.ValidType,
//		"mobile":         verified.Mobile,
//		"payAccountType": verified.PayAccountType,
//		"payAccount":     verified.PayAccount,
//		"bankName":       "",
//	}
//	bytesData, err := json.Marshal(post_data)
//
//	req, err := http.NewRequest("POST", config.Viper.GetString("VERIFY_API"), bytes.NewReader(bytesData))
//	var t VerifiedResp
//	if err != nil {
//		t = VerifiedResp{
//			Code: response.ErrorString,
//		}
//		return t, err
//
//	}
//	req.Header.Add("Content-Type", "application/json;charset=UTF-8")
//	client := http.Client{}
//	resp, err := client.Do(req)
//	if err != nil {
//		t = VerifiedResp{
//			Code: response.ErrorString,
//		}
//		return t, err
//	}
//
//	err = json.NewDecoder(resp.Body).Decode(&t)
//	if t.Code != "0000" {
//		t = VerifiedResp{
//			Code: response.ErrorString,
//		}
//		return t, err
//	}
//
//	return t, nil
//}
//
//func getRandomString() string {
//	str := "0123456789abcdefghijklmnopqrstuvwxyz"
//	info := []byte(str)
//	result := []byte{}
//	r := rand.New(rand.NewSource(time.Now().UnixNano()))
//	for i := 0; i < constant.LIMIT_NUM; i++ {
//		result = append(result, info[r.Intn(len(info))])
//	}
//	return string(result)
//}
