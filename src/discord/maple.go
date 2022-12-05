package main

import (
	"log"

	"github.com/sclevine/agouti"
)

const MAPLE_EVENT = "https://maplestory.nexon.co.jp/notice/list/event/"

func ScrapingEventInfo() (err error) {
	// ChromeDriver
	driver := agouti.ChromeDriver(
		agouti.ChromeOptions("args", []string{
			"no-sandbox",
			"--disable-gpu",
			"--disable-dev-shm-usage",
			"--window-size=1280,800",
		}),
		agouti.Debug,
	)

	if err := driver.Start(); err != nil {
		log.Print(err)
		return err
	}
	defer driver.Stop()

	page, err := driver.NewPage()
	if err != nil {
		log.Print(err)
		return err
	}

	if err := page.Navigate(MAPLE_EVENT); err != nil {
		return err
	}

	return nil
}
