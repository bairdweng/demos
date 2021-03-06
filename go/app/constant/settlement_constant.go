package constant

//状态
const (
	CONFIRM_STATUS  = 1
	UN_CONFIRM_STATUS  = 0
	//短信内容
	SETTLEMENT_SMS_CONTENT = "尊敬的用户您好，您参与的 %v 岗位/任务有一条新的结算单待确认，请及时使用[薪鸟小助手]确认结算单，以便我们及时为您结算。"

	READ_TYPE_INVITE = 1
	READ_TYPE_SETTLEMENT = 2
	XINNIAO_APPID = "xinniao"
	XINNIAO_RESUME_ID = "[{\"id\": 1},{\"id\": 2},{\"id\": 3},{\"id\": 4},{\"id\": 5},{\"id\": 6}]"
)

// 92
//var XINNIAO_COMPANY_ID_ARR = []int32{10001209,10001097,100001687,10002277,10002276,10002275}
//预发布
//var XINNIAO_COMPANY_ID_ARR = []int32{10001209,10001097,10001213,10001215,10001214}
//生产
var XINNIAO_COMPANY_ID_ARR = []int32{10001689,10001708,10001738,10001756,10001779,10001835,10001836,10001837,10002063,10002080,10002111,10002112,10002131,10002132,10002225,10002244,10002273,10002312,10002331,10002168}
