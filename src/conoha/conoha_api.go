package main

import (
	"encoding/json"
	"fmt"
	"grpc-conoha/config"
	conohapb "grpc-conoha/pkg/grpc"
	"log"
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

// トークンの取得
func GetToken(userName string, password string, tenantId string) string {
	body := fmt.Sprintf("{\"auth\":{\"passwordCredentials\":{\"username\":\"%s\",\"password\":\"%s\"},\"tenantId\":\"%s\"}}",
		userName, password, tenantId)
	resp, _, err := DoRequest("POST", identityURL, "tokens", "", body, map[string]string{})
	if err != nil {
		log.Fatal(err)
	}
	var access JsonAccess
	err = json.Unmarshal(resp, &access)
	return access.Access.Token.Id
}

func StartServer(token string, stream conohapb.ConohaService_MinecraftServer) (resBody []byte, statusCode int) {
	// log.Printf("Start server")
	// サーバーの状態を確認
	status, flavorId := GetServerStatus(token)
	// メモリを4gbに変更
	if flavorId != config.Config.Flavor4gb {
		_, statusCode = ChangeServerFlavor(token, "1gb", "4gb")
		fmt.Printf("status %d", statusCode)
		if statusCode != 202 {
			log.Print("error at ChangeServerFlavor")
			return nil, statusCode
		}
		// statusが"VERIFY_RESIZE"になってからconfirmする
		// 約6分かかる
		t := time.Now()
		t_expect := t.Add(time.Duration(8) * time.Minute)
		if err := stream.Send(&conohapb.MinecraftResponse{
			Message: fmt.Sprintf("リサイズ処理がスタートしました．予定時刻 %d:%d", t_expect.Hour(), t_expect.Minute()),
			Health:  true,
		}); err != nil {
			return nil, 503
		}
		for {
			time.Sleep(30 * time.Second)
			status, _ := GetServerStatus(token)
			if status == "VERIFY_RESIZE" {
				break
			}
		}

		_, resizeStatusCode := ConfirmResize(token)
		if resizeStatusCode != 204 {
			time.Sleep(10 * time.Second) // 10秒待ってから再リクエスト
			_, _ = ConfirmResize(token)
		}

		// VERIFY_RESIZE->SHUTOFF まで約90秒
		t = time.Now()
		t_expect = t.Add(time.Duration(2) * time.Minute)
		if err := stream.Send(&conohapb.MinecraftResponse{
			Message: "リサイズ処理が終了しました．起動処理を開始します．",
			Health:  true,
		}); err != nil {
			return nil, 503
		}
	}

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
	resBody, statusCode, err := DoRequest("POST", computeURL, url, token, body, map[string]string{})
	if err != nil {
		log.Print(err)
	}
	return resBody, statusCode // 202
}

func StopServer(token string, stream conohapb.ConohaService_MinecraftServer) (resBody []byte, statusCode int) {
	status, flavorId := GetServerStatus(token)
	if status == "ACTIVE" {
		t := time.Now()
		t_expect := t.Add(time.Duration(8) * time.Minute)
		if err := stream.Send(&conohapb.MinecraftResponse{
			Message: fmt.Sprintf("シャットダウンを開始します．予定時刻 %d:%d", t_expect.Hour(), t_expect.Minute()),
			Health:  true,
		}); err != nil {
			return nil, 503
		}

		url := config.Config.TenantId + "/servers/" + config.Config.ServerId + "/action"
		body := fmt.Sprintf("{\"os-stop\":\"null\"}")
		_, _, err := DoRequest("POST", computeURL, url, token, body, map[string]string{})
		if err != nil {
			log.Print(err)
		}
	}
	for {
		time.Sleep(10 * time.Second)
		status, _ := GetServerStatus(token)
		if status == "SHUTOFF" {
			break
		}
	}

	// メモリを1gbに変更
	if err := stream.Send(&conohapb.MinecraftResponse{
		Message: "リサイズ処理がスタートしました.",
		Health:  true,
	}); err != nil {
		return nil, 503
	}
	if flavorId == config.Config.Flavor4gb {
		_, statusCode = ChangeServerFlavor(token, "4gb", "1gb")
		fmt.Printf("status %d", statusCode)
		if statusCode != 202 {
			log.Print("error at ChangeServerFlavor")
			return nil, statusCode
		}
		// statusが"VERIFY_RESIZE"になってからconfirmする
		// 約6分かかる
		for {
			time.Sleep(30 * time.Second)
			status, _ := GetServerStatus(token)
			if status == "VERIFY_RESIZE" {
				break
			}
		}
		_, resizeStatusCode := ConfirmResize(token)
		if resizeStatusCode != 204 {
			time.Sleep(10 * time.Second) // 10秒待ってから再リクエスト
			_, _ = ConfirmResize(token)
		}

		// VERIFY_RESIZE->SHUTOFF まで約90秒待つ
		for {
			status, _ := GetServerStatus(token)
			if status == "SHUTOFF" {
				break
			}
			time.Sleep(10 * time.Second)
		}

		status, flavorId = GetServerStatus(token)
		if status != "SHUTOFF" || flavorId != config.Config.Flavor1gb {
			// リトライを求めるHTTPステータスは503(Service Unavailable)にしておく
			return []byte(status), 503
		}
	}

	return resBody, statusCode // 202
}

func RebootServer(token string) (resBody []byte, status int) {
	url := config.Config.TenantId + "/servers/" + config.Config.ServerId + "/action"
	body := fmt.Sprintf("{\"reboot\":{\"type\":\"SOFT\"}}")
	resBody, status, err := DoRequest("POST", computeURL, url, token, body, map[string]string{})
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
	resp, _, err := DoRequest("GET", computeURL, url, token, "", map[string]string{})
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
	resBody, status, err := DoRequest("POST", computeURL, url, token, data, map[string]string{})
	if err != nil {
		log.Fatal(err)
	}
	return resBody, status
}

func ConfirmResize(token string) (resBody []byte, status int) {
	url := config.Config.TenantId + "/servers/" + config.Config.ServerId + "/action"
	data := fmt.Sprintf("{\"confirmResize\": null}")
	resBody, status, err := DoRequest("POST", computeURL, url, token, data, map[string]string{})
	if err != nil {
		log.Print(err)
	}
	return resBody, status
}
