package lib

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
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

func Md5(input []byte) string {
	hash := md5.Sum(input)
	return hex.EncodeToString(hash[:])
	//hash := md5.New()
	//
	////Get the 16 bytes hash
	//hash.Write(input)
	//hashInBytes := hash.Sum(nil)[:16]
	//
	////Convert the bytes to a string
	//
	//return hex.EncodeToString(hashInBytes)
}
