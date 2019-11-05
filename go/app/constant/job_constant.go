package constant

const (
	TypeUnknown = iota
	TypeVideo
	TypeImage
	TypeText
)

const (
	WorkTypeTask = 1
	WorkTypeJob  = 2
)

const (
	PayTypePayNone = iota
	PayTypePayBefore
	PayTypePayAfter
)
//"""任务状态: 1进行中, 2过期, 3完成"""
const (
	ProgressUndefined = iota
	ProgressOnGoing
	ProgressExpire
	ProgressDone

)

const (
	//成员状态
	//任务状态:0:可参加,1:邀请中，2已申请，3已同意参加，4申请完成，5已同意完成，
	//6参与者已评分，7发布者已评分，8互评,9拒绝申请,10拒绝完成,20踢出任务,30:任务已截止招募,40:任务已完成,50:任务被下架,60:任务失效
	STATUS_CANJOIN = iota
	STATUS_INVITE
	STATUS_APPLY
	STATUS_APPROVE
	STATUS_APPLY_COMPLETE
	STATUS_APPROVE_COMPLETE
	STATUS_REFUSE_APPLY = 9
	STATUS_TASK_SHELVES = 50
	STATUS_TASK_EXPIRED = 60
)

const(
	PROGRESS_STATUS_CREATE = "create"; //发布任务,m
	PROGRESS_STATUS_APPLY = "apply"; //申请参加任务,g
	PROGRESS_STATUS_REFUSE = "refuse"; //拒绝参加任务,m
	PROGRESS_STATUS_APPROVE = "approve"; //同意参加任务,m
	PROGRESS_STATUS_INVITE = "invite"; //邀请任务,m
	PROGRESS_STATUS_REFUSE_INVITE = "refuse_invite"; //拒绝邀请,m
	PROGRESS_STATUS_APPROVE_INVITE = "approve_invite"; //同意邀请参加任务,g
	PROGRESS_STATUS_APPLY_COMPLETE = "apply_complete"; //申请完成任务/上传凭证,g
	PROGRESS_STATUS_JOB_UPLOAD = "job_upload"; //上传凭证,g
	PROGRESS_STATUS_COMPLETE = "complete"; //完成任务,g
	PROGRESS_STATUS_REFUSE_COMPLETE = "refuse_complete"; //拒绝完成任务,m
	PROGRESS_STATUS_MASTER_MARK = "master_mark"; //主态评价,m
	PROGRESS_STATUS_GUEST_MARK = "guest_mark"; //客态评价,g
	PROGRESS_STATUS_KICK_MEMBER = "kick_member"; //服务者超时未完成任务,g
	PROGRESS_STATUS_COMPLETE_AUTO = "complete_auto"; //需求者同意完成超时,m
)

const(
	ProofTypeNone  = iota //无凭证
	ProofTypeMedia        //视频凭证
	ProofTypeOther        //其他凭证
	ProofTypePic        //图片凭证
)

const(
	JOB_NOT_USE  = iota 	//未通过邀请
	JOB_NOEXIST				//1:岗位不存在(拒绝发放)
	JOB_INVALID				//2:岗位失效(拒绝发放)
	JOB_PASS				//3:人岗关系已通过(继续执行发放流程)
	JOB_NOT_SIGN			//4:用户未签约(继续执行发放流程)
	JOB_UNUSUAL	 = 9999		//9999:系统异常

	//重发间隔时间
	JOINS_SMS_RESENDTIME = 24 //h
	//短信内容
	JOIN_SMS_CONTENT = "尊敬的用户您好，您有一条来自 %v公司 的工作安排待接受，请及时使用[爱员工小助手]处理，超时未处理或将影响您工作绩效的下发，请知悉。"
)