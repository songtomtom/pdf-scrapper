package main

import (
	"fmt"
	"os"
	"songtomtom/pdf-scrapper/pdfscrapper"
)

const (
	baseURL = "https://www.suneung.re.kr/boardCnts/list.do?type=default&page=%d&m=0403&C06=&boardID=1500234&C05=&C04=&C03=&C02=%EC%98%81%EC%96%B4&searchType=S&C01=&s=suneung"
)

func main() {
	startPage := 1
	endPage := 20

	err := pdfscrapper.ScrapePages(baseURL, startPage, endPage)
	if err != nil {
		fmt.Printf("Error scraping pages: %v\n", err)
		os.Exit(1)
	}
}
