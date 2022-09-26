package ds

import (
	"fmt"
	"log"

	"github.com/shaineminkyaw/road-system-background/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DataSource struct {
	Sql *gorm.DB
}

var DB *gorm.DB

func ConnectToDB(host, port, dbname, dbuser, dbpassword string) *gorm.DB {

	db_host := host
	db_port := port
	db_name := dbname
	db_user := dbuser
	db_password := dbpassword

	// DB.Exec("create database ?", dbname)

	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local", db_user, db_password, db_host, db_port, db_name)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("error on connecting to database!")
	} else {
		log.Printf("Connected Database :::")
	}

	DB = db
	db.AutoMigrate(
		&model.Admin{},
		&model.AdminRole{},
		&model.Roles{},
		&model.CasbinPloicy{},
		&model.BetTable{},
		&model.UserLog{},
	)

	return db
}
