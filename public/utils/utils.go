package utils

import (
	"encoding/json"
	"github.com/golang-jwt/jwt"
	"io/ioutil"
	"net/http"
	"time"
)

type MyCustomClaims struct {
	UserId   string `json:"userId"`
	UserName string `json:"userName"`
	jwt.StandardClaims
}

var mySigningKey = []byte("dev123")

const duration = 24 * 30

func CreateToken(userId string, userName string) string {
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

func ParseResponse(response *http.Response) (map[string]interface{}, error) {
	var result map[string]interface{}
	body, err := ioutil.ReadAll(response.Body)
	if err == nil {
		err = json.Unmarshal(body, &result)
	}

	return result, err
}
