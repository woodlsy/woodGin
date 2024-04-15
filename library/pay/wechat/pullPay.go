package wechat

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/woodlsy/woodGin/helper"
	"strconv"
	"time"
)

type PullPay struct {
	AppId          string
	MchId          string
	TimeStamp      string
	NonceStr       string
	PrepayId       string
	Package        string
	SignType       string
	SslCertNumber  string
	PrivateKeyPath string
	Utils          Utils
}

func CreatePullPay(appId, mchId, prepayId, sslCertNumber, privateKeyPath string) *PullPay {
	return &PullPay{
		AppId:          appId,
		MchId:          mchId,
		TimeStamp:      strconv.FormatInt(time.Now().Unix(), 10),
		NonceStr:       helper.GenerateRandomString(32),
		PrepayId:       prepayId,
		SslCertNumber:  sslCertNumber,
		PrivateKeyPath: privateKeyPath,
		Package:        "prepay_id=" + prepayId,
		SignType:       "RSA",
	}
}

func (p PullPay) ResultMap() (map[string]interface{}, error) {
	if p.AppId == "" || p.MchId == "" || p.SslCertNumber == "" || p.PrivateKeyPath == "" || p.PrepayId == "" {
		return nil, errors.New("缺少必要参数")
	}
	message := fmt.Sprintf("%s\n%s\n%s\n%s\n", p.AppId, p.TimeStamp, p.NonceStr, p.Package)
	privateKey, err := p.Utils.GetPrivateKeyByPath(p.PrivateKeyPath)
	if err != nil {
		return nil, err
	}
	sign, err := p.Utils.Sign(p.SslCertNumber, privateKey, message)
	if err != nil {
		return nil, err
	}
	data := map[string]interface{}{
		"timeStamp": p.TimeStamp,
		"nonceStr":  p.NonceStr,
		"package":   p.Package,
		"signType":  p.SignType,
		"sign":      sign,
	}
	return data, nil
}
