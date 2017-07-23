package danmu

type DanmuClient interface {
	GetLiveStatus() (bool, error)
	Close()
	StartReceive()
}

type Monitor struct {
}

func Add(url string) (<-chan (string), error) {
	return nil, nil
}
