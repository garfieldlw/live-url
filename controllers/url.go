package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"live-url/lib"
	"live-url/page/service/url"
	"fmt"
)

type UrlController struct {
}

func (controller *UrlController) GetLiveUrl(c *gin.Context) {
	validateMap := map[string]string{
		"uri":  "not-nil",
	}
	resMap, errGet := lib.ValidateGetData(validateMap, c)
	if errGet != nil {
		c.String(http.StatusBadRequest, "")
		return
	}

	res, errRes := url.GetLiveUrl(resMap["uri"])
	if errRes != nil || res == nil {
		fmt.Println(errRes)
		c.String(http.StatusInternalServerError, "")
		return
	}
	c.String(http.StatusOK, res.RealUrl)
	return
}
