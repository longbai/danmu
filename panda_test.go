package danmu

import (
	"fmt"
	"testing"
)

func TestPandaLiveStatus(t *testing.T) {
	p := Panda{66666, nil}
	b, err := p.GetLiveStatus()
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	fmt.Println(b)
}

func TestPandaNoRoomLiveStatus(t *testing.T) {
	p := Panda{6666677777777, nil}
	_, err := p.GetLiveStatus()
	if err == nil {
		t.Fail()
	}
}

func TestPandaChatParam(t *testing.T) {
	p := Panda{66666, nil}
	param, err := p.getChatParam()
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	fmt.Println(param)
}
