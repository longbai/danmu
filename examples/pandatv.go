package main

import (
	"fmt"
	"github.com/longbai/danmu"
	// "time"
)

func main() {
	p, _ := danmu.NewPanda("https://www.panda.tv/66666")
	b, err := p.IsLive()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(b)
	ch, err := p.StartReceive()
	fmt.Println(err)

	for x := range ch {
		fmt.Println(*x)
	}
}
