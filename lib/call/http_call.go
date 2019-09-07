package call

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var (
	client *http.Client
)

func init() {

	client = &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				conn, err := net.DialTimeout(netw, addr, time.Second*5)
				if err != nil {
					return nil, err
				}
				conn.SetDeadline(time.Now().Add(time.Second * 5))
				return conn, nil
			},
			ResponseHeaderTimeout: time.Second * 5,
		},
	}
}

func HttpGetRequestUrl(urlStr string) (string, error) {
	resp, errGet := client.Get(urlStr)
	if errGet != nil {
		return "", errGet
	}
	defer resp.Body.Close()
	return resp.Request.URL.String(), nil
}

func HttpGet(urlStr string, queryList map[string]string) (response []byte, err error) {
	var temp = make([]string, 0, len(queryList))

	for key, value := range queryList {
		stringQuery := key + "=" + url.QueryEscape(value)
		temp = append(temp, stringQuery)
	}

	queryVar := strings.Join(temp, "&")
	fullQuery := urlStr + "?" + queryVar
	resp, errGet := client.Get(fullQuery)
	if errGet != nil {
		return response, errGet
		// handle error
	}

	defer resp.Body.Close()
	body, errRead := ioutil.ReadAll(resp.Body)

	if errRead != nil {
		return response, errRead
		// handle error
	}
	return body, nil
}

func HttpPost(urlStr string, formList map[string]string) (response []byte, err error) {

	var temp []string
	for key, value := range formList {
		temp = append(temp, key+"="+url.QueryEscape(value))
	}
	//fmt.Println(urlStr)
	implodedStr := strings.Join(temp, "&")
	//fmt.Println(implodedStr)
	resp, err := client.Post(urlStr,
		"application/x-www-form-urlencoded",
		strings.NewReader(implodedStr))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, errPost := ioutil.ReadAll(resp.Body)
	if errPost != nil {
		return response, errPost
	}
	return body, nil
}

func HttpPostJson(urlStr string, formList map[string]interface{}) (response []byte, err error) {
	jsonBytes, errJson := json.Marshal(formList)
	if errJson != nil {
		return nil, errJson
	}
	fmt.Println(urlStr)
	fmt.Println(string(jsonBytes))
	resp, err := client.Post(urlStr,
		"application/json",
		bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, errPost := ioutil.ReadAll(resp.Body)
	if errPost != nil {
		return response, errPost
	}
	return body, nil
}

// POST请求 -- 使用http.PostForm()方法
func HttpPostForm(urlStr string) {
	resp, err := client.PostForm(urlStr,
		url.Values{"key": {"Value"}, "id": {"123"}})

	if err != nil {
		// handle error
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	fmt.Println(string(body))
}

func HttpPostFile(urlStr string, path string) (response []byte, err error) {

	bodyBuf := bytes.NewBufferString("")
	bodyWriter := multipart.NewWriter(bodyBuf)

	_, err = bodyWriter.CreateFormFile("file", path)
	if err != nil {
		return nil, err
	}

	fh, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	boundary := bodyWriter.Boundary()

	closeBuf := bytes.NewBufferString(fmt.Sprintf("\r\n--%s--\r\n", boundary))

	requestReader := io.MultiReader(bodyBuf, fh, closeBuf)
	fi, err := fh.Stat()
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", urlStr, requestReader)
	if err != nil {
		return nil, err
	}

	// Set headers for multipart, and Content Length
	req.Header.Add("Content-Type", "multipart/form-data; boundary="+boundary)
	req.ContentLength = fi.Size() + int64(bodyBuf.Len()) + int64(closeBuf.Len())

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err

	}
	return body, nil
}

func HttpGetWithHeader(addr string, queryList map[string]string, headers map[string]string) (response []byte, err error) {
	var temp = make([]string, 0, len(queryList))

	for key, value := range queryList {
		stringQuery := key + "=" + url.QueryEscape(value)
		temp = append(temp, stringQuery)
	}

	var queryVar = strings.Join(temp, "&")
	fullQuery := addr + "?" + queryVar

	request, reqErr := http.NewRequest("GET", fullQuery, nil)
	if reqErr != nil {
		return nil, reqErr
	}
	for key, value := range headers {
		request.Header.Add(key, value)
	}
	resp, doErr := client.Do(request)
	defer resp.Body.Close()
	if doErr != nil {
		return nil, doErr
	}
	response, readErr := ioutil.ReadAll(resp.Body)
	return response, readErr
}

func HttpGetFake(urlStr string, queryList map[string]string) (response []byte, err error) {
	mapHeader := make(map[string]string, 0)
	mapHeader["Accept"] = "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"
	mapHeader["Accept-Charset"] = "UTF-8,*;q=0.5"
	mapHeader["Accept-Language"] = "en-US,en;q=0.8"
	mapHeader["User-Agent"] = "Mozilla/5.0 (X11; Linux x86_64; rv:13.0) Gecko/20100101 Firefox/13.0"

	return HttpGetWithHeader(urlStr, queryList, mapHeader)
}
