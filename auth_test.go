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
	var c config.Config
	err = json.Unmarshal(b, &c)
	if err != nil {
		panic(fmt.Sprintf("Parse config file error: %s", err.Error()))
	}
	ejetClient = NewEJet(c)
	m.Run()
}

func Test_authService_GetToken(t *testing.T) {
	req := AuthRequest{
		AppToken: ejetClient.config.AppToken,
		AppKey:   ejetClient.config.AppKey,
	}
	resp, err := ejetClient.Services.Auth.GetToken(req)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(resp.AccessToken)
}
