package Api

import (
	"bytes"
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/gogf/gf/util/gconv"
	"iQuest/app/eureka"
	"iQuest/app/graphql/model"
	"iQuest/config"
	"iQuest/library/response"
	"io/ioutil"
	"log"
	mr "math/rand"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type SignSubmitType struct {
	ExtrOrderId      string `json:"extrOrderId"`      //外部订单ID
	ExtrSystemId     string `json:"extrSystemId"`     //appid
	Identity         string `json:"identity"`         //证件号
	Name             string `json:"name"`             //姓名
	PersonalMobile   string `json:"personalMobile"`   //手机号
	TemplateId       string `json:"templateId"`       //模板ID
	IdentityType     string `json:"identityType"`     //证件类型
	Sign             string `json:"sign"`             //签名
	ServiceCompanyId string `json:"serviceCompanyId"` //签约服务商id
	UserId           string `json:"userId"`           //用户唯一标识
	CompanyId        string `json:"companyId"`        //企业id
}

type SignRspDataType struct {
	ResultCode    string `json:"resultCode"`
	State         string `json:"state"`
	StateDesc     string `json:"stateDesc"`
	ResultMessage string `json:"resultMessage"`
	PartybSignUrl string `json:"partybSignUrl"`
	PartycSignUrl string `json:"partycSignUrl"` //先用c， 为空， 判断b,b不为空则用B，bc为空是自动签约
}

type SignQueryType struct {
	ExtrOrderId  string `json:"extrOrderId"`  //外部订单ID
	ExtrSystemId string `json:"extrSystemId"` //appid
	Sign         string `json:"sign"`         //签名
}

type SignResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type CompayData struct {
	Id           int    `json:"id"`
	AppId        string `json:"appId"`
	CompanyName  string `json:"companyName"`
	IsFromOutApp int    `json:"isFromOutApp"`
}

func SignSubmit(submitType SignSubmitType) (*model.SignRspData, error) {

	//resultCode := [4]string{"ACCEPTED","AUTHING", "SIGNING", "CLOSED"}

	body, err := json.Marshal(submitType)
	if err != nil {
		return nil, err
	}
	log.Printf("SignSubmit:body:%v", string(body))

	var t model.SignRspData

	reader := bytes.NewReader(body)
	url := config.Viper.GetString("SIGN_API_URL") + "/extr/order/inner-submit"
	resp, err := eureka.Post(url, reader)
	if err != nil {
		return nil, err
	}
	log.Printf("SignSubmitres:Status:%v", resp.Status)
	if strings.Replace(resp.Status, " ", "", -1) != "200" {
		return nil, errors.New("请求签约接口失败")
	}

	err = json.NewDecoder(resp.Body).Decode(&t)
	if err != nil {
		return nil, err

	}

	//var codeResp SignResp
	//if t.State == nil {
	//	err = json.NewDecoder(resp.Body).Decode(&codeResp)
	//	if err != nil{
	//		return nil,err
	//
	//	}
	//	if codeResp.Code == 300000 {
	//
	//	}
	//}

	return &t, nil

}

func GetExtrOrderId(submitType SignSubmitType) string {
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	str := submitType.UserId + submitType.CompanyId + submitType.ServiceCompanyId + timeStr
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func SignQuery(queryType SignQueryType) (*model.SignRspData, error) {
	body, err := json.Marshal(queryType)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	reader := bytes.NewReader(body)
	url := config.Viper.GetString("SIGN_API_URL") + "/extr/order/inner-qry"
	log.Print(url)

	resp, err := eureka.Post(url, reader)

	//req, err := http.NewRequest("POST", url, reader)
	var t model.SignRspData
	if err != nil {

		return nil, err

	}
	//req.Header.Add("Content-Type", "application/json;charset=UTF-8")
	//client := http.Client{}
	//resp, err := client.Do(req)
	//if err != nil {
	//
	//	return nil, err
	//}
	if strings.Replace(resp.Status, " ", "", -1) != "200" {
		log.Print("签约查询失败")
		return nil, errors.New("签约查询失败")
	}

	err = json.NewDecoder(resp.Body).Decode(&t)
	log.Printf("res state:%s,StateDesc:%s", t.State, t.ResultCode)
	if err != nil {
		return nil, err
	}
	if t.ResultCode == nil {
		state := "FAIL"
		stateDec := "没有签约信息"
		return &model.SignRspData{
			State:         &state,
			StateDesc:     &stateDec,
			ResultCode:    &state,
			ResultMessage: &stateDec,
		}, nil
	}
	return &t, nil
}

func GetAppId(companyId string) (*string, error) {

	t := make([]CompayData, 3)

	//key := "ishouru_getAppId_by_companyId_" + companyId
	//value, _ := db.Redis().Get(key).Result()
	//if value != "" {
	//	return &value,nil
	//}

	url := config.Viper.GetString("COMPANY_APP_URL") + "/company-app/company-all-apps?companyId=" + fmt.Sprintf("%s", companyId)
	log.Print(url)
	resp, _ := eureka.Get(url)
	err := json.NewDecoder(resp.Body).Decode(&t)
	body, _ := ioutil.ReadAll(resp.Body)
	log.Printf("getappid:  %s \n ", string(body))
	log.Printf("getappid:%v", t[0].AppId)
	if err != nil {
		return nil, err
	}

	if len(t) == 0 {
		return nil, errors.New("查不到该企业的APPid")
	}
	app_id := t[0].AppId
	//db.Redis().Set(key,app_id,time.Minute)

	return &app_id, nil
}

type SignInfoResp struct {
	Code int
	Data []SignList
}

type SignList struct {
	Name       string         `json:"personalName"`
	ServerName string         `json:"serverName"`
	SignTime   int64          `json:"signTime"`
	CompanyId  string         `json:"companyId"`
	IsGroup    bool           `isGroup:"companyId"`
	GroupInfo  []GroupInfoArr `json:"personalOrderGroupInfo"`
}

type GroupInfoArr struct {
	FileName   string `json:"fileName"`
	PreviewUrl string `json:"previewUrl"`
	ServerName string `json:"serverName"`
}

func GetSignInfo(cardNo string) (*SignInfoResp, error) {

	var t SignInfoResp
	url := config.Viper.GetString("SIGN_API_URL") + "/extr/personal/inner-signlist"
	log.Println(url)
	post_data := map[string]interface{}{
		"identity": cardNo,
		"name":     "test",
		"qryFlag":  "1",
		"sign":     "123",
	}
	bytesData, err := json.Marshal(post_data)

	//resp, err := http.Post(url,
	//	"application/json;charset=UTF-8",
	//	bytes.NewReader(bytesData),
	//)
	//
	//if err != nil {
	//
	//	return nil, err
	//}
	//
	//err = json.NewDecoder(resp.Body).Decode(&t.Data)
	//return &t, nil

	resp, err := eureka.Post(url, bytes.NewReader(bytesData))

	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(resp.Body).Decode(&t.Data)
	log.Printf("%v,%v", t.Code, t.Data)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

type IsSignResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data bool   `json:"data"`
}

//传入 user_id company_id  service_company_id 返回是否已签约
func IsSign(userId int64, companyId int32, serviceCompanyId string) (*IsSignResp, error) {

	var t IsSignResp
	url := config.Viper.GetString("SIGN_API_URL") + "/inner/order/qrySign-by-userid"
	log.Printf(url)
	log.Printf("接口参数:userId:%d,companyId:%d,serviceCompanyId:%v", userId, companyId, serviceCompanyId)
	post_data := map[string]interface{}{
		"userId":           strconv.Itoa(int(userId)),
		"companyId":        strconv.Itoa(int(companyId)),
		"serviceCompanyId": serviceCompanyId,
	}
	bytesData, err := json.Marshal(post_data)

	resp, err := eureka.Post(url, bytes.NewReader(bytesData))

	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(resp.Body).Decode(&t)
	log.Printf("返回结果:code:%d,msg:%s,data:%v", t.Code, t.Msg, t.Data)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

type SignQueryByCompanyType struct {
	UserId           string `json:"userId"`
	CompanyId        string `json:"companyId"`
	ServiceCompanyId string `json:"serviceCompanyId"`
	ExtrSystemId     string `json:"extrSystemId"` //appid
	Sign             string `json:"sign"`
}

func SignQueryByCompany(companyType SignQueryByCompanyType) (*model.SignRspData, error) {
	log.Printf("qry-for-ishouru-params userId:%s,companyId:%s,extrSystemId:%s,sign:%s", companyType.UserId, companyType.CompanyId, companyType.ExtrSystemId, companyType.Sign)
	body, err := json.Marshal(companyType)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	reader := bytes.NewReader(body)
	url := config.Viper.GetString("SIGN_API_URL") + "/extr/order/inner-qry-for-ishouru"
	log.Print(url)

	resp, err := eureka.Post(url, reader)

	var t model.SignRspData
	if err != nil {

		return nil, err

	}

	if strings.Replace(resp.Status, " ", "", -1) != "200" {
		log.Print("签约查询失败")
		return nil, errors.New("签约查询失败")
	}

	err = json.NewDecoder(resp.Body).Decode(&t)
	log.Printf("qry-for-ishouru-res state:%s,StateDesc:%s", t.State, t.ResultCode)
	if err != nil {
		return nil, err
	}
	if t.ResultCode == nil {
		state := "FAIL"
		stateDec := "没有签约信息"
		return &model.SignRspData{
			State:         &state,
			StateDesc:     &stateDec,
			ResultCode:    &state,
			ResultMessage: &stateDec,
		}, nil
	}
	return &t, nil
}

type IdentityResp struct {
	Code interface{}      `json:"code"`
	Msg  string           `json:"msg"`
	Data IdentityRespData `json:"data"`
}

type IdentityRespData struct {
	CertResult    string `json:"certResult"`
	CertResultMsg string `json:"certResultMsg"`
}

func Identity(backFile string, frontFile string, identity string, identityType string, name string) (*IdentityResp, error) {

	var publicHeader = "\n-----BEGIN RSA PRIVATE KEY-----\n"
	var publicTail = "-----END RSA PRIVATE KEY-----\n"
	var temp string
	split(config.Viper.GetString("SIGN_RSA_KEY"), &temp)

	pk8 := []byte(publicHeader + temp + publicTail)

	fronte_id_card, _ := base64.StdEncoding.DecodeString(frontFile) //成图片文件并把文件写入到buffer
	id_card_front := bytes.NewBuffer(fronte_id_card)
	front_file_md5 := md5V(id_card_front.Bytes())
	front_file_md5 = fmt.Sprintf("%x", front_file_md5)

	back_id_card, _ := base64.StdEncoding.DecodeString(backFile) //成图片文件并把文件写入到buffer
	id_card_back := bytes.NewBuffer(back_id_card)
	back_file_md5 := md5V(id_card_back.Bytes())
	back_file_md5 = fmt.Sprintf("%x", back_file_md5)
	nonce := RandStringRunes(32)

	sign_str := "appId=" + config.Viper.GetString("PRIVATE_APP_ID") + "&backfile=" + back_file_md5 + "&frontfile=" + front_file_md5 + "&identity=" + identity + "&identityType=" + identityType + "&name=" + name + "&nonce=" + nonce
	sign, err := RsaSignWithSha1Hex(sign_str, pk8)

	if err != nil {
		return nil, err
	}

	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	fw, _ := w.CreateFormFile("frontfile", "test1.jpg")
	_, _ = fw.Write(id_card_front.Bytes())

	fw2, _ := w.CreateFormFile("backfile", "test2.jpg")
	_, _ = fw2.Write(id_card_back.Bytes())

	_ = w.Close()

	post_url := fmt.Sprintf(config.Viper.GetString("SIGN_IDENTITY")+"/econtract/extr/identity/upload?appId=%s&name=%s&identity=%s&identityType=%s&sign=%v&nonce="+nonce, config.Viper.GetString("PRIVATE_APP_ID"), url.QueryEscape(name), identity, identityType, sign)

	req, err := http.NewRequest("POST", post_url, &b)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	var t IdentityResp

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(res.Body).Decode(&t)
	log.Printf("上传身份证接口:code:%v,msg:%v", t.Code, t.Msg)
	if err != nil {
		return nil, err
	}

	if gconv.String(t.Code) != "0000" {
		t = IdentityResp{
			Code: response.Error,
			Msg:  t.Msg,
		}
		return &t, nil
	}

	t = IdentityResp{
		Code: response.Success,
		Msg:  t.Msg,
	}
	return &t, nil

}

func md5V(bytes_file []byte) string {
	h := md5.New()
	h.Write(bytes_file)
	return string(h.Sum([]byte("")))
}

//func Sha1WithRsa(data, privateKeyBytes []byte) (string, error) {
//	block, _ := pem.Decode(privateKeyBytes)
//	if block == nil {
//		return "", errors.New("私钥为空")
//	}
//
//	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
//	if err != nil {
//		return "", errors.New("解析私钥错误:" + err.Error())
//	}
//
//	encryptStr, err := RsaSign(data, privateKey)
//	if err != nil {
//		return "", err
//	}
//	encryptStr = url.QueryEscape(encryptStr)
//	return encryptStr, nil
//}

func RsaSign(data string, privateKey *rsa.PrivateKey) (string, error) {

	h := crypto.SHA256.New()
	h.Write([]byte(data))
	hashed := h.Sum(nil)

	//signature, err := rsa.SignPKCS1v15(nil, privateKey, crypto.SHA1, hash)
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed)
	if err != nil {
		return "", errors.New("rsa签名 错误: " + err.Error())
	}
	//signStr := hex.EncodeToString(signature)
	//base64.RawURLEncoding.EncodeToString(sign)
	signStr := base64.StdEncoding.EncodeToString(signature)
	return signStr, nil
}

func RsaSignWithSha1Hex(data string, privateKeyBytes []byte) (string, error) {
	block, _ := pem.Decode(privateKeyBytes)
	if block == nil {
		return "", errors.New("私钥为空")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", errors.New("解析私钥错误:" + err.Error())
	}

	encryptStr, err := RsaSign(data, privateKey.(*rsa.PrivateKey))

	if err != nil {
		return "", err
	}
	encryptStr = url.QueryEscape(encryptStr)
	return encryptStr, nil
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[mr.Intn(len(letterRunes))]
	}
	return string(b)
}

func split(key string, temp *string) {
	if len(key) <= 64 {
		*temp = *temp + key + "\n"
	}
	for i := 0; i < len(key); i++ {
		if (i+1)%64 == 0 {
			*temp = *temp + key[:i+1] + "\n"
			fmt.Println(len(*temp) - 1)
			key = key[i+1:]
			split(key, temp)
			break
		}
	}
}
func Certification(identity string, name string) (*IdentityResp, error) {
	timeUnix := time.Now().UnixNano() / 1e6
	type req_data struct {
		ExtrSystemId   string `json:"extrSystemId"`
		Idcard         string `json:"idcard"`
		Name           string `json:"name"`
		Nonce          string `json:"nonce"`
		PayAccountType string `json:"payAccountType"`
		RequestId      string `json:"requestId"`
		Sign           string `json:"sign"`
		SourceType     string `json:"sourceType"`
		Timestamp      string `json:"timestamp"`
		ValidType      int    `json:"validType"`
	}
	body := req_data{
		ExtrSystemId:   "hgt",
		Idcard:         identity,
		Name:           name,
		Nonce:          RandStringRunes(32),
		PayAccountType: "ALIPAY_USERID",
		RequestId:      strconv.Itoa(int(timeUnix)),
		Sign:           "hgt",
		SourceType:     "hgt",
		Timestamp:      strconv.Itoa(int(timeUnix)),
		ValidType:      1,
	}

	bodyJson, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(bodyJson)
	url := "http://prepare-service/sync/inner-certification"
	resp, err := eureka.Post(url, reader)
	if err != nil {
		return nil, err
	}
	var t IdentityResp

	err = json.NewDecoder(resp.Body).Decode(&t)
	log.Printf("实名认证接口:code:%v,msg:%v,code:%v", t.Code, t.Data.CertResultMsg, t.Data.CertResult)
	if err != nil {
		return nil, err
	}

	if gconv.String(t.Code) != "0000" {
		t = IdentityResp{
			Code: response.Error,
			Msg:  t.Msg,
		}
		return &t, nil
	}

	if t.Data.CertResult != "1" {
		t = IdentityResp{
			Code: response.Error,
			Msg:  t.Data.CertResultMsg,
		}
		return &t, nil
	}

	t = IdentityResp{
		Code: response.Success,
		Msg:  t.Data.CertResultMsg,
	}
	return &t, nil

}
