package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
)

func fetchHashesFromPage(page int) []string {
	hashes := []string{}

	url := fmt.Sprintf(
		"https://www.suneung.re.kr/boardCnts/list.do?type=default&m=0403&boardID=1500234&searchType=S&C02=%EC%98%81%EC%96%B4&s=suneung&page=%d",
		page,
	)
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("a[onclick^=\"fn_fileDown\"]").Each(
		func(i int, s *goquery.Selection) {
			res, exists := s.Attr("onclick")
			if exists {
				re := regexp.MustCompile(`'([\w\d]+)'`)
				matches := re.FindStringSubmatch(res)

				if len(matches) > 1 {
					hashes = append(hashes, matches[1])
				}
			}
		},
	)

	return hashes
}

func downloadFile(hash string) {
	url := fmt.Sprintf("https://www.suneung.re.kr/boardCnts/fileDown.do?fileSeq=%s", hash)
	response, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	// 'files' 디렉토리가 없으면 생성합니다.
	err = os.MkdirAll("files", 0755)
	if err != nil {
		panic(err)
	}

	// 파일 경로를 'files' 디렉토리로 변경합니다.
	filePath := filepath.Join("files", fmt.Sprintf("%s.zip", hash))
	err = ioutil.WriteFile(filePath, body, 0644)
	if err != nil {
		panic(err)
	}

	fmt.Printf("File %s downloaded\n", filePath)
}

func main() {
	startPage := 1
	endPage := 5

	for page := startPage; page <= endPage; page++ {
		hashes := fetchHashesFromPage(page)
		for _, hash := range hashes {
			downloadFile(hash)
		}
	}
}
