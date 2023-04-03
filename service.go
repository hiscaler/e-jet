package ejet

import (
	"github.com/go-resty/resty/v2"
	"github.com/hiscaler/e-jet-go/config"
)

type service struct {
	config     *config.Config // Config
	logger     Logger         // Logger
	httpClient *resty.Client  // HTTP client
}

// API Services
type services struct {
	Auth  authService
	Label labelService
}
