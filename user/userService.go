package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

func ServiceLogin(code string, context *gin.Context) {
	url := "https://api.weixin.qq.com/sns/jscode2session?appid=wx351cbd21ad881ad9&secret=eb43225a92bc2e7eb50a9555c0d38fbf&js_code=" + code + "&grant_type=authorization_code"
	var httpCode int
	var resError any
	var token string

	timeStart := time.Now()
	if res, err := http.Get(url); err != nil {
		httpCode = http.StatusBadRequest
		resError = "code is not good"
	} else {
		body, _ := parseResponse(res)
		httpCode = http.StatusOK
		resError = body["errcode"]

		openid := body["openid"]
		fmt.Println("openid:", openid)
		if openid, ok := openid.(string); ok {
			token = createToken(openid, openid)
		} else {
			log.Println("openid 不为 string！")
		}
	}
	fmt.Println("本次执行时间为:", time.Since(timeStart))
	fmt.Println("token: ", token)
	context.JSON(
		httpCode,
		gin.H{
			"error": resError,
			"token": token,
		},
	)
}
