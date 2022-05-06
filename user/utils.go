package user

import (
	"github.com/golang-jwt/jwt"
	"time"
)

type MyCustomClaims struct {
	UserId   string `json:"userId"`
	UserName string `json:"userName"`
	jwt.StandardClaims
}

var mySigningKey = []byte("dev123")

const duration = 24 * 30

func createToken(userId string, userName string) string {
	myToken := jwt.NewWithClaims(jwt.SigningMethodHS256, MyCustomClaims{
		userId,
		userName,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * duration).Unix(),
			NotBefore: time.Now().Unix(),
			Issuer:    "test",
		},
	})
	ss, _ := myToken.SignedString(mySigningKey)
	return ss
}
