package config

type Aliyun struct {
	Sms AliyunSms `json:"sms"`
	Oss AliyunOss `json:"oss"`
}

type AliyunSms struct {
	AccessKeyId     string `json:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret"`
	SignName        string `json:"signName"`
}

type AliyunOss struct {
	AccessKeyId     string `json:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret"`
	Endpoint        string `json:"endpoint"`
}
