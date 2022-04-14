package douDianSdk

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/cast"
	"net/http"
	"strconv"
	"strings"
)

type DouDianClient struct {
	AppKey      string // 应用key
	AppSecret   string // 应用秘钥
	ShopId      string // 店铺ID
	AccessToken string
	//Log 	  *Logger
}

// NewDouDianClient 创建sdk客户端
func NewDouDianClient(appKey, appSecret, shopId string) (*DouDianClient, error) {
	if appKey == "" {
		return nil, AppKeyEmptyError
	}
	if appSecret == "" {
		return nil, AppSecretEmptyError
	}
	return &DouDianClient{
		AppKey:    appKey,
		AppSecret: appSecret,
		ShopId:    shopId,
	}, nil
}

// SetAccessToken 设置AccessToken
func (c *DouDianClient) SetAccessToken(accessToken string) *DouDianClient {
	c.AccessToken = accessToken
	return c
}

// sign 计算签名
func (c *DouDianClient) sign(method, paramJson string, timestamp int64) string {

	paramPattern := "app_key" + c.AppKey + "method" + method + "param_json" + paramJson + "timestamp" + strconv.FormatInt(timestamp, 10) + "v2"
	signPattern := c.AppSecret + paramPattern + c.AppSecret
	return c.signHmac(signPattern, c.AppSecret)
}

// signHmac 计算hmac
func (c *DouDianClient) signHmac(params, appSecret string) string {
	h := hmac.New(sha256.New, []byte(appSecret))
	_, _ = h.Write([]byte(params))
	return hex.EncodeToString(h.Sum(nil))
}

// signMarshal 序列化参数
func (c *DouDianClient) signMarshal(params interface{}) string {

	paramsMap := cast.ToStringMap(params)
	buffer := bytes.NewBufferString("")
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	_ = encoder.Encode(paramsMap)
	marshal := strings.TrimSpace(buffer.String()) // Trim掉末尾的换行符
	return marshal
}

func (c *DouDianClient) buildJsonAndSign(params interface{}, method string, timestamp int64) (string, string) {
	paramsJson := c.signMarshal(params)
	sign := c.sign(method, paramsJson, timestamp)
	return paramsJson, sign
}

// callApi 调用抖音Api
func (c *DouDianClient) callApi(method, sign, paramJson, httpMethod string, timestamp int64) ([]byte, error) {

	methodPath := strings.Replace(method, ".", "/", -1)
	uri := "https://openapi-fxg.jinritemai.com/" + methodPath
	headers := map[string]string{
		"Accept":       "*/*",
		"Content-Type": "application/json;charset=UTF-8",
	}
	params := map[string]string{
		"v":            "2",
		"method":       method,
		"app_key":      c.AppKey,
		"access_token": c.AccessToken,
		"timestamp":    strconv.FormatInt(timestamp, 10),
		"sign":         sign,
		"sign_method":  "hmac-sha256",
		"param_json":   paramJson,
	}
	client := resty.New()
	resp := &resty.Response{}
	var err error
	if httpMethod == "get" || httpMethod == "GET" {
		resp, err = client.R().SetHeaders(headers).SetQueryParams(params).EnableTrace().Get(uri)
	} else {
		body := strings.NewReader(paramJson)
		resp, err = client.R().SetHeaders(headers).SetQueryParams(params).SetBody(body).EnableTrace().Post(uri)
	}
	if err != nil || resp.StatusCode() != http.StatusOK || resp.Body() == nil {
		return resp.Body(), err
	}
	type respBody struct {
		ErrNo   int64       `json:"err_no"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}
	var data respBody
	_ = json.Unmarshal(resp.Body(), &data)
	dataDataByte, _ := json.Marshal(data.Data)
	if data.ErrNo != 0 {
		errMsg := fmt.Sprintf("douDianSdk-->openHttp-->Fetch respBodyErr ErrNo:%s, message:%s", data.ErrNo, data.Message)
		return dataDataByte, errors.New(errMsg)
	}
	return dataDataByte, nil
}
