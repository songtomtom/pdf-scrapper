package pdfscrapper

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

type TableRow struct {
	Id         string
	Year       string
	Subject    string
	Contents   string
	Date       string
	ViewCount  string
	Files      []string
	FileHashes []string
}

func ScrapePages(baseURL string, startPage int, endPage int) error {
	err := os.RemoveAll("./files")
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	for page := startPage; page <= endPage; page++ {
		wg.Add(1)
		go func(page int) {
			defer wg.Done()

			url := fmt.Sprintf(baseURL, page)
			htmlString, err := FetchHTML(url)
			if err != nil {
				fmt.Printf("Error fetching HTML from %s: %v\n", url, err)
				return
			}

			err = FindDivCntBody(htmlString)
			if err != nil {
				fmt.Printf("Error processing page %d: %v\n", page, err)
			}
		}(page)
	}

	wg.Wait()

	return nil
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

func createDirectoryIfNotExists(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}

func downloadFile(url string, destPath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer out.Close()

	fileSize := resp.ContentLength
	fmt.Printf("Downloading %s, size: %d bytes\n", destPath, fileSize)

	progressChan := make(chan int64)
	go func() {
		for {
			select {
			case current := <-progressChan:
				percentage := float64(current) / float64(fileSize) * 100
				fmt.Printf("\rProgress: %d/%d bytes, %.2f%%", current, fileSize, percentage)
			}
		}
	}()

	reader := io.TeeReader(resp.Body, &writeCounter{0, progressChan})
	_, err = io.Copy(out, reader)
	if err != nil {
		return err
	}

	close(progressChan)
	fmt.Println("\nDownload complete")

	return nil
}

type writeCounter struct {
	total        int64
	progressChan chan int64
}

func (wc *writeCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.total += int64(n)
	wc.progressChan <- wc.total
	return n, nil
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
					onclick, _ := a.Attr("onclick")
					if onclick != "" {
						hashRegexp := regexp.MustCompile(`'([a-fA-F0-9]+)'`)
						hashMatches := hashRegexp.FindStringSubmatch(onclick)
						if len(hashMatches) > 1 {
							row.FileHashes = append(row.FileHashes, hashMatches[1])
						}

						fileURL := fmt.Sprintf(
							"https://www.suneung.re.kr/boardCnts/fileDown.do?fileSeq=%s", hashMatches[1],
						)
						row.FileHashes = append(row.FileHashes, hashMatches[1])

						dirPath := filepath.Join("./files", row.Year)

						err := createDirectoryIfNotExists(dirPath)
						if err != nil {
							fmt.Printf("Error creating directory: %v\n", err)
							return
						}

						destPath := filepath.Join(dirPath, title)
						err = downloadFile(fileURL, destPath)
						if err != nil {
							fmt.Printf("Error downloading file: %v\n", err)
							return
						}
					}
				},
			)
			rows = append(rows, row)
		},
	)

	return nil
}
