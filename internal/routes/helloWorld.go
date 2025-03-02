package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewHelloWorldRouter(group *gin.RouterGroup) {
	group.GET("/helloworld", Helloworld)
}

// PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /helloworld [get]
func Helloworld(g *gin.Context) {
	g.JSON(http.StatusOK, "helloworld")
}
