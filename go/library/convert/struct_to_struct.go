package convert

import (
	"github.com/fatih/structs"
	"github.com/mitchellh/mapstructure"
)

func Struct2Struct(input interface{}) (output interface{}, err error) {
	//	tmp := *input
	switch input.(type) {
	case struct{}:
		tmpMap := structs.Map(input)
		err := mapstructure.Decode(tmpMap, &output)
		if err != nil {
			return nil, err
		}
		return output, err
	}
	return
}
