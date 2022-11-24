package conoha

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
)

func DoRequest(method, base string, urlPath string, tokenId string, data string, query map[string]string) (body []byte, statuscode int, err error) {
	client := &http.Client{}
	baseURL, err := url.Parse(base)
	if err != nil {
		return
	}
	apiURL, err := url.Parse(urlPath)
	if err != nil {
		return
	}
	// 相対パス→絶対パス
	endpoint := baseURL.ResolveReference(apiURL).String()
	// log.Printf("action=doRequest endpoint=%s", endpoint)
	//リクエストの作成
	req, err := http.NewRequest(method, endpoint, bytes.NewBufferString(data))
	if err != nil {
		return
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if tokenId != "" {
		req.Header.Add("X-Auth-Token", tokenId)
	}
	// 渡されたクエリをAdd
	q := req.URL.Query()
	for key, value := range query {
		q.Add(key, value)
	}
	// クエリはエンコードが必要
	req.URL.RawQuery = q.Encode()

	// 実行
	resp, err := client.Do(req)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	defer resp.Body.Close()
	// 帰ってきた値のbodyを読み込む
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}

	return body, resp.StatusCode, nil
}
