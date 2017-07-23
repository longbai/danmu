package main

import (
	"fmt"
	"github.com/longbai/danmu"
)

func main() {
	p := danmu.NewPanda("https://www.panda.tv/66666")
	b, err := p.GetLiveStatus()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(b)
	err = p.StartReceive()
	fmt.Println(err)
}
