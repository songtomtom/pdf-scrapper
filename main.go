// main.go

package main

import (
	"fmt"
	"songtomtom/pdf-scrapper/pdfscrapper"
)

const (
	baseURL = "https://www.suneung.re.kr/boardCnts/list.do?type=default&m=0403&boardID=1500234&searchType=S&C02=%EC%98%81%EC%96%B4&s=suneung&page=%d"
)

func main() {

	startPage := 1
	endPage := 5 // 원하는 마지막 페이지 번호를 설정하세요.

	for page := startPage; page <= endPage; page++ {
		url := fmt.Sprintf(baseURL, page)
		htmlString, err := pdfscrapper.FetchHTML(url)
		if err != nil {
			fmt.Printf("Error fetching page %d: %v\n", page, err)
			continue
		}

		err = pdfscrapper.FindDivCntBody(htmlString)
		if err != nil {
			fmt.Printf("Error processing page %d: %v\n", page, err)
		}
	}
}
