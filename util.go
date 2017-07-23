package danmu

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func GetJson(url string, v interface{}) error {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("get error ", url, err)
		return err
	}

	if resp.StatusCode != http.StatusOK || resp.Body == nil {
		fmt.Println("status code error ", url, resp.StatusCode)
		return errors.New(fmt.Sprint("status code ", resp.StatusCode))
	}

	defer resp.Body.Close()
	d := json.NewDecoder(resp.Body)
	err = d.Decode(v)
	if err != nil {
		fmt.Println("json error ", url, err)
		return err
	}

	return nil
}
