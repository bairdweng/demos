package job

type CreateJobTemplateRequest struct {
	Appid              string `form:"appid" json:"appid" binding:"required"`
	Name               string `form:"name" json:"name" binding:"required"`
	Requirement        string `form:"requirement" json:"requirement" binding:"required"`
	SettlementRule     string `form:"settlementRule" json:"settlementRule" binding:"required"`
	ServiceTypeID      int32  `form:"serviceTypeId" json:"serviceTypeId" binding:"required"`
	ServiceTypeName    string `form:"serviceTypeName" json:"serviceTypeName" binding:"required"`
	CompanyId          int32  `form:"customerCompanyId" json:"customerCompanyId" binding:"required"`
	CompanyName        string `form:"customerCompanyName" json:"customerCompanyName" binding:"required"`
	ServiceCompanyId   int32  `form:"serviceCompanyId" json:"serviceCompanyId" binding:"required"`
	ServiceCompanyName string `form:"serviceCompanyName" json:"serviceCompanyName" binding:"required"`
	UserId             int32  `form:"userId" json:"userId" binding:"required"`
	SignTemplateId     int32  `form:"signTemplateId" json:"signTemplateId" binding:"required"`
	Remark             string `form:"remark" json:"remark" binding:"required"`
}

/////************** 以下是restful******************
type AuditJobTemplateCallbackRequest struct {
	Id int32 `form:"businessId" json:"businessId,string" binding:"required"`
	//SignTemplateId     int32  `form:"signTemplateId" json:"signTemplateId" `
	Remark  string `form:"remark" json:"remark"`
	Operate string `form:"operate" json:"operate" binding:"required"`
}

//批量创建模板
type BatchCreateTemplateRequest struct {
	CompanyId          int32         `form:"customCompanyId" json:"customCompanyId" binding:"required"`
	CompanyName        string        `form:"customCompanyName" json:"customCompanyName" `
	ServiceCompanyId   int32         `form:"serviceCompanyId" json:"serviceCompanyId" binding:"required"`
	ServiceCompanyName string        `form:"serviceCompanyName" json:"serviceCompanyName" `
	ServiceTypes       []ServiceBean `form:"servicePosList" json:"servicePosList" binding:"required"`
	ContractNo string `form:"contractNo" json:"contractNo" binding:"required"`
	ContractStartDate string `form:"contractStartDate" json:"contractStartDate" binding:"required"`
	ContractEndDate string `form:"contractEndDate" json:"contractEndDate" binding:"required"`
	ContractActiveDate int64 `form:"versionStartDate" json:"versionStartDate"`

}

type ServiceBean struct {
	ServiceTypeId   int32          `form:"serviceId" json:"serviceId" binding:"required"`
	ServiceTypeName string         `form:"serviceName" json:"serviceName" binding:"required"`
	Templates       []TemplateBean `form:"positions" json:"positions" binding:"required"`
}

type TemplateBean struct {
	Id int32 `form:"id" json:"id"`
	Name           string         `form:"posName" json:"posName" binding:"required"`
	Requirement    string         `form:"description" json:"description" binding:"required"`
	SettlementRule string         `form:"performance" json:"performance" binding:"required"`
	Attachment     AttachmentBean `form:"attachment" json:"attachment" binding:"required"`
}

type AttachmentBean struct {
	DownloadCode string `form:"downloadCode" json:"downloadCode" binding:"required"`
	DisplayName string `form:"displayname" json:"displayname" binding:"required"`
}

type BizContent struct {
	Appid       string `json:"appid"`
	Name        string `json:"name"`
	Requirement string `json:"requirement"`
	// 结算规则
	SettlementRule  string `json:"settlementRule"`
	ServiceTypeID   int32  `json:"serviceTypeId"`
	ServiceTypeName string `json:"serviceTypeName"`
}

type CreateJobTemplateInput struct {
	Attach               string   `json:"attach"`
	BizExtendData        string   `json:"bizExtendData"`
	BusinessID           string   `json:"businessId"`
	BusinessType         string   `json:"businessType"`
	CallBackUrl          string   `json:"callBackUrl"`
	CustomerCompanyId    int32    `json:"customerCompanyId"`
	CustomerCompanyName  string   `json:"customerCompanyName"`
	ProcessDefinitionKey string   `json:"processDefinitionKey"`
	ProcessParams        struct{} `json:"processParams"`
	ProfileId            int32    `json:"profileId"`
	ServiceCompanyId     int32    `json:"serviceCompanyId"`
	ServiceCompanyName   string   `json:"serviceCompanyName"`
	UserId               int32    `json:"userId"`
	UserName             string   `json:"userName"`
	Appid                string   `json:"appid"`
}

//审核结果请求
type AuditJobTemplateRequest struct {
	Id           int32  `json:"id" binding:"required"` //就是job_template id
	ProcessInsId string `form:"processInsId" json:"processInsId" binding:"required"`
	TaskId       string `form:"taskId" json:"taskId" binding:"required"`
	UserId       int32  `json:"userId" binding:"required"`
	UserName     string `json:"userName" binding:"required"`
	IsPass       bool   `json:"isPass" binding:"exists"`
	ProfileId    int32  `json:"profileId" binding:"required"`
	Remark       string `form:"remark" json:"remark" `
}

type OpenAdminRequest struct {
	MsgId string `json:"msgId"`
	Body  string `json:"body"`
}

type TemplateListRequest struct {
	CompanyId int32 `json:"company_id" binding:"required"`
	ServiceCompanyId  int32 `json:"service_company_id" binding:"required"`
}