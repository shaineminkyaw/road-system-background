package factory

import (
	"log"

	"github.com/shaineminkyaw/road-system-background/model"

	"github.com/shaineminkyaw/road-system-background/ds"
)

func FactoryRole() error {
	//
	role := &model.Roles{
		Slug: "admin",
		Name: "Administrator",
	}

	err := ds.DB.Model(&model.Roles{}).Create(&role).Error
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}
