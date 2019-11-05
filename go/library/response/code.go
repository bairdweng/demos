package response

// SUCCESS = 0, ERROR = 1 如此类推
const (
	Success = iota
	Error
	ValidateError
	NotFound
	ParamError
	DBError
	IllegalOperate
	TransactionError
	ConnectError
	MobileNotNull = "手机号不能为空"
	SmsOverNumber = "每天最多获取%s次验证码"
	SmsFrequent   = "获取验证码过于频繁，请稍后再试"
	ErrorString   = "1"
	ParamErrorMsg = "非法参数"
	RequestError  = "请求接口或参数相关错误"
)
