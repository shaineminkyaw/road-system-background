package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/scrypt"
	"gopkg.in/ini.v1"
	"gorm.io/gorm"
)

//@@@hash password
func HashPassword(pass string) (string, error) {
	//
	salt := make([]byte, 64)
	_, err := rand.Read(salt)
	if err != nil {
		return "", nil
	}

	unHash, err := scrypt.Key([]byte(pass), salt, 32768, 8, 1, 32)
	if err != nil {
		log.Fatalf("error on crypting password %v", err.Error())
	}
	hashed := fmt.Sprintf("%v.%v", hex.EncodeToString(unHash), hex.EncodeToString(salt))
	return hashed, nil

}

//validate hash
func ValidateHashedPassword(hash, plainPassword string) (bool, error) {
	//
	splitHash := strings.Split(hash, ".")
	salt, err := hex.DecodeString(splitHash[1])
	if err != nil {
		log.Fatalf("error on : %v", err.Error())
	}
	plainHashed, err := scrypt.Key([]byte(plainPassword), salt, 32768, 8, 1, 32)
	if err != nil {
		log.Fatalf("error on %v", err.Error())
	}

	return (hex.EncodeToString(plainHashed)) == (splitHash[0]), nil
}

//@@generate token

type AccessTokenCustomClaim struct {
	UserID uint64 `json:"user_id"`
	jwt.StandardClaims
}

//generateToken
func GetAccessToken(userID uint64, key *rsa.PrivateKey) (string, error) {
	//
	unixTime := time.Now().Unix()
	expTime := unixTime + (60 * 60 * 48) //2 day

	claims := &AccessTokenCustomClaim{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: expTime,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenData, err := token.SignedString(key)
	if err != nil {
		log.Println(err.Error())
	}
	return tokenData, err
}

type accessTokenCustomClaim struct {
	UserID uint64 `json:"user_id"`
	jwt.StandardClaims
}

//ValiadetAccessToken
func ValidateAccessToken(token string, key *rsa.PublicKey) (*accessTokenCustomClaim, error) {
	claims := &accessTokenCustomClaim{}

	accessToken, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return key, nil
	})

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	if !accessToken.Valid {
		log.Println("token is invalid")
	}
	claim, ok := accessToken.Claims.(*accessTokenCustomClaim)
	if !ok {
		log.Println("token is valid but could not parse claims")
	}

	// fmt.Println("finish ::::")
	return claim, nil
}

//start date end date
func BetweenDate(startDate, endDate string) func(*gorm.DB) *gorm.DB {
	//
	return func(db *gorm.DB) *gorm.DB {
		if len(startDate) > 0 {
			start, err := time.Parse(time.RFC1123, startDate)
			if err != nil {
				log.Println(err.Error())
			}
			db = db.Where("created_at => ?", start.Format("2006-01-02 15:04:05"))
		}
		if len(endDate) > 0 {
			end, err := time.Parse(time.RFC1123, endDate)
			if err != nil {
				log.Println(err.Error())
			}
			db = db.Where("created_at <= ?", end.Format("2006-01-02 15:04:05"))
		}
		return db
	}
}

func Paginate(page, pageSize int) func(*gorm.DB) *gorm.DB {
	//
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}
		if pageSize == 0 {
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

//
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
)

type Config struct {
	SQL       MySQL
	APP       App
	Private   *rsa.PrivateKey
	Public    *rsa.PublicKey
	SecretKey string
}

func File1() *Config {
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

	//rsa
	rsa := iniData.Section("rsa")

	key := rsa.Key("Secret_Key").String()
	prv := rsa.Key("Private_Key").String()
	private, err := ioutil.ReadFile(prv)
	if err != nil {
		log.Println(err.Error())
	}
	p, err := jwt.ParseRSAPrivateKeyFromPEM(private)
	if err != nil {
		log.Println(err.Error())
	}

	pub := rsa.Key("Public_Key").String()
	public, err := ioutil.ReadFile(pub)
	if err != nil {
		log.Println(err.Error())
	}

	pu, err := jwt.ParseRSAPublicKeyFromPEM(public)
	if err != nil {
		log.Println(err.Error())
	}

	db := &MySQL{
		Host:     MyDB.Host,
		Port:     MyDB.Port,
		DB:       MyDB.DB,
		User:     MyDB.User,
		Password: MyDB.Password,
	}

	return &Config{
		SQL:       *db,
		APP:       *hostApp,
		Private:   p,
		Public:    pu,
		SecretKey: key,
	}

}
