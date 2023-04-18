package pdfscrapper

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func FetchHTML(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func FindDivCntBody(htmlString string) error {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlString))
	if err != nil {
		return err
	}

	doc.Find("div.cntBody.clearfix").Each(
		func(i int, s *goquery.Selection) {
			// 여기에서 필요한 정보를 추출하고 처리합니다.
			// 예를 들어, 다음 코드는 선택한 태그의 모든 텍스트를 출력합니다.
			fmt.Println(s.Text())
		},
	)

	return nil
}
