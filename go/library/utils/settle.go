package utils

import (
	"context"
	"crypto/md5"
	"fmt"
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
	"iQuest/app/constant"
	"math/rand"
	"path"
	"strconv"
	"time"
)

func Unzip(companyId int32) bool{
	for i := 0; i < len(constant.XINNIAO_COMPANY_ID_ARR); i++ {
		if constant.XINNIAO_COMPANY_ID_ARR[i] == companyId {
			return true
		}
	}
	return false
}

func UploadToQiNiu(fileUrl string) (string,error) {
	host := "https://file.ishouru.com/"
	bucket := "iquest-file"
	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}
	putPolicy.Expires = 7200 //示例2小时有效期
	mac := qbox.NewMac("LS8BEp91TfmH5GV5hk5ct-ojvs5fZMd1EIPl89iZ", "K24kSASmfjpBmCUnVk1HQFdv9sI-WCRQark3fZs3")
	upToken := putPolicy.UploadToken(mac)

	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Zone = &storage.ZoneHuadong
	// 是否使用https域名
	cfg.UseHTTPS = true
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	putExtra := storage.PutExtra{
		Params: map[string]string{
			"x:name": "github logo",
		},
	}

	unix := time.Now().UnixNano()
	b := make([]byte, 10)
	n, _ := rand.Read(b)

	timeByte := []byte(strconv.Itoa(int(unix)) + string(b) + strconv.Itoa(n))
	filename := "settlement/" + fmt.Sprintf("%x", md5.Sum(timeByte))
	var filenameWithSuffix string
	filenameWithSuffix = path.Base(fileUrl) //获取文件名带后缀
	fmt.Println("filenameWithSuffix =", filenameWithSuffix)
	var fileSuffix string
	fileSuffix = path.Ext(filenameWithSuffix) //获取文件后缀

	err := formUploader.PutFile(context.Background(), &ret, upToken, filename + fileSuffix, fileUrl, &putExtra)
	if err != nil {
		return "",err
	}
	return host + filename + fileSuffix,err
}
