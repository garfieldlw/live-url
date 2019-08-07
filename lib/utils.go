package lib

import (
	"github.com/gin-gonic/gin"
	"fmt"
)

func ValidateGetData(keyMap map[string]string, c *gin.Context) (map[string]string, error) {
	returnMap := make(map[string]string, len(keyMap))
	for keyName, defaultVal := range keyMap {
		data, isExist := c.GetQuery(keyName)
		if defaultVal == "not-nil" {
			if !isExist {
				return nil, fmt.Errorf("%s not found", keyName)
			}
		}
		if !isExist {
			returnMap[keyName] = defaultVal
		} else {
			returnMap[keyName] = data
		}
	}
	return returnMap, nil
}
