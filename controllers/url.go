package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"live-url/lib"
	"live-url/page/service"
)

type UrlController struct {
}

func (controller *UrlController) GetLiveUrl(c *gin.Context) {
	validateMap := map[string]string{
		"uri":  "not-nil",
		"room": "not-nil",
	}
	var res string
	resMap, errGet := lib.ValidateGetData(validateMap, c)
	if errGet != nil {
		c.String(http.StatusBadRequest, "")
		return
	}

	res, errRes := service.GetLiveUrl(resMap["uri"], resMap["room"])
	if errRes != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}
	c.String(http.StatusOK, res)
	return
}
