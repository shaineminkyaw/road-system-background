package config

import (
	"log"
	"os"

	"gopkg.in/ini.v1"
)

type MySQL struct {
	Host     string
	Port     string
	DB       string
	User     string
	Password string
}

var MyDB MySQL

type Config struct {
	SQL MySQL
}

func init() {
	//

}

func Init() *Config {
	//
	iniFile := "config/config.ini"

	args := os.Args

	if len(args) > 1 {
		iniFile = args[1]
	}

	iniData, err := ini.Load(iniFile)
	if err != nil {
		log.Println(err.Error())
	}

	sql := iniData.Section("mysql")
	MyDB.Host = sql.Key("host").String()
	MyDB.Port = sql.Key("port").String()
	MyDB.DB = sql.Key("db").String()
	MyDB.User = sql.Key("user").String()
	MyDB.Password = sql.Key("password").String()

	db := &MySQL{
		Host:     MyDB.Host,
		Port:     MyDB.Port,
		DB:       MyDB.DB,
		User:     MyDB.User,
		Password: MyDB.Password,
	}

	return &Config{
		SQL: *db,
	}

}
