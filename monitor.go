package danmu

import (
	"fmt"
	"net/url"
	"sync"
	"time"
)

var (
	LiveCheckInterval = 10 * time.Minute
	KeepAliveInterval = time.Minute
)

const (
	DefaultDanmuPool = 1024
)

type DanmuClient interface {
	IsLive() (bool, error)
	Close()
	StartReceive() (<-chan *string, error)
}

type Factory func(url string) (DanmuClient, error)

type Monitor struct {
	factorys map[string]Factory
	agents   map[string]DanmuClient
	sync.Mutex
}

func NewMonitor() *Monitor {
	return &Monitor{factorys: make(map[string]Factory), agents: make(map[string]DanmuClient)}
}

func (m *Monitor) Register(domain string, factory Factory) {
	m.Lock()
	defer m.Unlock()
	m.factorys[domain] = factory
}

func (m *Monitor) Obeserve(_url string) (<-chan *string, error) {
	u, err := url.Parse(_url)
	if err != nil {
		return nil, err
	}
	factory := m.factorys[u.Host]
	if factory == nil {
		return nil, fmt.Errorf("not found client for %s", _url)
	}
	client, err := factory(_url)
	if err != nil {
		return nil, err
	}

	return client.StartReceive()
}

var defaultMonitor = NewMonitor()

func Obeserve(_url string) (<-chan *string, error) {
	return defaultMonitor.Obeserve(_url)
}

func Register(domain string, factory Factory) {
	defaultMonitor.Register(domain, factory)
}

func init() {
	defaultMonitor.Register("www.pandatv.com", NewPanda)
}
