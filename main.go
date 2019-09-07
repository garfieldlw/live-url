package main

import (
	"github.com/gin-gonic/gin"
	"live-url/routers"
)

func main() {
	engine := gin.Default()

	routers.LoadRouter(engine)

	engine.Run(":3005")
}
