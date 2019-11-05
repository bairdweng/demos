package constant

//work状态
const (
	WORK_STATUS_NORMAL  = 1
	WORK_STATUS_CLOSED  = 3
	WORK_STATUS_Reject  = 2
	WORK_STATUS_DELETED = 5
	WORK_STATUS_UNAUDIT = 7 //待审核
	//私密状态
	WORK_STATUS_PUBLIC  = 1
	WORK_STATUS_PRIVATE = 0

	//任务关注状态
	TASK_FOCUS_NUM_DETAIL = "detail"
	TASK_FOCUS_NUM_APPLY  = "apply"

	//任务/任务来源
	WORK_SOURCE_APPLY  = 1 //申请
	WORK_SOURCE_INVITE = 2 //邀请

)

const (
	WorkProgressUndefined = iota
	WorkProgressInviting //邀请中
	WorkProgressApplying //申请请中
	WorkProgressApprove //已同意参加
	WorkProgressApplyComplete //申请完成
	WorkProgressCompleted //已完成
	WorkProgressParticipantMark //已评分
	WorkProgressPublisherMark //发布者已评分
	WorkProgressBothMark //互评
	WorkProgressRejectApply //拒绝申请
	WorkProgressRejectComplete //拒绝完成
	WorkProgressKickOut //提出任务
)

const (
	WorkSourceUndefined   = iota
	WorkSourceApp         //app
	WorkSourceXiaoShan    //小善
	WorkSourceAppShuiChou //税筹

)

const (
	WorkTypePrivate = 0
	WorkTypePublic  = 1
)

//任务进程
//进行中状态(任务关系状态)
var JOB_STATUS_ING = []int32{3, 4, 5, 6, 7, 8, 10}

//任务已参加列表状态
var JOB_STATUS_JOINED = []int32{3, 4, 5, 6}

//可申请列表
var JOB_STATUS_CAN_APPLY = []int32{1, 3, 4, 5, 6, 7, 8, 10}

//完成任务状态
var JOB_STATUS_FINISH_TASK = []int32{5, 6, 7, 8}
