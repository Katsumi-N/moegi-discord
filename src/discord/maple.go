package main

import (
	"log"
	"math"
	"strings"

	"github.com/sclevine/agouti"
)

const MAPLE_EVENT = "https://maplestory.nexon.co.jp/notice/list/event/"

type MapleInfo struct {
	Title       string
	Date        string
	Description string
	Url         string
}

func ScrapingEventInfo() (*[]MapleInfo, int, error) {
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
		return nil, 0, err
	}
	defer driver.Stop()

	page, err := driver.NewPage()
	if err != nil {
		log.Print(err)
		return nil, 0, err
	}

	if err := page.Navigate(MAPLE_EVENT); err != nil {
		return nil, 0, err
	}

	if err := page.FirstByClass("notice-list").First("tbody").First("tr").Find("a").Click(); err != nil {
		return nil, 0, err
	}

	menu := page.FindByClass("menu")
	menu_num, _ := menu.All("li").Count()
	retInfo := make([]MapleInfo, menu_num)

	for i := 0; i < menu_num; i++ {
		eventUrl, _ := menu.All("li").At(i).First("a").Attribute("href")
		eventClass := strings.Split(eventUrl, "#")
		date, _ := page.FindByClass(eventClass[1]).First("h3").Text()
		title, _ := page.FindByClass(eventClass[1]).First("img").Attribute("alt")
		txt, err := page.FindByClass(eventClass[1]).First("p").Text()
		if err != nil {
			return nil, menu_num, err
		}

		info := MapleInfo{
			Title:       title,
			Description: txt[:int(math.Min(float64((len(txt))), 1000))],
			Url:         eventUrl,
			Date:        date,
		}
		retInfo[i] = info
	}
	return &retInfo, menu_num, nil

	// noticeList := page.FindByClass("notice-list").First("tr").

	// for i := 0; i < eventNum; i++ {
	// 	tr := noticeList.At(i + 1)
	// 	url, err := tr.Find("a").Attribute("href")
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	row, err := tr.Text()
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	title := strings.Split(row, "\n")[1]
	// 	date := strings.Split(row, "\n")[2]

	// 	// イベント本文を取得
	// 	if err := page.Navigate(url); err != nil {
	// 		return nil, err
	// 	}
	// 	desc, _ := page.FindByClass("txt-area").Text()
	// 	info := MapleInfo{
	// 		Title:       title,
	// 		Description: desc,
	// 		Url:         url,
	// 		Date:        date,
	// 	}
	// 	retInfo[i] = info
	// 	page.Back()
	// 	time.Sleep(300 * time.Millisecond)

	// }

	// return &retInfo, nil
}
