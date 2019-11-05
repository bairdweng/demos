/***
 * 工单系统api
 */
package Api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/satori/go.uuid"
	"iQuest/app/request/job"
	"iQuest/config"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	AuditUrl             = "/api/taxplan-workflow/bizProcess/restStartInstance"
	BusinessType         = "asr-create-job"
	ProcessDefinitionKey = "asr-create-job-flow"
)

type CreateJobTemplateResponse struct {
	Code int  `json:"code"`
	Data bool `json:"data"`
}

func CreateJobTemplate(data job.CreateJobTemplateInput, token string) (*CreateJobTemplateResponse, error) {
	uuid := uuid.Must(uuid.NewV4(), nil)

	dataStr, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	requestData := map[string]interface{}{
		"msgId": uuid.String(),
		"body":  string(dataStr),
	}

	requestByte, err := json.Marshal(requestData)
	if err != nil {
		log.Println("json marshal err: ", err)
		log.Println("json : ", string(requestByte))
		return nil, err
	}
	auditUrl := config.Viper.GetString("OPENADMIN_DOMAIN") + AuditUrl
	log.Println("送审url:", auditUrl)

	request, err := http.NewRequest(http.MethodPost, auditUrl, bytes.NewReader(requestByte))
	if err != nil {
		log.Println("构造请求出错: ", err)
		return nil, err
	}
	/*
		if config.Viper.GetBool("DEBUG") {
			token = config.Viper.GetString("AYG_SESSIONID") //TODO 暂时写死,依赖前端传过来
		}*/
	cookie := &http.Cookie{Name: config.Viper.GetString("HEADER_AUTH_RAW"), Value: token}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set(config.Viper.GetString("HEADER_AUTH_RAW"), token)
	request.AddCookie(cookie)
	client := &http.Client{}
	log.Printf("请求头:%v, 数据:%s\n", request, string(requestByte))

	resp, err := client.Do(request)

	if err != nil {
		log.Println("POST请求出错: ", err)
		return nil, err
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	log.Printf("送审结果: %s \n ", string(body))

	if resp.StatusCode != http.StatusOK {
		log.Printf("POST请求status code : %d, err: %s \n ", resp.StatusCode, err)

		return nil, errors.New(fmt.Sprintf("status code:%d, \t;body:%s", resp.StatusCode, string(body)))
	}

	var r *CreateJobTemplateResponse
	err = json.Unmarshal(body, &r)
	if err != nil {
		log.Printf("POST请求响应错误:  %s \n ", string(body))
		return nil, err
	}

	if http.StatusOK != r.Code {
		log.Printf("业务错误:  %s \n ", string(body))
		return nil, errors.New(string(body))
	}

	return r, nil
}
