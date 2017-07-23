package danmu

import (
	"fmt"
	"net"
	purl "net/url"
	"strconv"
	"time"
)

type Panda struct {
	room int64
	conn *net.TCPConn
}

func NewPanda(url string) *Panda {
	u, err := purl.Parse(url)
	if err != nil || len(u.Path) <= 1 {
		return nil
	}
	fmt.Println(u.Path)
	room := u.Path[1:]
	id, err := strconv.ParseInt(room, 10, 64)
	if err != nil {
		return nil
	}
	return &Panda{id, nil}
}

type PandaVideoInfo struct {
	Status string `json:"status"`
}

type PandaLiveData struct {
	VideoInfo PandaVideoInfo `json:"videoinfo"`
}

type PandaLive struct {
	Data PandaLiveData `json:"data"`
}

type PandaChatParam struct {
	Rid          int64    `json:"rid"`
	Appid        string   `json:"appid"`
	Ts           int64    `json:"ts"`
	Sign         string   `json:"sign"`
	AuthType     string   `json:"authType"`
	ChatAddrList []string `json:"chat_addr_list"`
}

type pandaChatData struct {
	Data PandaChatParam `json:"data"`
}

func (p *Panda) getChatParam() (*PandaChatParam, error) {
	u1 := fmt.Sprintf("http://www.panda.tv/ajax_chatinfo?roomid=%d&_=%d", p.room, time.Now().Unix()*1000)
	var chatData pandaChatData
	err := GetJson(u1, &chatData)
	if err != nil {
		return nil, err
	}
	u2 := fmt.Sprintf("http://api.homer.panda.tv/chatroom/getinfo?rid=%d&roomid=%d&retry=0&sign=%s&ts=%d&_=%d",
		chatData.Data.Rid, p.room, chatData.Data.Sign, chatData.Data.Ts, time.Now().Unix()*1000)
	var chatData2 pandaChatData
	err = GetJson(u2, &chatData2)
	if err != nil {
		return nil, err
	}
	return &chatData2.Data, nil
}

func (p *Panda) GetLiveStatus() (bool, error) {
	url := fmt.Sprintf("http://www.panda.tv/api_room?roomid=%d&pub_key=&_=%d", p.room, time.Now().Unix())
	var pandaLive PandaLive
	err := GetJson(url, &pandaLive)
	if err != nil {
		return false, err
	}
	return pandaLive.Data.VideoInfo.Status == "2", nil
}

func (p *Panda) init() error {
	param, err := p.getChatParam()
	if err != nil {
		return err
	}

	conn, err := net.Dial("tcp", param.ChatAddrList[0])
	if err != nil {
		fmt.Println(err)
		return err
	}
	p.conn = conn.(*net.TCPConn)
	p.handshake(param)
	return nil
}

var start = []byte{0x00, 0x06, 0x00, 0x02}

func (p *Panda) handshake(param *PandaChatParam) error {
	data := fmt.Sprintf("u:%d@%s\nk:1\nt:300\nts:%d\nsign:%s\nauthtype:%s", param.Rid, param.Appid, param.Ts, param.Sign, param.AuthType)
	l := len(data)

	msg := make([]byte, 4+2+l)
	copy(msg, []byte{0x00, 0x06, 0x00, 0x02})
	copy(msg[4:], []byte{byte(l >> 8), byte(l & 0xff)})
	copy(msg[4+2:], []byte(data))
	_, err := p.conn.Write(msg)
	return err
}

var heatbeat = []byte{0x00, 0x06, 0x00, 0x00}

func (p *Panda) keepAlive() {
	go func() {
		p.conn.Write(heatbeat)
		time.Sleep(time.Minute)
	}()
}

func (p *Panda) StartReceive() error {
	err := p.init()
	if err != nil {
		return err
	}
	p.keepAlive()
	b := make([]byte, 512*1024)
	for {
		n, err := p.conn.Read(b)
		if err != nil {
			return err
		}
		fmt.Println(n, string(b[:n]))
	}
	return nil
}

func (p *Panda) Close() {
	p.conn.Close()
}
