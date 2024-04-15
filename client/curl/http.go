package curl

import (
	"bytes"
	"fmt"
	"github.com/woodlsy/woodGin/helper"
	"github.com/woodlsy/woodGin/log"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"strings"
)

type Request struct {
	Response    *http.Response
	Body        []byte
	Request     *http.Request
	Url         string
	Header      http.Header
	RequestBody *bytes.Buffer
	Data        map[string]interface{}
}

func Instance() Request {
	return Request{Header: make(http.Header)}
}

func (r *Request) Get(url string) string {
	return r.FetchString(url, "GET")
}

func (r *Request) Post(url string) string {
	if len(r.Data) > 0 {
		r.RequestBody  = bytes.NewBuffer([]byte(helper.JsonEncode(r.Data)))
	}
	return r.FetchString(url, "POST")
}

func (r *Request) PostLocalFile(url string, filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		log.Logger.Error("待上传文件打不开", filePath)
		return ""
	}
	defer file.Close()
	body := new(bytes.Buffer)

	writer := multipart.NewWriter(body)

	//for key, val := range r.Data {
	//	_ = writer.WriteField(key, val)
	//}
	formFile, err := r.createFormFile(writer, "file", filePath)
	if err != nil {
		log.Logger.Error("CreateFormFile err: %v, file: %s", err, file)
		return ""
	}
	_, err = io.Copy(formFile, file)
	if err != nil {
		return ""
	}
	writer.Close()
	helper.P(writer.FormDataContentType(), " writer.FormDataContentType()")
	r.Header.Set("Content-Type", writer.FormDataContentType())
	r.RequestBody = body
	return r.FetchString(url, "POST")
}

//
// createFormFile
// @Description: 重新writer.CreateFormFile
// @receiver r
// @param writer
// @param file
// @param filedName
// @param filePath
// @return io.Writer
// @return error
//
func (r *Request) createFormFile(writer *multipart.Writer, filedName string, filePath string) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			escapeQuotes(filedName), escapeQuotes(filePath)))
	contentType := helper.GetFileMiMeType(filePath)
	h.Set("Content-Type", contentType)
	return writer.CreatePart(h)
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

func (r *Request) SetHeader(key string, value string) *Request {
	r.Header.Add(key, value)
	return r
}

func (r *Request) HeaderJson() *Request {
	r.Header.Set("Content-Type", "application/json")
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
	if r.RequestBody == nil {
		r.Request, err = http.NewRequest(method, url, nil)
	} else {
		r.Request, err = http.NewRequest(method, url, r.RequestBody)
	}

	if err != nil {
		log.Logger.Error("创建curl请求失败", url, err)
	}
	r.Request.Header = r.Header
	return r, nil
}
