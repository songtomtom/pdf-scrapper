package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"regexp"
)

func main() {
	res, err := http.Get("https://www.suneung.re.kr/boardCnts/list.do?type=default&page=1&searchStr=&m=0403&C06=&boardID=1500234&C05=&C04=&C03=&C02=&searchType=S&C01=&s=suneung")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// HTML 문서를 파싱합니다.
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// 원하는 DOM 요소를 찾고 조작합니다.
	doc.Find("a[onclick^=\"fn_fileDown\"]").Each(
		func(i int, s *goquery.Selection) {
			// 선택된 각 요소에서 onclick 속성 값을 가져옵니다.
			res, exists := s.Attr("onclick")
			if exists {
				// 정규식을 사용하여 필요한 문자열만 추출합니다.
				re := regexp.MustCompile(`'([\w\d]+)'`)
				matches := re.FindStringSubmatch(res)

				if len(matches) > 1 {
					fmt.Printf("Found result %d: %s\n", i, matches[1])
				}
			}
		},
	)
}
