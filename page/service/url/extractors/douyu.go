package extractors

import (
	"encoding/json"
	"fmt"
	"github.com/garfieldlw/live-url/lib/call"
	"github.com/google/uuid"
	"github.com/robertkrimen/otto"
	"regexp"
	"strings"
	"time"
)

type DouyuH5EncInfo struct {
	Error int32             `json:"error"`
	Data  map[string]string `json:"data"`
}

type DouyuLiveBase struct {
	Error int32         `json:"error"`
	Msg   string        `json:"msg"`
	Data  DouyuDataInfo `json:"data"`
}

type DouyuDataInfo struct {
	RtmpUrl  string `json:"rtmp_url"`
	RtmpLive string `json:"rtmp_live"`
}

type LiveDouyu struct {
}

func (live *LiveDouyu) Do(uri string) (*LiveInfo, error) {
	roomId := extractRoomId(uri)
	body, errRequest := call.HttpGetFake("https://open.douyucdn.cn/api/RoomApi/room/"+roomId, nil)
	if errRequest != nil {
		return nil, errRequest
	}

	baseMap := make(map[string]interface{})
	errJson := json.Unmarshal(body, &baseMap)
	if errJson != nil {
		return nil, errJson
	}

	bodyH5enc, errH5enc := call.HttpGetFake("https://www.douyu.com/swf_api/homeH5Enc?rids="+roomId, nil)
	if errH5enc != nil {
		return nil, errH5enc
	}

	h5enc := new(DouyuH5EncInfo)
	errJson = json.Unmarshal(bodyH5enc, &h5enc)
	js_enc := h5enc.Data["room"+roomId]

	reg, errReg := regexp.Compile("function ub98484234\\(.+?\\Weval\\((\\w+)\\);")
	if errReg != nil {
		return nil, errReg
	}

	workflow := reg.FindString(js_enc)

	bodyCrypto, errCrypto := call.HttpGetFake("https://cdnjs.cloudflare.com/ajax/libs/crypto-js/3.1.9-1/crypto-js.min.js", nil)
	if errCrypto != nil {
		return nil, errCrypto
	}
	jsCrypto := string(bodyCrypto)

	jsDom := getDom()
	jsDebug := getDebug()
	jsPatch := getPatch(workflow)

	js := js_enc + jsCrypto + jsDom + jsDebug + jsPatch
	vm := otto.New()
	vm.Run(js)
	did := uuid.New().String()
	tt := fmt.Sprintf("%v", time.Now().Unix())
	ub98484234, _ := vm.Call("ub98484234", nil, roomId, did, tt)

	reg1, errReg := regexp.Compile("v=(\\d+)")
	if errReg != nil {
		return nil, errReg
	}

	v := reg1.FindString(ub98484234.String())
	if len(v) > 0 {
		v = v[2:]
	}

	reg2, errReg := regexp.Compile("sign=(\\w{32})")
	if errReg != nil {
		return nil, errReg
	}
	sign := reg2.FindString(ub98484234.String())
	if len(sign) > 0 {
		sign = sign[5:]
	}

	queryMap := make(map[string]string, 0)
	queryMap["v"] = v
	queryMap["did"] = did
	queryMap["tt"] = tt
	queryMap["sign"] = sign
	queryMap["cdn"] = ""
	queryMap["iar"] = "0"
	queryMap["ive"] = "0"
	queryMap["rate"] = "0"

	bodyLive, errLive := call.HttpPost("https://www.douyu.com/lapi/live/getH5Play/"+roomId, queryMap)
	if errLive != nil {
		return nil, errLive
	}

	douyuInfo := new(DouyuLiveBase)
	errJson = json.Unmarshal(bodyLive, &douyuInfo)

	info := new(LiveInfo)
	info.RealUrl = douyuInfo.Data.RtmpUrl + "?" + douyuInfo.Data.RtmpLive
	return info, nil
}

func extractRoomId(url string) (string) {
	index := strings.LastIndex(url, "/") + 1
	return url[index:]
}

func getDebug() string {
	return "var ub123 = ub98484234;" +
		"ub98484234 = function(p1, p2, p3) {{ " +
		" try {{ " +
		"		var resoult = ub123(p1, p2, p3);" +
		"	debug123.resoult123 = resoult;" +
		"	}} catch(e) {{" +
		"		debug123.resoult123 = e.message;" +
		"	}}" +
		"		return debug123;" +
		"	}};"
}

func getDom() string {
	return "debug123 = {{code123: []}}; " +
		"	if (!this.window) {{window = {{}};}}" +
		"	if (!this.document) {{document = {{}};}}"
}

func getPatch(workflow string) string {
	return strings.ReplaceAll("debug123.code123.push({workflow}); "+
		"var patchCode = function(workflow) {{ "+
		"	var testVari = /(\\w+)=(\\w+)\\([\\w\\+]+\\);.*?(\\w+)=\"\\w+\";/.exec(workflow); "+
		"	if (testVari && testVari[1] == testVari[2]) {{ "+
		"		{workflow} += testVari[1] + \"[\" + testVari[3] + \"] = function() {{return true;}};\"; "+
		"	}} "+
		"}}; "+
		"patchCode({workflow}); "+
		"var subWorkflow = /(?:\\w+=)?eval\\((\\w+)\\)/.exec({workflow}); "+
		"if (subWorkflow) {{ "+
		"	var subPatch = ( "+
		"		\"debug123.code123.push('sub workflow: ' + subWorkflow);\" + "+
		"			\"patchCode(subWorkflow);\" "+
		"	).replace(/subWorkflow/g, subWorkflow[1]) + subWorkflow[0]; "+
		"	{workflow} = {workflow}.replace(subWorkflow[0], subPatch); "+
		"}} "+
		"eval({workflow}); ", "{workflow}", workflow)
}
