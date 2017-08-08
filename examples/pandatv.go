package main

import (
	"encoding/json"
	"fmt"

	"github.com/longbai/danmu"
	// "time"
)

func main() {
	p, _ := danmu.NewPanda("https://www.panda.tv/135069")
	b, err := p.IsLive()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(b)
	ch, err := p.StartReceive()
	fmt.Println(err)

	for x := range ch {
		msg, err := danmu.PandaMessageParse(*x)
		if err != nil {
			// fmt.Println(err)
		} else {
			if msg.Type() == danmu.PandaText {
				msgText := msg.(*danmu.PandaTextMessage)
				// fmt.Printf("\"%s\",\"%s\"\n", msgText.Data.User.NickName, msgText.Data.Content)
				x, _ := json.Marshal(msgText)
				fmt.Println(string(x))
			}
		}
	}
}
