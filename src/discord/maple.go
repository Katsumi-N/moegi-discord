package main

import (
	"log"
	"strings"
	"time"

	"github.com/sclevine/agouti"
)

const MAPLE_EVENT = "https://maplestory.nexon.co.jp/notice/list/event/"

type MapleInfo struct {
	Title       string
	Date        string
	Description string
	Url         string
}

func ScrapingEventInfo(eventNum int) (*[]MapleInfo, error) {
	// ChromeDriver
	driver := agouti.ChromeDriver(
		agouti.ChromeOptions("args", []string{
			"--headless",
			"no-sandbox",
			"--disable-gpu",
			"--disable-dev-shm-usage",
			"--window-size=1280,800",
		}),
		agouti.Debug,
	)

	if err := driver.Start(); err != nil {
		log.Print(err)
		return nil, err
	}
	defer driver.Stop()

	page, err := driver.NewPage()
	if err != nil {
		log.Print(err)
		return nil, err
	}

	if err := page.Navigate(MAPLE_EVENT); err != nil {
		return nil, err
	}

	retInfo := make([]MapleInfo, eventNum)
	noticeList := page.FindByClass("notice-list").All("tr")
	for i := 0; i < eventNum; i++ {
		tr := noticeList.At(i + 1)
		url, err := tr.Find("a").Attribute("href")
		if err != nil {
			return nil, err
		}
		row, err := tr.Text()
		if err != nil {
			return nil, err
		}
		title := strings.Split(row, "\n")[1]
		date := strings.Split(row, "\n")[2]

		// イベント本文を取得
		if err := page.Navigate(url); err != nil {
			return nil, err
		}
		desc, _ := page.FindByClass("txt-area").Text()
		info := MapleInfo{
			Title:       title,
			Description: desc[:300],
			Url:         url,
			Date:        date,
		}
		retInfo[i] = info
		page.Back()
		time.Sleep(300 * time.Millisecond)

	}

	return &retInfo, nil
}
