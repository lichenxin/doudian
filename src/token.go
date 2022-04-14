package douDianSdk

import (
	"encoding/json"
	"time"
)

type (
	TokenCreateRequest struct {
		Code      string `json:"code"`
		GrantType string `json:"grant_type"`
		TestShop  string `json:"test_shop"`
		ShopId    string `json:"shop_id"`
	}

	TokenCreateResponse struct {
		AccessToken  string `json:"access_token"`
		AuthorityID  string `json:"authority_id"`
		ExpiresIn    string `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
		Scope        string `json:"scope"`
		ShopID       string `json:"shop_id"`
		ShopName     string `json:"shop_name"`
	}

	TokenRefreshRequest struct {
		RefreshToken string `json:"refresh_token"`
		GrantType    string `json:"grant_type"`
	}

	TokenRefreshResponse struct {
		AccessToken  string `json:"access_token"`
		AuthorityID  string `json:"authority_id"`
		ExpiresIn    string `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
		Scope        string `json:"scope"`
		ShopID       string `json:"shop_id"`
		ShopName     string `json:"shop_name"`
	}
)

// TokenCreate 创建token
func (c *DouDianClient) TokenCreate(code, grantType, testShop string) (*TokenCreateResponse, error) {

	if !(grantType == "authorization_code") || !(grantType == "authorization_self") {
		return nil, GrantTypeIllegalError
	}
	callMethod := "token.create"
	timestamp := time.Now().Unix()
	reqParams := &TokenCreateRequest{
		Code:      code,
		GrantType: grantType,
		TestShop:  testShop,
		ShopId:    c.ShopId,
	}
	paramsJson, sign := c.buildJsonAndSign(reqParams, callMethod, timestamp)
	respBodyData, _ := c.callApi(callMethod, sign, paramsJson, "get", timestamp)
	var ret TokenCreateResponse
	_ = json.Unmarshal(respBodyData, &ret)
	return &ret, nil
}

// TokenRefresh 刷新Token
func (c *DouDianClient) TokenRefresh(refreshToken string) (*TokenCreateResponse, error) {
	callMethod := "token.refresh"
	timestamp := time.Now().Unix()
	reqParams := &TokenRefreshRequest{
		RefreshToken: refreshToken,
		GrantType:    "refresh_token",
	}
	paramsJson, sign := c.buildJsonAndSign(reqParams, callMethod, timestamp)
	respBodyData, _ := c.callApi(callMethod, sign, paramsJson, "get", timestamp)
	var ret TokenCreateResponse
	_ = json.Unmarshal(respBodyData, &ret)
	return &ret, nil
}
