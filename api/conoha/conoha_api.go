package conoha

import (
	"bytes"
	"encoding/json"
	"fmt"
	"grpc-conoha/config"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

const identityURL = "https://identity.tyo2.conoha.io/v2.0/"
const computeURL = "https://compute.tyo2.conoha.io/v2/"
const accountURL = "https://account.tyo2.conoha.io/v1/"
const imageURL = "https://image-service.tyo2.conoha.io/v2/"

type JsonAccess struct {
	Access JsonToken `json:"access"`
}
type JsonToken struct {
	Token TokenInfo `json:"token"`
}

type TokenInfo struct {
	Id       string `json:"id"`
	IssuedAt string `json:"issued_at"`
	Expires  string `json:"expires"`
}

func doRequest(method, base string, urlPath string, tokenId string, data string, query map[string]string) (body []byte, statuscode int, err error) {
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
	log.Printf("action=doRequest endpoint=%s", endpoint)
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

// トークンの取得
func GetToken(userName string, password string, tenantId string) string {
	body := fmt.Sprintf("{\"auth\":{\"passwordCredentials\":{\"username\":\"%s\",\"password\":\"%s\"},\"tenantId\":\"%s\"}}",
		userName, password, tenantId)
	resp, _, err := doRequest("POST", identityURL, "tokens", "", body, map[string]string{})
	if err != nil {
		log.Fatal(err)
	}
	var access JsonAccess
	err = json.Unmarshal(resp, &access)
	return access.Access.Token.Id
}

func StartServer(token string) (resBody []byte, statusCode int) {
	log.Printf("Start server")
	// サーバーの状態を確認
	status, flavorId := GetServerStatus(token)
	log.Print(status)
	// メモリを4gbに変更
	if flavorId != config.Config.Flavor4gb {
		_, statusCode = ChangeServerFlavor(token, "1gb", "4gb")
		fmt.Printf("status %d", statusCode)
		if statusCode != 202 {
			log.Print("error at ChangeServerFlavor")
			return nil, statusCode
		}
		// now := time.Now()
		// statusが"VERIFY_RESIZE"になってからconfirmする
		// 約6分かかる
		for {
			time.Sleep(30 * time.Second)
			status, _ := GetServerStatus(token)
			if status == "VERIFY_RESIZE" {
				break
			}
		}
		// log.Printf("statusの変更にかかった時間: %vms", time.Since(now))

		_, resizeStatusCode := ConfirmResize(token)
		if resizeStatusCode != 204 {
			time.Sleep(10 * time.Second) // 10秒待ってから再リクエスト
			_, _ = ConfirmResize(token)
		}
	}

	// VERIFY_RESIZE->SHUTOFF まで約90秒
	for {
		status, _ := GetServerStatus(token)
		if status == "SHUTOFF" {
			break
		}
		time.Sleep(10 * time.Second)
	}

	status, flavorId = GetServerStatus(token)
	if status != "SHUTOFF" || flavorId != config.Config.Flavor4gb {
		// リトライを求めるHTTPステータスは503(Service Unavailable)にしておく
		return []byte(status), 503
	}

	url := config.Config.TenantId + "/servers/" + config.Config.ServerId + "/action"
	body := fmt.Sprintf("{\"os-start\":\"null\"}")
	resBody, statusCode, err := doRequest("POST", computeURL, url, token, body, map[string]string{})
	if err != nil {
		log.Print(err)
	}
	return resBody, statusCode
}

func StopServer(token string) (resBody []byte, status int) {
	url := config.Config.TenantId + "/servers/" + config.Config.ServerId + "/action"
	body := fmt.Sprintf("{\"os-stop\":\"null\"}")
	resBody, status, err := doRequest("POST", computeURL, url, token, body, map[string]string{})
	if err != nil {
		log.Print(err)
	}
	return resBody, status
}

func RebootServer(token string) (resBody []byte, status int) {
	url := config.Config.TenantId + "/servers/" + config.Config.ServerId + "/action"
	body := fmt.Sprintf("{\"reboot\":{\"type\":\"SOFT\"}}")
	resBody, status, err := doRequest("POST", computeURL, url, token, body, map[string]string{})
	if err != nil {
		log.Print(err)
	}
	return resBody, status
}

// サーバーの状態を取得
type ServerInfo struct {
	Server struct {
		Status string `json:"status"`
		Flavor struct {
			ID string `json:"id"`
		}
	} `json:"server"`
}

// サーバーの状態を取得
func GetServerStatus(token string) (status string, flavorId string) {
	url := config.Config.TenantId + "/servers/" + config.Config.ServerId
	resp, _, err := doRequest("GET", computeURL, url, token, "", map[string]string{})
	if err != nil {
		log.Fatal(err)
	}
	var server ServerInfo
	err = json.Unmarshal(resp, &server)
	return server.Server.Status, server.Server.Flavor.ID
}

//サーバーのflavor(メモリプラン)の変更
func ChangeServerFlavor(token string, now string, to string) (resBody []byte, status int) {
	url := config.Config.TenantId + "/servers/" + config.Config.ServerId + "/action"
	var changeFlavor string
	if now == "1gb" && to == "4gb" {
		changeFlavor = config.Config.Flavor4gb
	} else if now == "4gb" && to == "1gb" {
		changeFlavor = config.Config.Flavor1gb
	}
	data := fmt.Sprintf("{\"resize\": {\"flavorRef\": \"%s\"}}", changeFlavor)
	resBody, status, err := doRequest("POST", computeURL, url, token, data, map[string]string{})
	if err != nil {
		log.Fatal(err)
	}
	return resBody, status
}

func ConfirmResize(token string) (resBody []byte, status int) {
	url := config.Config.TenantId + "/servers/" + config.Config.ServerId + "/action"
	data := fmt.Sprintf("{\"confirmResize\": null}")
	resBody, status, err := doRequest("POST", computeURL, url, token, data, map[string]string{})
	if err != nil {
		log.Print(err)
	}
	return resBody, status
}
