package main

import (
	"encoding/json"
	"fmt"
	"github.com/longbai/danmu"
	// "time"
)

func main() {
	p, _ := danmu.NewDouyu("https://www.douyu.com/271934")
	b, err := p.IsLive()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(b)
	ch, err := p.StartReceive()
	fmt.Println(err)

	for x := range ch {
		// fmt.Print(*x)
		msg, err := danmu.DouyuMessageParse(*x)
		if err != nil {
			fmt.Println(err)
		} else {
			// fmt.Println("msg", msg.Type())
			if msg.Type() == danmu.DouyuText {
				msgText := msg.(*danmu.DouyuTextMessage)
				x, _ := json.Marshal(msgText)
				fmt.Println(string(x))
			} else if msg.Type() == danmu.DouyuGift {
				// msgGift := msg.(*danmu.DouyuGiftMessage)
				// fmt.Println(msgGift.NickName, "送出:", msgGift.GiftId)
			} else if msg.Type() == danmu.DouyuVisitor {
				// msgVisitor := msg.(*danmu.DouyuVisitorMessage)
				// fmt.Println(msgVisitor.NickName, "enter")
			}
		}
	}
}
