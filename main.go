package main

import (
	"github.com/garfieldlw/live-url/routers"
	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()

	routers.LoadRouter(engine)

	engine.Run(":3005")
}
