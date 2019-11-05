package commonly

import (
	"iQuest/app/model/commonly"
	"iQuest/db"
)

func Get() ([]commonly.CommonlyUsed, error) {
	var commonly_used []commonly.CommonlyUsed
	if err := db.Get().Last(&commonly_used).Error; err != nil {
		return nil, err
	}
	return commonly_used, nil
}

func Create(commonly *commonly.CommonlyUsed) error {

	if err := db.Get().Create(&commonly).Error; err != nil {
		return err
	}
	return nil
}

func Delete(id string) error {
	if err := db.Get().Where("id = (?)", id).Delete(commonly.CommonlyUsed{}).Error; err != nil {
		return err
	}
	return nil
}
