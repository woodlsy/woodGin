package wechatApi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/url"
	"os"

	"github.com/woodlsy/woodGin/client/curl"
	"github.com/woodlsy/woodGin/helper"
)

type Request struct {
	AppId       string
	Secret      string
	Url         string
	CustomQuery map[string]string
	CustomData  map[string]interface{}
}

const domain = "https://api.weixin.qq.com/"

func (r *Request) Get(result interface{}) {
	if r.Url == "" {
		panic("未配置请求地址")
	}

	u := domain + r.Url + "?" + r.getParams()

	request := curl.Instance()
	resp := request.Get(u)
	fmt.Println("url GET:", u)
	fmt.Println("result:", resp)
	if resp == "" {
		fmt.Println("请求", u, "失败")
		return
	}
	err := json.Unmarshal(request.Body, &result)
	if err != nil {
		fmt.Println(helper.Join("", domain, r.Url))
		fmt.Println(err)
		// log.Logger.Error("json 解析天气接口数据失败", err, string(request.Body))
		return
	}
}

func (r *Request) Post(result interface{}) {
	if r.Url == "" {
		panic("未配置请求地址")
	}

	u := domain + r.Url + "?" + r.getParams()

	request := curl.Instance()
	if len(r.CustomData) > 0 {
		request.Data = r.CustomData
	}
	resp := request.Post(u)
	fmt.Println("url POST:", u)
	fmt.Println("请求报文:", helper.JsonEncode(r.CustomData))
	fmt.Println("result:", resp)
	if resp == "" {
		fmt.Println("请求", u, "失败")
		return
	}

	err := json.Unmarshal(request.Body, &result)
	if err != nil {
		fmt.Println(helper.Join("", domain, r.Url))
		fmt.Println(err)
		return
	}
}

func (r *Request) PostLocalFile(result interface{}, fileUrl string, fileName string) {
	if r.Url == "" {
		panic("未配置请求地址")
	}

	u := domain + r.Url + "?" + r.getParams()

	request := curl.Instance()

	// 获取文件类型
	file, err := os.Open(fileUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile(fileName, fileUrl)
	if err != nil {
		log.Fatal(err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		log.Fatal(err)
	}

	err = writer.Close()
	if err != nil {
		log.Fatal(err)
	}

	request.RequestBody = body
	request.Header.Add("Content-Type", writer.FormDataContentType())

	resp := request.Post(u)
	fmt.Println("url POST:", u)
	fmt.Println("请求报文:", helper.JsonEncode(r.CustomData))
	fmt.Println("result:", resp)
	if resp == "" {
		fmt.Println("请求", u, "失败")
		return
	}

	err = json.Unmarshal(request.Body, &result)
	if err != nil {
		fmt.Println(helper.Join("", domain, r.Url))
		fmt.Println(err)
		return
	}
}

func (r *Request) getParams() string {
	params := url.Values{}
	if r.AppId != "" {
		params.Add("appid", r.AppId)
	}
	if r.Secret != "" {
		params.Add("secret", r.Secret)
	}

	if len(r.CustomQuery) > 0 {
		for k, v := range r.CustomQuery {
			params.Add(k, v)
		}
	}
	return params.Encode()
}
