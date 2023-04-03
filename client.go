package ejet

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/hiscaler/e-jet-go/config"
	"strings"
	"time"
)

const (
	Version   = "0.0.1"
	userAgent = "e-jet API Client-Golang/" + Version + " (https://github.com/hiscaler/e-jet-go)"
)

const (
	OK                   = 200 // 无错误
	BadRequestError      = 400 // Bad Request
	UnauthorizedError    = 401 // 授权失败
	ServiceNotFoundError = 404 // 服务不存在
	InternalError        = 500 // 内部错误，数据库异常
)

type EJet struct {
	config      *config.Config // 配置
	logger      Logger         // 日志
	httpClient  *resty.Client  // Resty Client
	accessToken string         // Access Token
	Services    services       // API Services
}

func NewEJet(cfg config.Config) *EJet {
	ejetClient := &EJet{
		config: &cfg,
		logger: createLogger(),
	}
	httpClient := resty.
		New().
		SetDebug(cfg.Debug).
		SetHeaders(map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json",
			"User-Agent":   userAgent,
		})
	httpClient.SetBaseURL("http://e-jet.cn/api/svc")

	httpClient.
		SetTimeout(time.Duration(cfg.Timeout) * time.Second).
		OnBeforeRequest(func(client *resty.Client, request *resty.Request) error {
			if request.URL == "/getToken" {
				return nil
			}

			if ejetClient.accessToken == "" {
				authResponse, e := ejetClient.Services.Auth.GetToken(AuthRequest{
					AppToken: cfg.AppToken,
					AppKey:   cfg.AppKey,
				})
				if e != nil {
					ejetClient.logger.Errorf("Auth error: %s", e.Error())
					return e
				}
				ejetClient.accessToken = authResponse.AccessToken
			}
			request.SetHeaders(map[string]string{
				"Authorization": ejetClient.accessToken,
			})
			return nil
		}).
		OnAfterResponse(func(client *resty.Client, response *resty.Response) (err error) {
			if response.IsError() {
				r := struct {
					Code    int    `json:"code"`
					Message string `json:"msg,omitempty"`
				}{}
				if err = json.Unmarshal(response.Body(), &r); err == nil {
					err = ErrorWrap(r.Code, r.Message)
				}
			}
			return
		}).
		SetRetryCount(2).
		SetRetryWaitTime(2 * time.Second).
		SetRetryMaxWaitTime(10 * time.Second)
	ejetClient.httpClient = httpClient
	xService := service{
		config:     &cfg,
		logger:     ejetClient.logger,
		httpClient: ejetClient.httpClient,
	}
	ejetClient.Services = services{
		Auth:  (authService)(xService),
		Label: (labelService)(xService),
	}
	return ejetClient
}

// SetDebug 设置是否开启调试模式
func (ejet *EJet) SetDebug(v bool) *EJet {
	ejet.config.Debug = v
	ejet.httpClient.SetDebug(v)
	return ejet
}

// SetLogger 设置日志器
func (ejet *EJet) SetLogger(logger Logger) *EJet {
	ejet.logger = logger
	return ejet
}

// ErrorWrap 错误包装
func ErrorWrap(code int, message string) error {
	if code == OK || code == 0 {
		return nil
	}

	switch code {
	case BadRequestError:
		if message == "" {
			message = "Bad Request"
		}
	case UnauthorizedError:
		if message == "" {
			message = "Unauthorized"
		}
	case ServiceNotFoundError:
		if message == "" {
			message = "服务不存在"
		}
	default:
		if code == InternalError {
			if message == "" {
				message = "内部错误，请联系 e-jet"
			}
		} else {
			message = strings.TrimSpace(message)
			if message == "" {
				message = "未知错误"
			}
		}
	}
	return fmt.Errorf("%d: %s", code, message)
}
