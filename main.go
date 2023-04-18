// main.go

package main

import (
	"fmt"
	"songtomtom/pdf-scrapper/pdfscrapper"
)

const (
	homepageUrl = "https://www.suneung.re.kr/boardCnts/list.do?type=default&page=1&searchStr=&m=0403&C06=&boardID=1500234&C05=&C04=&C03=&C02=&searchType=S&C01=&s=suneung"
)

func main() {

	htmlString, err := pdfscrapper.FetchHTML(homepageUrl)
	if err != nil {
		panic(err)
	}

	err = pdfscrapper.FindDivCntBody(htmlString)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
