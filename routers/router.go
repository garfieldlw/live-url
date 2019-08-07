package routers

import (
	"github.com/gin-gonic/gin"
	"live-url/controllers"
)

func LoadRouter(engine *gin.Engine) {
	urlCon := new(controllers.UrlController)
	urlGroup := engine.Group("/api/live")
	urlGroup.GET("/getliveurl", urlCon.GetLiveUrl)
}
