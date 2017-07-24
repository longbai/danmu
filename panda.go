package danmu

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	purl "net/url"
	"strconv"
	"time"
)

type Panda struct {
	room int64
	conn *net.TCPConn
	exit bool
	pool chan *string
}

func NewPanda(url string) (d DanmuClient, err error) {
	d, err = newPanda(url)
	return
}

func newPanda(url string) (*Panda, error) {
	u, err := purl.Parse(url)
	if err != nil || len(u.Path) <= 1 {
		return nil, fmt.Errorf("invalid room url %s", url)
	}
	room := u.Path[1:]
	id, err := strconv.ParseInt(room, 10, 64)
	if err != nil {
		return nil, err
	}
	return &Panda{id, nil, false, make(chan *string, DefaultDanmuPool)}, nil
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

func (p *Panda) IsLive() (bool, error) {
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

var pandaStart = []byte{0x00, 0x06, 0x00, 0x02}
var pandaHeartbeat = []byte{0x00, 0x06, 0x00, 0x00}
var pandaResponse = []byte{0x00, 0x06, 0x00, 0x06} //连接弹幕服务器响应
var pandaReceiveMsg = []byte{0x00, 0x06, 0x00, 0x03}
var pandaHeartbeatResponse = []byte{0x00, 0x06, 0x00, 0x01}

const pandaIgnoreByteLength = 16 //弹幕消息体忽略的字节数

func (p *Panda) handshake(param *PandaChatParam) error {
	data := fmt.Sprintf("u:%d@%s\nk:1\nt:300\nts:%d\nsign:%s\nauthtype:%s", param.Rid, param.Appid, param.Ts, param.Sign, param.AuthType)
	l := len(data)

	msg := make([]byte, 4+2+l)
	copy(msg, pandaStart)
	copy(msg[4:], []byte{byte(l >> 8), byte(l & 0xff)})
	copy(msg[4+2:], []byte(data))
	_, err := p.conn.Write(msg)
	if err != nil {
		return err
	}
	buff := make([]byte, 6)
	n, err := p.conn.Read(buff)
	if err != nil {
		return err
	}
	if n != 6 || !bytes.Equal(buff[:4], pandaResponse) {
		return errors.New("response error")
	}

	length := int((uint(buff[4]) << 8) + uint(buff[5]))
	if length > 255 {
		return errors.New("invalid response length flag")
	}
	buff2 := make([]byte, length)
	p.conn.Read(buff2)
	return nil
}

func (p *Panda) keepAlive() {
	go func() {
		p.conn.Write(pandaHeartbeat)
		time.Sleep(KeepAliveInterval)
	}()
}

var typeStart = []byte(`{"type"`)

func (p *Panda) dealBuffer(buff []byte) int {
	if len(buff) <= 4 {
		return 0 // no deal
	}

	if bytes.Equal(buff[:4], pandaReceiveMsg) {
		if len(buff) < 4+2 { // msg length + body not enough, wait next
			return 0
		}
		length := uint(buff[4]<<8) + uint(buff[5])
		pos := int(4 + 2 + length)
		if len(buff) < pos+4+pandaIgnoreByteLength {
			return 0
		}

		msgLen := int((uint(buff[pos]) << 24) + (uint(buff[pos+1]) << 16) + (uint(buff[pos+2]) << 8) + uint(buff[pos+3]))
		if len(buff) < pos+4+msgLen {
			return 0
		}
		pos += 4 + pandaIgnoreByteLength
		strBytes := buff[pos : pos+msgLen-pandaIgnoreByteLength]

		// 弹幕有时有bug，多条消息并在一起，需要拆开
		var n = 0
		for {
			n = bytes.LastIndex(strBytes, typeStart)
			if n == -1 {
				fmt.Println("invalid string", string(strBytes))
				break
			}
			str := string(strBytes[n:])
			p.pool <- &str
			if n == 0 {
				break
			}
			strBytes = strBytes[:n-pandaIgnoreByteLength]
		}

		return pos + msgLen - pandaIgnoreByteLength

	} else if bytes.Equal(buff[:4], pandaHeartbeatResponse) {
		// fmt.Println("heartbeat")
		return 4
	}
	fmt.Println(hex.EncodeToString(buff[:4]))
	return 4
}

func (p *Panda) startReceive(live bool) {
	var err error
	if !live {
		for !p.exit {
			time.Sleep(LiveCheckInterval)
			live, err = p.IsLive()
			if err != nil {
				return
			}
		}
	}

	err = p.init()
	if err != nil {
		return
	}
	p.keepAlive()
	b := make([]byte, 512*1024)
	start := 0
	end := 0
	lastOffset := 0
	for !p.exit {
		start = 0
		end = 0
		n, err := p.conn.Read(b[lastOffset:])
		if err != nil {
			return
		}
		end = lastOffset + n

		for {
			dealed := p.dealBuffer(b[start:end])
			start += dealed
			if dealed == 0 {
				copied := copy(b, b[start:end])
				lastOffset = copied
				break
			}
		}
	}
}

func (p *Panda) StartReceive() (<-chan *string, error) {
	b, err := p.IsLive()
	if err != nil {
		return nil, err
	}

	go p.startReceive(b)
	return p.pool, nil
}

func (p *Panda) Close() {
	p.exit = true
	p.conn.Close()
}
