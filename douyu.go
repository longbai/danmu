package danmu

import (
	"fmt"
	"net"
	purl "net/url"
	"strconv"
	// "time"
)

type Douyu struct {
	room int64
	conn *net.TCPConn
}

func NewDouyu(url string) *Douyu {
	u, err := purl.Parse(url)
	if err != nil || len(u.Path) <= 1 {
		return nil
	}
	room := u.Path[1:]
	id, err := strconv.ParseInt(room, 10, 64)
	if err != nil {
		return nil
	}
	return &Douyu{id, nil}
}

type liveResultData struct {
	Status string `json:"room_status"`
}

type liveResult struct {
	Error int            `json:"error"`
	Data  liveResultData `json:"data"`
}

func (p *Douyu) GetLiveStatus() (bool, error) {
	url := fmt.Sprintf("http://open.douyucdn.cn/api/RoomApi/room/%d", p.room)
	var result liveResult
	err := GetJson(url, &result)
	if err != nil {
		return false, err
	}
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

func (d *Douyu) send() {

}

func (d *Douyu) init() error {
	conn, err := net.Dial("tcp", "openbarrage.douyutv.com:8601")
	if err != nil {
		fmt.Println(err)
		return err
	}
	d.conn = conn.(*net.TCPConn)
	return nil
}
