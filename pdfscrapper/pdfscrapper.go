package pdfscrapper

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type TableRow struct {
	Id        string
	Year      string
	Subject   string
	Contents  string
	Date      string
	ViewCount string
	Files     []string
}

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

	var rows []TableRow

	doc.Find("div.cntBody.clearfix tr").Each(
		func(i int, s *goquery.Selection) {
			row := TableRow{}

			// 여기에서 필요한 정보를 추출하고 처리합니다.
			// 예를 들어, 다음 코드는 선택한 태그의 모든 텍스트를 출력합니다.

			if s.Find("td").Eq(2).Text() != "영어" {
				return
			}

			row.Id = s.Find("td").Eq(0).Text()
			row.Year = s.Find("td").Eq(1).Text()
			row.Subject = s.Find("td").Eq(2).Text()
			row.Contents = s.Find("td .link a").Text()
			row.Date = s.Find("td").Eq(4).Text()
			row.ViewCount = s.Find("td").Eq(5).Text()

			s.Find("td a").Each(
				func(i int, a *goquery.Selection) {
					title, _ := a.Attr("title")
					if title != "" {
						row.Files = append(row.Files, title)
					}
				},
			)
			rows = append(rows, row)
		},
	)

	fmt.Print(rows)

	return nil
}
