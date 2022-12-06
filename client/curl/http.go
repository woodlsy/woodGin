package curl

import (
	"github.com/woodlsy/woodGin/log"
	"io/ioutil"
	"net/http"
	"strings"
)

type Request struct {
	Response *http.Response
	Body     []byte
	Request  *http.Request
	Url      string
	Header   http.Header
}

func Instance() Request {
	return Request{Header: make(http.Header)}
}

func (r *Request) Get(url string) string {
	return r.FetchString(url, "GET")
}

func (r *Request) Post(url string) string {
	return r.FetchString(url, "POST")
}

func (r *Request) SetHeader(key string, value string) *Request {
	r.Header.Add(key, value)
	return r
}

func (r *Request) FetchString(url string, method string) string {
	_, err := r.NewRequest(url, method)
	if err == nil {
		r.Fetch()
	}
	return string(r.Body)
}

//
// Fetch
// @Description: 执行请求
// @receiver r
// @return *Request
// @return error
//
func (r *Request) Fetch() (*Request, error) {
	var err error
	var client http.Client

	r.Response, err = client.Do(r.Request)
	if err != nil {
		log.Logger.Error("url")
		log.Logger.Error("get failed, err:", err)
		return r, err
	}
	defer r.Response.Body.Close()
	r.Body, err = ioutil.ReadAll(r.Response.Body)
	if err != nil {
		log.Logger.Error(r.Url)
		log.Logger.Error("读取请求body失败, err:", err)
		return r, err
	}
	return r, nil
}

//
// NewRequest
// @Description: 配置请求
// @receiver r
// @param url
// @param method
// @return *Request
// @return error
//
func (r *Request) NewRequest(url string, method string) (*Request, error) {
	r.Url = url
	var err error
	method = strings.ToUpper(method)
	r.Request, err = http.NewRequest(method, url, nil)
	if err != nil {
		log.Logger.Error("创建curl请求失败", url, err)
	}
	r.Request.Header = r.Header
	return r, nil
}
