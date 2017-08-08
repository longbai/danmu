package danmu

import (
	"os"
)

var file *os.File

func init() {
	file, _ = os.OpenFile("/Users/long/temp/13fdsfasfas", os.O_RDWR|os.O_CREATE, 0755)
}

func dump(data []byte) {
	file.Write(data)
	file.Sync()
}
