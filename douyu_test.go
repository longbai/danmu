package danmu

import (
	"fmt"
	"testing"
)

func TestDouyuLiveStatus(t *testing.T) {
	d := Douyu{271934, nil}
	b, err := d.GetLiveStatus()
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	fmt.Println(b)
}
