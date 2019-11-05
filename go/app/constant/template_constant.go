package constant

const (
	//# 是否通过审核 0待审核, 1 通过, 2 拒绝 , 3 (已审核,合同生效期在未来)未生效, 4.已失效(合同过期)
	JobTemplateUnaided     = iota //未审核
	JobTemplateEnable             // 通过
	JobTemplateReject             //拒绝
	JobTemplateExpired            //已失效
	JobTemplateUnactivated        //未生效

	JobTemplateUnAudit = 0 //未审核
	JobTemplateAudited = 1 //已审核

	JobTemplateSourceCompany = 0 //企业自主创建(爱收入)

	JobTemplateSourceAudit = 1 //风控审核后创建(税筹)

	JobJournalUnhandled = 0 //未处理
	JobJournalHandled   = 1 //已处理
)
