package ejet

import (
	"encoding/json"
	"fmt"
	"github.com/hiscaler/e-jet-go/config"
	"os"
	"testing"
)

var ejetClient *EJet

func TestMain(m *testing.M) {
	b, err := os.ReadFile("./config/config.json")
	if err != nil {
		panic(fmt.Sprintf("Read config error: %s", err.Error()))
	}
	c := struct {
		Debug    bool
		Sandbox  bool
		AppToken string
		AppKey   string
	}{}
	err = json.Unmarshal(b, &c)
	if err != nil {
		panic(fmt.Sprintf("Parse config file error: %s", err.Error()))
	}

	cfg := config.Config{
		Debug:    c.Debug,
		AppToken: c.AppToken,
		AppKey:   c.AppKey,
	}
	ejetClient = NewEJet(cfg)
	m.Run()
}
