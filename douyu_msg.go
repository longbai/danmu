package danmu

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	DouyuText            = 1
	DouyuGift            = 2
	DouyuLogin           = 3
	DouyuKeepLive        = 4
	DouyuVisitor         = 5
	DouyuSynExp          = 6
	DouyuRankList        = 7
	DouyuError           = 8
	DouyuQausrespond     = 9
	DouyuOnlineNobleList = 10
	DouyuSpbc            = 11
	DouyuBlab            = 12
)

// type@=chatmsg/rid@=271934/ct@=1/uid@=36761277/nn@=石s乐l志z/txt@=春天我们在一起吧/cid@=d39fac1d2fa64108a4eb100000000000/ic@=avanew@Sface@S201705@S19@S15@Sbba4deb95fd9128c0aa672370f4bb12b/level@=18/sahf@=0/col@=2/bnn@=久哥哥/bl@=8/brid@=271934/hc@=1f9644816bc8fdc6a4b50dd3b2d214c1/ifs@=1/el@=eid@AA=1500000113@ASetp@AA=1@ASsc@AA=1@ASef@AA=0@AS@S/

type DouyuTextMessage struct {
	RoomId   int64  `json:"rid"`
	Ct       int64  `json:"ct"`
	Uid      int64  `json:"uid"`
	NickName string `json:"nn"`
	Content  string `json:"txt"`
	Cid      string `json:"cid"`
	Ic       string `json:"ic"`
	Level    int64  `json:"level"`
	Sahf     int64  `json:"sahf"`
	Col      int64  `json:"col"`
	Bnn      string `json:"bnn"`
	Bl       int64  `json:"bl"`
	Brid     int64  `json:"brid"`
	Hc       string `json:"hc"`
	Ifs      int64  `json:"ifs"`
	El       string `json:"el"`
}

func (d *DouyuTextMessage) Type() int {
	return DouyuText
}

// type@=dgb/rid@=271934/gfid@=519/gs@=2/uid@=58200943/nn@=丿顺其自然丶了/ic@=avatar@Sface@S201607@S21@S0ea652c6e8855bd701eb0f3adb40678e/eid@=0/level@=8/dw@=187375540/hits@=3/ct@=0/cm@=0/bnn@=久哥哥/bl@=6/brid@=271934/hc@=1f9644816bc8fdc6a4b50dd3b2d214c1/sahf@=0/fc@=0/
type DouyuGiftMessage struct {
	RoomId   int64  `json:"rid"`
	GiftId   int64  `json:"gfid"`
	Gs       int64  `json:"gs"`
	Uid      int64  `json:"uid"`
	NickName string `json:"nn"`
	Ic       string `json:"ic"`
	Eid      int64  `json:"eid"`
	Level    int64  `json:"level"`
	Dw       int64  `json:"dw"`
	Hits     int64  `json:"hits"`
	Ct       int64  `json:"ct"`
	Cm       int64  `json:"cm"`
	Bnn      string `json:"bnn"`
	Bl       int64  `json:"bl"`
	Brid     int64  `json:"brid"`
	Hc       string `json:"hc"`
	Sahf     int64  `json:"sahf"`
	Fc       int64  `json:"fc"`
}

func (d *DouyuGiftMessage) Type() int {
	return DouyuGift
}

type DouyuLoginMessage struct {
}

func (d *DouyuLoginMessage) Type() int {
	return DouyuLogin
}

type DouyuKeepLiveMessage struct {
}

func (d *DouyuKeepLiveMessage) Type() int {
	return DouyuKeepLive
}

type DouyuVisitorMessage struct {
	NickName string `json:"nn"`
}

func (d *DouyuVisitorMessage) Type() int {
	return DouyuVisitor
}

type DouyuSyncExpMessage struct {
}

func (d *DouyuSyncExpMessage) Type() int {
	return DouyuSynExp
}

type DouyuRankListMessage struct {
}

func (d *DouyuRankListMessage) Type() int {
	return DouyuRankList
}

type DouyuErrorMessage struct {
	Code int `json:"code"`
}

func (d *DouyuErrorMessage) Type() int {
	return DouyuError
}

type DouyuQausrespondMessage struct {
}

func (d *DouyuQausrespondMessage) Type() int {
	return DouyuQausrespond
}

type DouyuOnlineNobleListMessage struct {
}

func (d *DouyuOnlineNobleListMessage) Type() int {
	return DouyuOnlineNobleList
}

type DouyuSpbcMessage struct {
}

func (d *DouyuSpbcMessage) Type() int {
	return DouyuSpbc
}

type DouyuBlabMessage struct {
}

func (d *DouyuBlabMessage) Type() int {
	return DouyuBlab
}

func douyuSplit(str string) [][]string {
	x := strings.Split(str, "/")
	var ret [][]string
	for _, z := range x {
		z = strings.Replace(z, "@S", "/", -1)
		temp := strings.Split(z, "@=")
		if len(temp) < 2 {
			continue
		}
		for i := 0; i < len(temp); i++ {
			temp[i] = strings.Replace(temp[i], "@A", "@", -1)
		}
		ret = append(ret, temp)
	}
	return ret
}

func douyuTextMessageBuild(str string) *DouyuTextMessage {
	var msg DouyuTextMessage
	kvs := douyuSplit(str)
	for _, it := range kvs {
		if it[0] == "nn" {
			msg.NickName = it[1]
		} else if it[0] == "txt" {
			msg.Content = it[1]
		} else if it[0] == "rid" {
			msg.RoomId, _ = strconv.ParseInt(it[1], 10, 64)
		} else if it[0] == "uid" {
			msg.Uid, _ = strconv.ParseInt(it[1], 10, 64)
		} else if it[0] == "level" {
			msg.Level, _ = strconv.ParseInt(it[1], 10, 64)
		}
	}
	return &msg
}

func douyuGiftMessageBuild(str string) *DouyuGiftMessage {
	var msg DouyuGiftMessage
	kvs := douyuSplit(str)
	for _, it := range kvs {
		if it[0] == "nn" {
			msg.NickName = it[1]
		} else if it[0] == "gfid" {
			msg.GiftId, _ = strconv.ParseInt(it[1], 10, 64)
		}
	}
	return &msg
}

func douyuVisitorMessageBuild(str string) *DouyuVisitorMessage {
	var msg DouyuVisitorMessage
	kvs := douyuSplit(str)
	for _, it := range kvs {
		if it[0] == "nn" {
			msg.NickName = it[1]
		}
	}
	return &msg
}

func DouyuMessageParse(str string) (Message, error) {
	if strings.HasPrefix(str, `type@=chatmsg`) {
		return douyuTextMessageBuild(str), nil
	} else if strings.HasPrefix(str, `type@=dgb`) {
		return douyuGiftMessageBuild(str), nil
	} else if strings.HasPrefix(str, `type@=loginres`) {
		var msg DouyuLoginMessage
		return &msg, nil
	} else if strings.HasPrefix(str, `type@=keeplive`) {
		var msg DouyuKeepLiveMessage
		return &msg, nil
	} else if strings.HasPrefix(str, `type@=uenter`) {
		return douyuVisitorMessageBuild(str), nil
	} else if strings.HasPrefix(str, `type@=synexp`) {
		var msg DouyuSyncExpMessage
		return &msg, nil
	} else if strings.HasPrefix(str, `type@=ranklist`) {
		var msg DouyuRankListMessage
		return &msg, nil
	} else if strings.HasPrefix(str, `type@=error`) {
		var msg DouyuErrorMessage
		return &msg, nil
	} else if strings.HasPrefix(str, `rid@=`) { //  rid@=2500395/sc@=888300/sctn@=100/rid@=-1/type@=qausrespond/
		var msg DouyuQausrespondMessage
		return &msg, nil
	} else if strings.HasPrefix(str, `type@=online_noble_list`) {
		var msg DouyuOnlineNobleListMessage
		return &msg, nil
	} else if strings.HasPrefix(str, `type@=spbc`) {
		var msg DouyuSpbcMessage
		return &msg, nil
	} else if strings.HasPrefix(str, `type@=blab`) {
		var msg DouyuBlabMessage
		return &msg, nil
	} else {
		return nil, fmt.Errorf("unkown msg %s", str)
	}
}
