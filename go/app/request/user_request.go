package request

type GetUserRequest struct {
	UserId int64 `form:"user_id" json:"user_id" binding:"required"`
}

type GetUserTokenRequest struct {
	Token string `form:"token" json:"token" binding:"required"`
}

type VerifiedInput struct {
	Name           string `form:"name" json:"name" binding:"required"`
	IdCard         string `form:"idcard" json:"idcard" binding:"required"`
	ValidType      int    `form:"validType" json:"validType" binding:"required"`
	PayAccountType string `form:"payAccountType" json:"payAccountType" binding:"required"`
	PayAccount     string `form:"payAccount" json:"payAccount"`
	Mobile         string `form:"mobile" json:"mobile"`
}

type GetImportLogsRequest struct {
	FileHash string `form:"id" json:"id" binding:"required"`
	PageSize int32  `form:"page_size" json:"page_size" `
	PageNum  int32  `form:"page_num" json:"page_num" `
}

type IdentityRequest struct {
	FrontFile    string `form:"frontfile" json:"frontfile" binding:"required"`
	BackFile     string `form:"backfile" json:"backfile" binding:"required"`
	Name         string `form:"name" json:"name" binding:"required"`
	Identity     string `form:"identity" json:"identity" binding:"required"`
	IdentityType string `form:"identityType" json:"identityType" binding:"required"`
}
