package danmu

import (
	// "bytes"
	"encoding/binary"
	// "encoding/hex"
	"fmt"
	"net"
	purl "net/url"
	"strconv"
	"time"
)

type Douyu struct {
	room int64
	conn *net.TCPConn
	exit bool
	pool chan *string
}

func NewDouyu(url string) (*Douyu, error) {
	u, err := purl.Parse(url)
	if err != nil || len(u.Path) <= 1 {
		return nil, err
	}
	room := u.Path[1:]
	id, err := strconv.ParseInt(room, 10, 64)
	if err != nil {
		return nil, err
	}
	return &Douyu{id, nil, false, make(chan *string, DefaultDanmuPool)}, nil
}

type liveResultData struct {
	Status string `json:"room_status"`
	RoomId string `json:"room_id"`
}

type liveResult struct {
	Error int            `json:"error"`
	Data  liveResultData `json:"data"`
}

func (d *Douyu) IsLive() (bool, error) {
	url := fmt.Sprintf("http://open.douyucdn.cn/api/RoomApi/room/%d", d.room)
	var result liveResult
	err := GetJson(url, &result)
	if err != nil {
		return false, err
	}
	// fmt.Println(result.Data.RoomId)
	return result.Data.Status == "1", nil
}

// def _prepare_env(self):
//      return ('openbarrage.douyutv.com', 8601), {'room_id': self.roomId}
//  def _init_socket(self, danmu, roomInfo):
//      self.danmuSocket = _socket()
//      self.danmuSocket.connect(danmu)
//      self.danmuSocket.settimeout(3)
//      self.danmuSocket.communicate('type@=loginreq/roomid@=%s/'%roomInfo['room_id'])
//      self.danmuSocket.push('type@=joingroup/rid@=%s/gid@=-9999/'%roomInfo['room_id'])

// class _socket(socket.socket):
//     def communicate(self, data):
//         self.push(data)
//         return self.pull()
//     def push(self, data):
//         s = pack('i', 9 + len(data)) * 2
//         s += b'\xb1\x02\x00\x00' # 689
//         s += data.encode('ascii') + b'\x00'
//         self.sendall(s)

var douyuFlag = []byte{0xb1, 0x02, 0x00, 0x00}

func (d *Douyu) send(data []byte) error {
	l := len(data)
	buffer := make([]byte, 13+l)
	binary.LittleEndian.PutUint32(buffer, uint32(l)+9)
	binary.LittleEndian.PutUint32(buffer[4:], uint32(l)+9) //dup write
	copy(buffer[8:], douyuFlag)
	copy(buffer[12:], data)
	buffer[l+13-1] = 0x00
	_, err := d.conn.Write(buffer)
	return err
}

func (d *Douyu) handshake() error {
	err := d.send([]byte(fmt.Sprintf("type@=loginreq/roomid@=%d/", d.room)))
	if err != nil {
		return err
	}
	return d.send([]byte(fmt.Sprintf("type@=joingroup/rid@=%d/gid@=-9999/", d.room)))
}

func (d *Douyu) KeepAlive() {
	d.send([]byte(fmt.Sprintf("type@=keeplive/tick@=%d/", time.Now().Unix())))
	time.Sleep(KeepAliveInterval)
}

func (d *Douyu) init() error {
	conn, err := net.Dial("tcp", "openbarrage.douyutv.com:8601")
	if err != nil {
		fmt.Println(err)
		return err
	}
	d.conn = conn.(*net.TCPConn)
	err = d.handshake()
	if err != nil {
		return err
	}
	go d.KeepAlive()
	return nil
}

// def get_danmu(self):
//             if not select.select([self.danmuSocket], [], [], 1)[0]: return
//             content = self.danmuSocket.pull()
//             for msg in re.findall(b'(type@=.*?)\x00', content):
//                 try:
//                     msg = msg.replace(b'@=', b'":"').replace(b'/', b'","')
//                     msg = msg.replace(b'@A', b'@').replace(b'@S', b'/')
//                     msg = json.loads((b'{"' + msg[:-2] + b'}').decode('utf8', 'ignore'))
//                     msg['NickName'] = msg.get('nn', '')
//                     msg['Content']  = msg.get('txt', '')
//                     msg['MsgType']  = {'dgb': 'gift', 'chatmsg': 'danmu',
//                         'uenter': 'enter'}.get(msg['type'], 'other')
//                 except Exception as e:
//                     pass
//                 else:
//                     self.danmuWaitTime = time.time() + self.maxNoDanMuWait
//                     self.msgPipe.append(msg)

func (d *Douyu) dealBuffer(buff []byte) int {
	newBuff := buff
	for {
		l := len(newBuff)
		if l <= 8 {
			break // no deal
		}
		x := binary.LittleEndian.Uint32(newBuff)
		if l < 4+int(x) {
			break
		}
		str := string(newBuff[8+4 : 4+x])
		d.pool <- &str
		newBuff = newBuff[4+x:]
	}
	return len(buff) - len(newBuff)
}

func (d *Douyu) startReceive(live bool) {
	var err error
	if !live {
		for !d.exit {
			time.Sleep(LiveCheckInterval)
			live, err = d.IsLive()
			if err != nil {
				return
			}
		}
	}

	err = d.init()
	if err != nil {
		return
	}
	b := make([]byte, 512*1024)
	start := 0
	end := 0
	lastOffset := 0
	for !d.exit {
		start = 0
		end = 0
		n, err := d.conn.Read(b[lastOffset:])
		if err != nil {
			return
		}
		end = lastOffset + n
		for {
			dealed := d.dealBuffer(b[start:end])
			start += dealed
			if dealed == 0 {
				copied := copy(b, b[start:end])
				lastOffset = copied
				break
			}
		}
	}
}

func (d *Douyu) StartReceive() (<-chan *string, error) {
	b, err := d.IsLive()
	if err != nil {
		return nil, err
	}

	go d.startReceive(b)
	return d.pool, nil
}

func (d *Douyu) Close() {
	d.exit = true
	d.conn.Close()
}
