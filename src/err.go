package douDianSdk

import "errors"

var (
	AppKeyEmptyError      = errors.New("app key is require")
	AppSecretEmptyError   = errors.New("app secret is require")
	GrantTypeIllegalError = errors.New("grant type illegal")
)
