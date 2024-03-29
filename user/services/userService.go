package services

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
	"wxProjectDev/public/utils"
	"wxProjectDev/user/daos"
)

// ServiceLogin 用户登录服务，生成Token。并且在数据库中检索该用户是否已经注册，如果没有则还会在数据库中创建该用户
func ServiceLogin(code, userName string, context *gin.Context) {
	url := "https://api.weixin.qq.com/sns/jscode2session?appid=wxf02df1739164ba92&secret=35ca848d9e2f3ba0b88ff9eedc4b2eb2&js_code=" + code + "&grant_type=authorization_code"
	var httpCode int
	var resError any
	var token string
	var openid string
	createUser := false

	timeStart := time.Now()
	// 通过code获取用户的唯一标识符openid
	if res, err := http.Get(url); err != nil {
		httpCode = http.StatusBadRequest
		resError = "code is not good"
	} else {
		// 解析微信服务端的response，获得openid并查询是否已经存入数据库，如果没有则在数据库中生成一个user
		body, _ := utils.ParseResponse(res)
		httpCode = http.StatusOK
		resError = body["errcode"]
		openidAny := body["openid"]

		var ok bool
		if openid, ok = openidAny.(string); ok {
			// 生成token。如果数据库里没有该用户，则在该数据库生成该user
			token = utils.CreateToken(openid, openid)
			if _, err := daos.SelectUser(openid); err == sql.ErrNoRows {
				log.Printf("生成用户：%s\n", userName)
				err = daos.CreateUser(openid, userName)
				createUser = true
			}
			if err != nil {
				log.Println(err)
			}
		} else {
			log.Println("openid 不为 string！")
		}
	}
	fmt.Println("本次执行时间为:", time.Since(timeStart))
	fmt.Println("token: ", token)
	context.JSON(
		httpCode,
		gin.H{
			"error":      resError,
			"token":      token,
			"openid":     openid,
			"createUser": createUser,
		},
	)
}

// ServiceUpdateUserName 更新用户的昵称
func ServiceUpdateUserName(openid, userName string, context *gin.Context) {
	if err := daos.UpdateUser(openid, userName); err == nil {
		context.JSON(
			http.StatusOK,
			gin.H{
				"msg": "ok",
			},
		)
	} else {
		context.JSON(
			http.StatusServiceUnavailable,
			gin.H{
				"msg": "error",
			},
		)
	}
}

// ServiceSelectUserName 获得用户的昵称
func ServiceSelectUserName(openid string, context *gin.Context) {
	if user, err := daos.SelectUser(openid); err == nil {
		context.JSON(
			http.StatusOK,
			gin.H{
				"msg":      "ok",
				"userName": user.UserName,
			},
		)
	} else {
		context.JSON(
			http.StatusServiceUnavailable,
			gin.H{
				"msg":      "error",
				"userName": "",
			},
		)
	}
}
