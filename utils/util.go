package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/scrypt"
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
	UserID uint64
	jwt.StandardClaims
}

func GetAccessToken(userID uint64) (string, error) {
	//
	unixTime := time.Now()
	expTime := unixTime.Add(60 * 15) // 15min
	key := "mysecret"

	claims := &AccessTokenCustomClaim{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  unixTime.Unix(),
			ExpiresAt: expTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenData, err := token.SignedString([]byte(key))
	if err != nil {
		log.Println(err.Error())
	}
	return tokenData, err
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
