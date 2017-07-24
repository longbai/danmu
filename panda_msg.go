package danmu

import (
	"encoding/json"
	"strconv"
	"strings"
)

const (
	PandaText         = 1
	PandaBamboo       = 206
	PandaVisitor      = 207
	PandaVisitorTotal = 205
	PandaGift         = 306
	PandaGiftOther    = 311
)

type PandaUser struct {
	Identity   string `json:"identity"`
	NickName   string `json:"nickName"`
	Badge      string `json:"badge"`
	Rid        string `json:"rid"`
	SpIdentity string `json:"sp_identity"`
	Level      string `json:"level"`
	IsPay      int    `json:"ispay"`
	Platform   string `json:"__plat"`
	UserName   string `json:"userName"`
	MsgColor   string `json:"msgcolor"`
	Hl         string `json:"hl"`
}

type PandaMessageToRoom struct {
	Room string `json"toroom"`
}

type PandaTextMessageData struct {
	User    PandaUser          `json:"from"`
	Room    PandaMessageToRoom `json:"to"`
	Content string             `json:"content"`
}

type PandaTextMessage struct {
	Time int64                `json:"time"`
	Data PandaTextMessageData `json:"data"`
}

// gift content
// "content": {
//       "avatar": "http://i9.pdim.gs/b43b5244c6c13f8b13b515007b8c2de8.jpeg",
//       "combo": "16",
//       "count": "1",
//       "effective": "2",
//       "group": "5975ebadb5061c5de7433be1",
//       "id": "59008f433c74f35f06e6589a",
//       "name": "666",
//       "newBamboos": "7620",
//       "newExp": "114.2",
//       "pic": {
//         "pc": {
//           "chat": "http://i6.pdim.gs/1ba49540dbaaf83dfef22264fc062846.png",
//           "effect": "http://i5.pdim.gs/1a45a13f262ec35df565cd4b7f81007c.gif"
//         }
//       },
//       "position": "2",
//       "price": "1"
//     }

type PandaGiftContent struct {
	Avatar     string `json:"avatar"`
	Combo      string `json:"combo"`
	Count      string `json:"count"`
	Effective  string `json:"effective"`
	Group      string `json:"group"`
	Id         string `json:"id"`
	Name       string `json:"name"`
	NewBamboos string `json:"newBamboos"`
	NewExp     string `json:"newExp"`
	Position   string `json:"position"`
	Price      string `json:"price"`
}

type PandaGiftData struct {
	User    PandaUser          `json:"from"`
	Room    PandaMessageToRoom `json:"to"`
	Content PandaGiftContent   `json:"content"`
}

type PandaGiftMessage struct {
	Time int64         `json:"time"`
	Data PandaGiftData `json:"data"`
}

type PandaCommonData struct {
	User    PandaUser          `json:"from"`
	Room    PandaMessageToRoom `json:"to"`
	Content string             `json:"content"`
}

type PandaBambooMessage struct {
	Time int64           `json:"time"`
	Data PandaCommonData `json:"data"`
}

type PandaVisitorMessage struct {
	Time int64           `json:"time"`
	Data PandaCommonData `json:"data"`
}

type PandaVisitorTotalContent struct {
	ShowNum int64 `json:"show_num"`
	Total   int64 `json:"total"`
}

type PandaVisitorTotalData struct {
	User    PandaUser                `json:"from"`
	Room    PandaMessageToRoom       `json:"to"`
	Content PandaVisitorTotalContent `json:"content"`
}

type PandaVisitorTotalMessage struct {
	Time int64                 `json:"time"`
	Data PandaVisitorTotalData `json:"data"`
}

//other gift
// {
//   "type": "311",
//   "time": 1500895554,
//   "data": {
//     "from": {
//       "identity": "30",
//       "nickName": "三生烟火",
//       "rid": "3039568",
//       "sp_identity": "0"
//     },
//     "to": {
//       "nickName": "阿狸师妹",
//       "rid": "84067672",
//       "roomid": "1022918",
//       "toroom": "0"
//     },
//     "content": {
//       "begintime": "1500895554",
//       "countdown": "150",
//       "eventid": "mall-bsw",
//       "extra": "",
//       "group": "5975d94268b7642beb1c1acb",
//       "platshow": "1",
//       "position": "1,2,3,4",
//       "roomshow": "1",
//       "times": "1",
//       "usercombo": "1"
//     }
//   }
// }

type PandaUnknownMessage struct {
	Kind string `json:"type"`
	Time int64  `json:"time"`
}

type PandaGiftOtherRoom struct {
	NickName string `json:"nickname"`
	Rid      string `json:"rid"`
	RoomId   string `json:"roomid"`
	ToRoom   string `json:"toroom"`
}

type PandaGiftOtherConent struct {
	BeginTime string `json:"begintime"`
	CountDown string `json:"countdown"`
	EventId   string `json:"eventid"`
	Extra     string `json:"extra"`
	Group     string `json:"group"`
	PlatShow  string `json:"platshow"`
	Position  string `json:"position"`
	RoomShow  string `json:"roomshow"`
	Times     string `json:"times"`
	UserCombo string `json:"usercombo"`
}

type PandaGiftOtherData struct {
	User    PandaUser            `json:"from"`
	Room    PandaGiftOtherRoom   `json:"to"`
	Content PandaGiftOtherConent `json:"content"`
}

type PandaGiftOtherMessage struct {
	Time int64              `json:"time"`
	Data PandaGiftOtherData `json:"data"`
}

func (p *PandaTextMessage) Type() int {
	return PandaText
}

func (p *PandaGiftMessage) Type() int {
	return PandaGift
}

func (p *PandaGiftOtherMessage) Type() int {
	return PandaGiftOther
}

func (p *PandaBambooMessage) Type() int {
	return PandaBamboo
}

func (p *PandaVisitorMessage) Type() int {
	return PandaVisitor
}

func (p *PandaVisitorTotalMessage) Type() int {
	return PandaVisitorTotal
}

func (p *PandaUnknownMessage) Type() int {
	x, err := strconv.ParseInt(p.Kind, 10, 64)
	if err != nil {
		return -1
	}
	return int(x)
}

func PandaMessageParse(str string) (Message, error) {
	if strings.HasPrefix(str, `{"type":"1"`) {
		var msg PandaTextMessage
		err := json.Unmarshal([]byte(str), &msg)
		if err != nil {
			return nil, err
		}
		return &msg, nil
	} else if strings.HasPrefix(str, `{"type":"206"`) {
		var msg PandaBambooMessage
		err := json.Unmarshal([]byte(str), &msg)
		if err != nil {
			return nil, err
		}
		return &msg, nil
	} else if strings.HasPrefix(str, `{"type":"205"`) {
		var msg PandaVisitorTotalMessage
		err := json.Unmarshal([]byte(str), &msg)
		if err != nil {
			return nil, err
		}
		return &msg, nil
	} else if strings.HasPrefix(str, `{"type":"207"`) {
		var msg PandaVisitorMessage
		err := json.Unmarshal([]byte(str), &msg)
		if err != nil {
			return nil, err
		}
		return &msg, nil
	} else if strings.HasPrefix(str, `{"type":"306"`) {
		var msg PandaGiftMessage
		err := json.Unmarshal([]byte(str), &msg)
		if err != nil {
			return nil, err
		}
		return &msg, nil
	} else if strings.HasPrefix(str, `{"type":"311"`) {
		var msg PandaGiftOtherMessage
		err := json.Unmarshal([]byte(str), &msg)
		if err != nil {
			return nil, err
		}
		return &msg, nil
	} else {
		var msg PandaUnknownMessage
		err := json.Unmarshal([]byte(str), &msg)
		if err != nil {
			return nil, err
		}
		return &msg, nil
	}
}
