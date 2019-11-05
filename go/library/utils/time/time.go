package time

import (
	"iQuest/app/constant"
	"time"
)

//检查target是否大于source
func AfterTime(target string, source time.Time) (gt bool, err error) {
	t, err := time.Parse(constant.DateTimeLayout, target)
	if nil != err {
		return false, err
	}

	if t.After(source) {
		return true, nil
	}

	return false, nil
}

//当前东八区时间
func PRCNow() (datetime string) {
	return time.Now().Format(constant.DateTimeLayout)
}

//解析时间到 Y-m-d H:i:s的格式
func ParseToPRCTime(input string) (time.Time, error) {
	t, err := time.Parse(constant.DateTimeLayout, input)
	if nil != err {
		return time.Time{}, err
	}

	return t, nil
}