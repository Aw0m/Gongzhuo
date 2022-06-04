package controllers

import (
	"github.com/gin-gonic/gin"
	"log"
	"wxProjectDev/user/services"
)

func ConLogin(context *gin.Context) {
	code := context.PostForm("code")
	userName := context.PostForm("userName")
	log.Println("code is ", code, userName)
	services.ServiceLogin(code, userName, context)
}

func ConUpdateUserName(context *gin.Context) {
	openid := context.GetHeader("userID")
	userName := context.PostForm("userName")
	services.ServiceUpdateUserName(openid, userName, context)
}

func ConGetUserName(context *gin.Context) {
	openid := context.GetHeader("userID")
	services.ServiceSelectUserName(openid, context)
}
