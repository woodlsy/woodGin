package config

type Aliyun struct {
	Sms AliyunSms              `json:"sms"`
	Oss AliyunOss              `json:"oss"`
	Ocr map[string]interface{} `json:"ocr"`
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
	Bucket          string `json:"bucket"`
	Domain          string `json:"domain"`
	PrivateBucket   string `json:"privateBucket"`
	PrivateDomain   string `json:"privateDomain"`
}
