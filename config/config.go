package config

import (
	"crypto/rsa"
	"io/ioutil"
	"log"
	"os"

	"github.com/golang-jwt/jwt"
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
	MyDB       MySQL
	AppHost    App
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
	SecretKey  string
)

// type Config struct {
// 	SQL       MySQL
// 	APP       App
// 	Private   *rsa.PrivateKey
// 	Public    *rsa.PublicKey
// 	SecretKey string
// }

func init() {
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

	//sql
	sql := iniData.Section("mysql")
	MyDB.Host = sql.Key("host").String()
	MyDB.Port = sql.Key("port").String()
	MyDB.DB = sql.Key("db").String()
	MyDB.User = sql.Key("user").String()
	MyDB.Password = sql.Key("password").String()

	//rsa
	rsa := iniData.Section("rsa")
	SecretKey = rsa.Key("Secret_Key").String()
	//
	prv := rsa.Key("Private_Key").String()
	private, err := ioutil.ReadFile(prv)
	if err != nil {
		log.Println(err.Error())
	}
	p, err := jwt.ParseRSAPrivateKeyFromPEM(private)
	if err != nil {
		log.Println(err.Error())
	}
	PrivateKey = p

	//
	pub := rsa.Key("Public_Key").String()
	public, err := ioutil.ReadFile(pub)
	if err != nil {
		log.Println(err.Error())
	}

	pu, err := jwt.ParseRSAPublicKeyFromPEM(public)
	if err != nil {
		log.Println(err.Error())
	}
	PublicKey = pu
}

// func File() *Config {
// 	//
// 	iniFile := "config/config.ini"

// 	args := os.Args

// 	if len(args) > 1 {
// 		iniFile = args[1]
// 	}

// 	iniData, err := ini.Load(iniFile)
// 	if err != nil {
// 		log.Println(err.Error())
// 	}

// 	//app
// 	app := iniData.Section("app")
// 	AppHost.Host = app.Key("host").String()
// 	AppHost.Port = app.Key("port").String()

// 	hostApp := &App{
// 		Host: AppHost.Host,
// 		Port: AppHost.Port,
// 	}

// 	//sql
// 	sql := iniData.Section("mysql")
// 	MyDB.Host = sql.Key("host").String()
// 	MyDB.Port = sql.Key("port").String()
// 	MyDB.DB = sql.Key("db").String()
// 	MyDB.User = sql.Key("user").String()
// 	MyDB.Password = sql.Key("password").String()

// 	//rsa
// 	rsa := iniData.Section("rsa")

// 	key := rsa.Key("Secret_Key").String()
// 	prv := rsa.Key("Private_Key").String()
// 	private, err := ioutil.ReadFile(prv)
// 	if err != nil {
// 		log.Println(err.Error())
// 	}
// 	p, err := jwt.ParseRSAPrivateKeyFromPEM(private)
// 	if err != nil {
// 		log.Println(err.Error())
// 	}

// 	pub := rsa.Key("Public_Key").String()
// 	public, err := ioutil.ReadFile(pub)
// 	if err != nil {
// 		log.Println(err.Error())
// 	}

// 	pu, err := jwt.ParseRSAPublicKeyFromPEM(public)
// 	if err != nil {
// 		log.Println(err.Error())
// 	}

// 	db := &MySQL{
// 		Host:     MyDB.Host,
// 		Port:     MyDB.Port,
// 		DB:       MyDB.DB,
// 		User:     MyDB.User,
// 		Password: MyDB.Password,
// 	}

// 	return &Config{
// 		SQL:       *db,
// 		APP:       *hostApp,
// 		Private:   p,
// 		Public:    pu,
// 		SecretKey: key,
// 	}

// }
