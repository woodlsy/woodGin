package wechat

import (
	"context"
	"crypto/rsa"
	"errors"
	"github.com/wechatpay-apiv3/wechatpay-go/core/auth/signers"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
	"github.com/woodlsy/woodGin/log"
)

type Utils struct {
	Signer signers.SHA256WithRSASigner
}

func (u Utils) Sign(sslCertNumber string, privateKey *rsa.PrivateKey, message string) (string, error) {
	u.Signer.CertificateSerialNo = sslCertNumber
	u.Signer.PrivateKey = privateKey
	sign, err := u.Signer.Sign(context.Background(), message)
	if err != nil {
		log.Logger.Error("加签失败：", err)
		return "", errors.New("加签失败")
	}
	return sign.Signature, nil
}

func (u Utils) GetPrivateKeyByPath(sslKeyPath string) (*rsa.PrivateKey, error) {
	return utils.LoadPrivateKeyWithPath(sslKeyPath)
}

func (u Utils) GetPrivateKeyByString(sslKey string) (*rsa.PrivateKey, error) {
	return utils.LoadPrivateKey(sslKey)
}

/*
func verifyNotEmpty(params ...interface{}) {
	if len(params) == 0 {
		return nil
	}
	for _, param := range params {
		switch param.(type) {
		case string:
			if param == "" {
				return
			}
		}
	}
}
*/
