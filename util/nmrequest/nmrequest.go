package nmrequest

import (
	"net/http"
	"strings"

	"github.com/Zhousiru/Waver/util/nmcrypto"
)

// DoRequest - 发起请求
func DoRequest(url string, params map[string]string, reqCallback func(*http.Request)) (*http.Response, error) {
	postData, err := nmcrypto.EncryptRequest(params)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(postData))
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()

	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/68.0.3440.106 Safari/537.36")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if reqCallback != nil {
		reqCallback(req)
	}

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
