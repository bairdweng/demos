package Api

import (
	"bytes"
	"encoding/json"
	"iQuest/app/eureka"
	"iQuest/app/model/user"
	"iQuest/config"
)

type VerifiedResp struct {
	Code int              `json:"code"`
	Data VerifiedListInfo `json:"userId"`
	Msg  string           `json:"msg"`
}

type VerifiedListInfo struct {
	Code   int    `json:"code"`
	UserId int64  `json:"userId"`
	Msg    string `json:"msg"`
}

func UserVerifiedAndCreate(verified user.UserServiceVerfiedInput) (*VerifiedListInfo, error) {

	post_data := map[string]interface{}{
		"companyId":       verified.CompanyId,
		"companyName":     verified.CompanyName,
		"realName":        verified.RealName,
		"credentialsNo":   verified.IdCardNo,
		"mobilePhone":     verified.MobilePhone,
		"bankCardNo":      verified.BankCardNo,
		"companyType":     "company",
		"credentialsType": "idcard",
	}
	bytesData, err := json.Marshal(post_data)

	resp, err := eureka.Post(config.Viper.GetString("USER_SERVICE_VERIFY"), bytes.NewReader(bytesData))

	var t VerifiedListInfo
	if err != nil {
		return nil, err
	}

	_ = json.NewDecoder(resp.Body).Decode(&t)

	return &t, nil
}
