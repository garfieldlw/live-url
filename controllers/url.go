package controllers

import (
	"fmt"
	"github.com/garfieldlw/live-url/lib"
	"github.com/garfieldlw/live-url/page/service/url"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UrlController struct {
}

func (controller *UrlController) GetLiveUrl(c *gin.Context) {
	validateMap := map[string]string{
		"uri": "not-nil",
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
