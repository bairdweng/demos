package utils

import "strconv"

//*int(指针类型的int)转换为 *int32
func PointInt2PointInt32(input *int, defaultWhenNil int) (output *int32) {
	var i32 int32
	if input == nil {
		i32 = int32(defaultWhenNil)
	} else {
		i32 = int32(*input)
	}

	return &i32
}

//int 转换为 *int32
func Int2PointInt32(input int) (output *int32) {
	i32 := int32(input)
	return &i32
}

func Int322PointInt(input int32) (output *int) {
	i := int(input)
	return &i
}

func Int322PointString(input int32) (output *string) {
	str := Int322String(input)
	return &str
}

func Int322String(input int32) (output string) {
	str := strconv.FormatInt(int64(input), 10)
	return str
}


func Int642PointString(input int64) (output *string) {
	str := Int642String(input)
	return &str
}
func Int642String(input int64) (output string) {
	str := strconv.FormatInt(int64(input), 10)
	return str
}

func String2PointInt32(input string) (output int32) {
	defaultVal := int32(0)
	if input == "" {
		return defaultVal
	}

	i64, err := strconv.ParseInt(input, 10, 32)
	if err != nil {
		return defaultVal
	}

	i32 := int32(i64)
	return i32
}

