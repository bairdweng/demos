package Api

import (
	"bytes"
	"encoding/json"
	"iQuest/app/eureka"
	"iQuest/config"
	"iQuest/library/response"
	"log"
	"net/http"

	"github.com/speps/go-hashids"
)

type AYGResp struct {
	Code int `json:"code"`
}

func SendSms(mobile string, content string, appId string) (AYGResp, error) {
	var t AYGResp
	urls := config.Viper.GetString("COMMON_APP_URL") + "/api/common/sms/send"
	log.Print(urls)
	post_data := map[string]interface{}{
		"mobile":   mobile,
		"contents": content,
		"smsNo":    "",
		"type":     "sms",
		"appId":    appId,
	}
	bytesData, err := json.Marshal(post_data)

	resp, err := eureka.Post(urls, bytes.NewReader(bytesData))

	if err != nil {
		t = AYGResp{
			Code: response.Error,
		}
		return t, err
	}

	err = json.NewDecoder(resp.Body).Decode(&t)
	log.Printf("res code:%d", t.Code)
	if t.Code != http.StatusOK {
		t = AYGResp{
			Code: response.Error,
		}
		return t, nil
	}
	return t, nil

}

//生成短网址
type DwzResp struct {
	Code     int    `json:"Code"`
	ShortUrl string `json:"ShortUrl"`
}

func Dwz(taskMemberID int64) (string, error) {
	var t DwzResp
	token := "ea36e73ae6904e7b222bd06fd100320c"
	host := "https://dwz.cn"
	path := "/admin/v2/create"
	urls := host + path
	hd := hashids.NewData()
	hd.MinLength = 30
	h, _ := hashids.NewWithData(hd)
	e, _ := h.Encode([]int{int(taskMemberID)})

	longUrl := "https://open.aiyuangong.com/asr_admin/admin/index.html#/proof?task_member_id=" + e
	log.Print(longUrl)
	post_data := map[string]interface{}{
		"url": longUrl,
	}
	bytesData, err := json.Marshal(post_data)
	request, err := http.NewRequest(http.MethodPost, urls, bytes.NewReader(bytesData))
	if err != nil {
		log.Println("构造请求出错: ", err)
		return "", err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Token", token)
	client := &http.Client{}
	log.Printf("请求头:%v, 数据:%s\n", request, string(bytesData))

	resp, err := client.Do(request)

	if err != nil {
		log.Println("POST请求出错: ", err)
		return "", err
	}
	err = json.NewDecoder(resp.Body).Decode(&t)
	log.Printf("res code:%d,shortUrl:%s,err:%v", t.Code, t.ShortUrl, err)
	if t.Code == 0 {
		return t.ShortUrl, nil
	}
	return longUrl, nil
}
