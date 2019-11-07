package routers

import (
	"github.com/garfieldlw/live-url/controllers"
	"github.com/gin-gonic/gin"
)

func LoadRouter(engine *gin.Engine) {
	urlCon := new(controllers.UrlController)
	urlGroup := engine.Group("/api/live")
	urlGroup.GET("/getliveurl", urlCon.GetLiveUrl)
}
