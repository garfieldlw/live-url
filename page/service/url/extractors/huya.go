package extractors

import (
	"live-url/lib/call"
	"regexp"
	"errors"
	"encoding/json"
	"fmt"
	"strings"
)

type LiveHuya struct {
}

func (live *LiveHuya) Do(uri string) (*LiveInfo, error) {
	body, errRequest := call.HttpGetFake(uri, nil)
	if errRequest != nil {
		return nil, errRequest
	}

	reg, errReg := regexp.Compile("\"stream\": ({.+?})\\s*};")
	if errReg != nil {
		return nil, errReg
	}

	content := reg.FindString(string(body))
	if len(content) == 0 {
		return nil, errors.New("offline")
	}

	content = content[9 : len(content)-2]

	baseMap := new(HuyaBase)
	errBase := json.Unmarshal([]byte(content[:]), &baseMap)
	if errBase != nil {
		return nil, errBase
	}

	if baseMap.Status != 200 {
		return nil, errors.New("")
	}

	if len(baseMap.Data) == 0 {
		return nil, errors.New("")
	}

	infoData := baseMap.Data[0]
	if infoData == nil {
		return nil, errors.New("")
	}

	info := new(LiveInfo)
	if infoData.GameLiveInfo != nil {
		info.RoomName = infoData.GameLiveInfo.RoomName
		info.OwnerName = infoData.GameLiveInfo.Nick
	}

	if infoData.GameStreamInfo != nil && len(infoData.GameStreamInfo) > 0 {
		for _, item := range infoData.GameStreamInfo {
			if item == nil {
				continue
			}

			if item.SCdnType != "WS" {
				continue
			}

			info.RealUrl = fmt.Sprintf("%s/%s.%s?%v", item.SHlsUrl, item.SStreamName, item.SHlsUrlSuffix, strings.ReplaceAll(item.SHlsAntiCode, "amp;", ""))
		}
	}

	return info, nil
}

type HuyaBase struct {
	Status int32       `json:"status"`
	Msg    string      `json:"msg"`
	Data   []*HuyaData `json:"data"`
}

type HuyaData struct {
	GameLiveInfo   *HuyaGameLiveInfo     `json:"gameLiveInfo"`
	GameStreamInfo []*HuyaGameStreamInfo `json:"gameStreamInfoList"`
}

type HuyaGameLiveInfo struct {
	RoomName string `json:"roomName"`
	Nick     string `json:"nick"`
}

type HuyaGameStreamInfo struct {
	SCdnType      string `json:"sCdnType"`
	SStreamName   string `json:"sStreamName"`
	SHlsUrl       string `json:"sHlsUrl"`
	SHlsUrlSuffix string `json:"sHlsUrlSuffix"`
	SHlsAntiCode  string `json:"sHlsAntiCode"`
}
