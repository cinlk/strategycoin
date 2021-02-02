package requestBuilder

import (
	"BitCoinProfitStrategy/huobi/getRequest"
	"BitCoinProfitStrategy/huobi/signer"
	"fmt"
	"net/url"
	"time"
)

type PrivateUrlBuilder struct {
	host    string
	akKey   string
	akValue string
	smKey   string
	smValue string
	svKey   string
	svValue string
	tKey    string

	signer *signer.Signer
}

func (p *PrivateUrlBuilder) Init(accessKey string, secretKey string, host string) *PrivateUrlBuilder {
	p.akKey = "AccessKeyId"
	p.akValue = accessKey
	p.smKey = "SignatureMethod"
	p.smValue = "HmacSHA256"
	p.svKey = "SignatureVersion"
	p.svValue = "2"
	p.tKey = "Timestamp"

	p.host = host
	p.signer = new(signer.Signer).Init(secretKey)

	return p
}

func (p *PrivateUrlBuilder) Build(method string, path string, request *getRequest.GetRequest) string {
	t := time.Now().UTC()
	return p.BuildWithTime(method, path, t, request)
}

func (p *PrivateUrlBuilder) BuildWithTime(method string, path string, utcDate time.Time, request *getRequest.GetRequest) string {
	t := utcDate.Format("2006-01-02T15:04:05")

	req := new(getRequest.GetRequest).InitFrom(request)
	req.AddParam(p.akKey, p.akValue)
	req.AddParam(p.smKey, p.smValue)
	req.AddParam(p.svKey, p.svValue)
	req.AddParam(p.tKey, t)

	parameters := req.BuildParams()

	signature := p.signer.Sign(method, p.host, path, parameters)

	u := fmt.Sprintf("https://%s%s?%s&Signature=%s", p.host, path, parameters, url.QueryEscape(signature))

	return u
}
