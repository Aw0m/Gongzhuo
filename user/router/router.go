package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wxProjectDev/user/controllers"
)

func UserRouter() http.Handler {
	e := gin.New()
	e.Use(gin.Recovery(), gin.Logger())

	v1 := e.Group("/user")
	{
		v1.GET("/", func(c *gin.Context) {
			c.JSON(
				http.StatusOK,
				gin.H{
					"code":  http.StatusOK,
					"error": "Welcome Service User!",
				},
			)
		})

		// 登录
		v1.POST("/login", controllers.ConLogin)

		// 更新用户名
		v1.POST("/userName", controllers.ConUpdateUserName)

		// 获得用户名
		v1.GET("/userName", controllers.ConGetUserName)
	}

	v2 := e.Group("/team")
	{
		v2.POST("/createTeam", controllers.ConCreateTeam)

		v2.POST("/updateTeam", controllers.ConUpdateTeam)

		v2.GET("/selectTeam", controllers.ConSelectTeam)

		v2.GET("/selectMembers", controllers.ConSelectMembers)

		v2.GET("/selectTeamCode", controllers.ConSelectTeamCode)

		v2.POST("/updateTeamCode", controllers.ConUpdateTeamCode)

		v2.POST("/joinTeam", controllers.ConJoinTeam)

		// 获得指定用户的所有已加入的team信息
		v2.GET("/selectAllTeams", controllers.ConSelectAllTeams)

		// 增加新的admin
		v2.POST("/addAdmin", controllers.ConAddAdmin)
	}
	return e
}
