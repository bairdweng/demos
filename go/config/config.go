package config

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// Viper Viper
var Viper = viper.GetViper()

func init() {

	/*viper.SetConfigName("config")
	viper.AddConfigPath("./config")*/

	viper.SetConfigType("env")
	viper.SetConfigName("env")
	viper.SetConfigFile(".env")
	viper.AddConfigPath("./")

	//replacer := strings.NewReplacer(".", "_")
	//viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()
	err := viper.ReadInConfig() // 搜索路径，并读取配置数据
	if err != nil {
		log.Panic().Str("err", err.Error()).Msg("Fatal error config file")
	}
	setDefault()
}

func setDefault() {
	// 跨域设置
	viper.SetDefault("ALLOW_ORIGINS", "*")
	viper.SetDefault("ALLOW_HEADERS", "*")

	// 从客户端接收的请求头
	viper.SetDefault("HEADER_AUTH", "Ayg-Sessionid")
	viper.SetDefault("HEADER_AUTH_RAW", "AYG_SESSIONID")
	viper.SetDefault("HEADER_COMPANY_ID", "Company-Id")
	viper.SetDefault("HEADER_PROFILE_ID", "Profile-Id")
	viper.SetDefault("HEADER_COMPANY_NAME", "Company-Name")

	viper.SetDefault("COOKIE_AUTH", "x-access-token")
	viper.SetDefault("COOKIE_PROFILE_ID", "x-sec-profile")
	viper.SetDefault("COOKIE_COMPANY_ID", "x-sec-subject-company-id")
	viper.SetDefault("COOKIE_COMPANY_NAME", "x-sec-subject-company-name")
	viper.SetDefault("COOKIE_APP_ID", "x-sec-lvl1subject-app-id")
	viper.SetDefault("COOKIE_APP_NAME", "x-sec-lvl1subject-app-name")

	// 日志相关
	viper.SetDefault("LOG_LEVEL", 0)                         // 日志等级 Debug:0 Info:1 Warn:2 Error:3 Fatal:4 Panic:5 No:6 Disabled:7
	viper.SetDefault("LOG_FILENAME", "logs/ishouru-job.log") // 日志文件名
	viper.SetDefault("LOG_MAX_SIZE", 100)                    // 日志单文件大小 mb
	viper.SetDefault("LOG_MAX_AGE", 7)                       // 日志天数
	viper.SetDefault("LOG_MAX_BACKUPS", 10)                  // 日志文件数
	viper.SetDefault("LOG_JSON", false)                      // 日志保存JSON
	viper.SetDefault("LOG_CONSOLE", false)                   // 日志输出到stdout

	// 发放相关
	viper.SetDefault("PAY_FOR_UNCONFIRMED", true)
	viper.SetDefault("PAYROLL_TEMPLATE", "./assets/payroll_template.xlsx")

	// passport接口地址
	viper.SetDefault("PASS_PORT_URL", "http://passport-service/user/")

	// 签名接口
	viper.SetDefault("SIGN_API_URL", "http://econtract")

	// 获取app_id 地址
	viper.SetDefault("COMPANY_APP_URL", "http://sysmgr-web")

	// 短信接口
	viper.SetDefault("COMMON_APP_URL", "http://cloud-apigateway")
	viper.SetDefault("SMS_APP_ID", "21")

	// 文件上传地址
	viper.SetDefault("UPLOAD_PATH", "./upload/")

	// 岗位模版审核回调地址
	viper.SetDefault("JOB_TEMPLATE_AUDIT_CALLBACK_URL", "/callback/jobTemplate/audit")

	// 添加用户接口
	viper.SetDefault("USER_SERVICE_VERIFY", "http://user-service/user/add-realname-user")
	viper.SetDefault("USER_IMPORT", "http://user-service/user/add-realname-user-list")
	viper.SetDefault("USER_POSITION_FLOW_COUNT", "http://datacenter/post-evidence/user-post/flow/count")
	viper.SetDefault("USER_FLOW_COUNT", "http://datacenter/post-evidence/user/flow/count")
	viper.SetDefault("SERVICE_NAME", "http://contract-web/service-mgr/get-service-type-options")

	// passport用户接口地址
	viper.SetDefault("PASS_PORT_URL_USER", "http://user-service")

	// 发放相关接口
	viper.SetDefault("POST_NOTIFY_URL", "http://console-dlv")
	viper.SetDefault("POST_EVIDENCE_URL", "http://datacenter")

	// 消息服务地址
	viper.SetDefault("MESSAGE_URL", "http://ishouru-message")
}
