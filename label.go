package ejet

import (
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// 获取面单

type labelService service

type LabelRequest struct {
	OrderCode string `json:"order_code"` // 订单号
}

func (m LabelRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.OrderCode, validation.Required.Error("订单号不能为空")),
	)
}

type Label struct {
	TrackingNumber string `json:"tracking_number"` // 跟踪号
	LabelURL       string `json:"label_url"`       // 面单链接
}

type LabelResponse struct {
	ReferenceNo      string  `json:"reference_no"`       // 参考号
	OrderCode        string  `json:"order_code"`         // 订单号
	OrderAddressType string  `json:"order_address_type"` // 地址类型
	Labels           []Label `json:"labels"`             // 面单信息
}

// GetLabel 获取面单
func (s labelService) GetLabel(req LabelRequest) (labelResponse LabelResponse, err error) {
	if err = req.Validate(); err != nil {
		return
	}

	res := struct {
		Code    int           `json:"code"`
		Result  LabelResponse `json:"result"`
		Message string        `json:"msg"`
	}{}

	resp, err := s.httpClient.R().
		SetBody(req).
		Post("/getLabel")
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

	labelResponse = res.Result

	return
}
