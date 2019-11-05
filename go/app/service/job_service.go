package job

import (
	"iQuest/app/constant"
	"math/rand"
	"time"
)

func GetFocusCount(typeValue string) int32 {
	randNum := 0
	switch {
	case typeValue == constant.TASK_FOCUS_NUM_DETAIL:
		randNum = RandInt(1, 10)
	case typeValue == constant.TASK_FOCUS_NUM_APPLY:
		randNum = RandInt(10, 20)
	default:
		randNum = 0
	}
	return int32(randNum)
}

func GetPrismaPageParam(pageNumber int, pageItem *int) (int, int) {
	size := int(*pageItem)
	//页码为空或小于1的处理
	if int(pageNumber) == 0 || pageNumber < 1 {
		pageNumber = 1
	}
	if size == 0 || size < 1 {
		size = 10
	}
	skip := (pageNumber - 1) * size
	return size, skip
}

func RandInt(min, max int) int {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	rand.Seed(time.Now().UnixNano()) // UnixNano()表示纳秒
	return rand.Intn(max-min) + min
}

func DateTimeToTimestamp(datetime string) *int {
	tm2, err := time.Parse(constant.DateTimeLayoutWithTimeZone, datetime)
	if err != nil {
		errTimestamp := -1
		return &errTimestamp
	}
	tm3 := int(tm2.Unix())
	return &tm3
}

//获取传入的时间所在月份的第一天，即某月第一天的0点。如传入time.Now(), 返回当前月份的第一天0点时间。
func GetFirstDateOfMonth(d time.Time) time.Time {
	d = d.AddDate(0, 0, -d.Day()+1)
	return GetZeroTime(d)
}

//获取某一天的0点时间
func GetZeroTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}
