package danmu

import (
	"fmt"
	"testing"
)

func TestPandaLiveStatus(t *testing.T) {
	p, err := newPanda("http://www.pandatv.com/66666")
	fmt.Println("xxxxx", p, err)
	b, err := p.IsLive()
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	fmt.Println(b)
}

func TestPandaNoRoomLiveStatus(t *testing.T) {
	p, _ := newPanda("http://www.pandatv.com/66666777777")
	_, err := p.IsLive()
	if err == nil {
		t.Fail()
	}
}

func TestPandaChatParam(t *testing.T) {
	p, _ := newPanda("http://www.pandatv.com/66666")
	param, err := p.getChatParam()
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	fmt.Println(param)
}
