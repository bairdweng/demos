package request

type IsJoinJobRequest struct {
	ID  int32 `form:"id" json:"id" binding:"required"`
	Uid int64 `form:"uid" json:"uid" binding:"required"`
}

type JobInfo struct {
	ID int32 `form:"id" json:"id" binding:"required"`
}

type JobsRequest struct {
	CompanyId        int32 `form:"companyId" json:"companyId" binding:"required"`
	ServiceCompanyId int32 `form:"serviceCompanyId" json:"serviceCompanyId" binding:"required"`
}

type IsBindPhoneRequest struct {
	UserId int64 `form:"user_id" json:"user_id" binding:"required"`
}

type BindPhoneRequest struct {
	UserId int64 `form:"user_id" json:"user_id" binding:"required"`
	Mobile int64 `form:"mobile" json:"mobile" binding:"required"`
}

// DownloadPayrollRequest 下载工资单请求
type DownloadPayrollRequest struct {
	WorkID int   `form:"work_id" json:"work_id" binding:"required"`
	IDs    []int `form:"ids" json:"ids" binding:"required"`
}

// 上传绩效数据
type CreateSettlementRecords struct {
	Name string                 `form:"name" json:"name" binding:"required"`
	Path map[string]interface{} `form:"path" json:"path" binding:"required"`
}

// DownloadPayrollRequest 下载工资单请求
type DownloadSettlementRequest struct {
	WorkID   int    `form:"work_id" json:"work_id" binding:"required"`
	FileName string `form:"file_name" json:"file_name" binding:"required"`
}

// 上传绩效数据
type CreateSettlementAndDownloadFileRequest struct {
	Key 	string   `form:"key" json:"key" binding:"required"`
	WorkID 	int   `form:"work_id" json:"work_id" binding:"required"`
}

//下载单人结算模板
type DownLoadSingleSettleTemplateRequest struct {
	MemberId int `form:"member_id" json:"member_id" binding:"required"`
}

//下载单人结算模板
type DownLoadBatchSettleTemplateRequest struct {
	WorkId int `form:"work_id" json:"work_id" binding:"required"`
}
