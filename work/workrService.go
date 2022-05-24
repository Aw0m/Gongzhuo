package work

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ServiceCreateReport(userID, teamID, done, toDo, problem string, context *gin.Context) {
	fmt.Printf("done:\n%s\n", done)
	fmt.Printf("\ntoDO:\n%s\n", toDo)
	fmt.Printf("\nproblem:\n%s\n", problem)

	repID, err := createReport(userID, teamID, done, toDo, problem)
	if err != nil {
		context.JSON(
			http.StatusBadRequest,
			gin.H{
				"msg":   "error",
				"repID": repID,
			},
		)
	} else {
		context.JSON(
			http.StatusOK,
			gin.H{
				"msg":   "ok",
				"repID": repID,
			},
		)
	}
}
