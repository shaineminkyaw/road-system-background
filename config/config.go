package config

import (
	"crypto/rsa"
	"io/ioutil"
	"log"
	"os"

	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
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
	Loggers    *logrus.Logger
)

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

	//private key
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

	//public key
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

	//logging
	Loggers = LogConf()

}

func LogConf() *logrus.Logger {
	file, err := os.OpenFile("logs/access.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("err %v\n", err)
	}
	// defer file.Close()

	log := logrus.New()
	log.Out = file

	//Set logs level
	log.SetLevel(logrus.GetLevel())

	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	return log
}
