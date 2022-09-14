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
type App struct {
	Host string
	Port string
}

var (
	MyDB    MySQL
	AppHost App
)

type Config struct {
	SQL MySQL
	APP App
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

	//app
	app := iniData.Section("app")
	AppHost.Host = app.Key("host").String()
	AppHost.Port = app.Key("port").String()

	hostApp := &App{
		Host: AppHost.Host,
		Port: AppHost.Port,
	}

	//sql
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
		APP: *hostApp,
	}

}
