package novelai

import (
	"fmt"
	"grpc-conoha/config"
	"log"
	"time"

	"github.com/sclevine/agouti"
)

func ExecuteNovelAI(spell string) {
	// ChromeDriver
	driver := agouti.ChromeDriver(
		agouti.ChromeOptions("args", []string{
			// "--headless", コンテナ環境だと正常に動作しない
			"no-sandbox",
			"--disable-gpu",
			"--disable-dev-shm-usage",
			"--window-size=1280,800",
		}),
		agouti.Debug,
	)

	if err := driver.Start(); err != nil {
		log.Fatal(err)
	}
	defer driver.Stop()

	page, err := driver.NewPage()
	if err != nil {
		log.Fatal(err)
	}

	// NovelAIホームページに遷移
	if err := page.Navigate("https://novelai.net/"); err != nil {
		log.Fatal(err)
	}
	time.Sleep(3 * time.Second)
	// ログイン
	if err := page.FindByLink("Login").Click(); err != nil {
		log.Fatal(err)
	}
	emailForm := page.FindByID("username")
	passwordForm := page.FindByID("password")

	emailForm.Fill(config.Config.NovelAIEmail)
	passwordForm.Fill(config.Config.NovelAIPassword)
	time.Sleep(3 * time.Second)

	if err := page.FindByButton("Sign In").Click(); err != nil {
		log.Fatal(err)
	}
	if err := page.Navigate("https://novelai.net/image"); err != nil {
		log.Fatal(err)
	}

	time.Sleep(3 * time.Second)
	promptForm := page.FindByID("prompt-input-0")
	promptForm.Fill(spell)

	time.Sleep(3 * time.Second)
	generate := page.Find("div.sc-75a56bc9-36.jaHANM")
	if err := generate.Click(); err != nil {
		log.Fatal(err)
	}
	// 画像生成まで待つ
	time.Sleep(time.Second * 10)
	// ダウンロード
	imageDownload := page.FindByXPath("//*[@id=\"__next\"]/div[2]/div[3]/div/div[1]/div[2]/div[2]/div/div/div[2]/div[2]/button[2]/div")
	if err := imageDownload.Click(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("save")
	time.Sleep(10 * time.Second)
}
