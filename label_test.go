package ejet

import (
	"fmt"
	"testing"
)

func Test_labelService_GetLabel(t *testing.T) {
	labelResponse, err := ejetClient.Services.Label.GetLabel(LabelRequest{OrderCode: "YJ0105202303317378452"})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(fmt.Sprintf("%#v", labelResponse))

}
