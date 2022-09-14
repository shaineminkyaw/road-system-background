package factory

import (
	"log"

	"github.com/shaineminkyaw/road-system-background/ds"
	"github.com/shaineminkyaw/road-system-background/model"
	"github.com/shaineminkyaw/road-system-background/utils"
)

func FactoryAdmin() error {
	//
	password, err := utils.HashPassword("password-123")
	if err != nil {
		log.Println(err.Error())
	}

	//create default admin
	admin := &model.Admin{
		UserName: "admin",
		Email:    "administrator@gmail.com",
		Password: password,
	}

	err = ds.DB.Model(&model.Admin{}).Create(&admin).Error
	if err != nil {
		log.Println(err.Error())
		return err
	}

	//assign role
	role := &model.AdminRole{
		UserID: admin.ID,
		RoleID: 1,
	}

	err = ds.DB.Model(&model.AdminRole{}).Create(&role).Error
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}
