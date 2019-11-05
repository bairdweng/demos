package user

// SessionUser session对象
type SessionUser struct {
	UserID             int32                  `json:"userId"`
	UserName           string                 `json:"userName"`
	Mobile             string                 `json:"mobile"`
	UserType           string                 `json:"userType"`
	SourceName         string                 `json:"sourceName"`
	SessionID          string                 `json:"sessionId"`
	OpenID             string                 `json:"openId"`
	ExpiresIn          int32                  `json:"expiresIn"`
	ProfileID          int32                  `json:"profileId"`
	CompanyID          int32                  `json:"companyId"`
	CompanyName        string                 `json:"companyName"`
	ExpiresAt          int64                  `json:"expiresAt"`
	RoleIds            []int32                `json:"roleIds"`
	Permissions        []string               `json:"permissions"`
	ProfilePermissions map[string]interface{} `json:"profilePermissions"`
	SubjectID          string                 `json:"subjectId"`
	UniqueUserId       string                 `json:"uniqueUserId"`
}

type User struct {
	UserID               int64  `json:"id" binding:"required"`
	UserName             string `json:"username" binding:"required"`
	Name                 string `json:"name" binding:"required"`
	Email                string `json:"email"`
	Mobile               string `json:"mobile"`
	SourceType           string `json:"sourceType"`
	Avatar               string `json:"avatar"`
	WxUnionId            string `json:"wxUnionId"`
	PaymentUserId        int64  `json:"paymentUserId"`
	Active               bool   `json:"actived"`
	VerifyMobileRequired bool   `json:"verifyMobileRequired"`
	VerifiedLevel        int64  `json:"verifiedLevel"`
	IdNumber             string `json:"idNumber"`
	LoginName            string `json:"loginName"`
	Channel              string `json:"channel"`
	RegIp                string `json:"regIp"`
	RegAt                int64  `json:"regAt"`
	LastLoginIp          string `json:"lastLoginIp"`
	LastLoginAt          string `json:"lastLoginAt"`
}

type Verified struct {
	ExtraSystemId  string `json:"extrSystemId"`
	RequestId      int    `json:"requestId"`
	SignType       string `json:"signType"`
	Sign           string `json:"sign"`
	Nonce          string `json:"nonce"`
	Timestamp      int64  `json:"timestamp"`
	NotifyUrl      string `json:"notifyUrl"`
	Name           string `json:"name"`
	IdCard         string `json:"idcard"`
	ValidType      int    `json:"validType"`
	Mobile         string `json:"mobile"`
	PayAccountType string `json:"payAccountType"`
	PayAccount     string `json:"payAccount"`
	BankName       string `json:"bankName"`
}

type VerifiedInput struct {
	Name           string `json:"name"`
	IdCard         string `json:"idcard"`
	ValidType      int    `json:"validType"`
	PayAccountType string `json:"payAccountType"`
	PayAccount     string `json:"payAccount"`
	Mobile         string `json:"mobile"`
}

type VerifiedResp struct {
	UserId int64 `json:"userId"`
}

type UserServiceVerfiedInput struct {
	CompanyId   int    `json:"CompanyId" binding:"required"`
	CompanyName string `json:"CompanyName" binding:"required"`
	RealName    string `json:"realName" binding:"required"`
	IdCardNo    string `json:"idCardNo" binding:"required"`
	MobilePhone string `json:"mobilePhone" binding:"required"`
	BankCardNo  string `json:"bankCardNo"`
}
