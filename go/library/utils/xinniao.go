package utils

import (
	"iQuest/app/constant"
)

//判断那几个b是不是薪鸟
func IsInXinNiaoCompanyIdArr(companyId int32) bool{
	for i := 0; i < len(constant.XINNIAO_COMPANY_ID_ARR); i++ {
		if constant.XINNIAO_COMPANY_ID_ARR[i] == companyId {
			return true
		}
	}
	return false
}


//获取显示岗位还是任务
func GetWorkTitle(company_id int32) string {
	if(IsInXinNiaoCompanyIdArr(company_id)){
		return "任务"
	}
	return "岗位"
}


