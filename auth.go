package ejet

import (
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// 获取面单

type authService service

type AuthRequest struct {
	AppToken string `json:"app_token"` // App Token
	AppKey   string `json:"app_key"`   // App Key
}

func (m AuthRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.AppToken, validation.Required.Error("App Token 不能为空")),
		validation.Field(&m.AppKey, validation.Required.Error("App Key 不能为空")),
	)
}

type UserInformation struct {
	UID           string `json:"u_id"`            // ID
	UAccount      string `json:"u_account"`       // 帐号
	UCustomerCode string `json:"u_customer_code"` // 客户代码
}

type AuthResponse struct {
	AccessToken string          `json:"access_token"` // 授权 Token
	UserInfo    UserInformation `json:"user_info"`    // 用户信息
}

// GetLabel 获取面单
func (s authService) GetLabel(req AuthRequest) (authResponse AuthResponse, err error) {
	if err = req.Validate(); err != nil {
		return
	}

	res := struct {
		Code    int          `json:"code"`
		Result  AuthResponse `json:"result"`
		Message string       `json:"msg"`
	}{}

	resp, err := s.httpClient.R().
		SetBody(req).
		Post("/getToken")
	if err != nil {
		return
	}

	if err = json.Unmarshal(resp.Body(), &res); err != nil {
		return
	}
	if err != nil {
		return
	}
	err = ErrorWrap(res.Code, res.Message)
	if err != nil {
		return
	}

	return
}
