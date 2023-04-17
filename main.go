package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"rsc.io/pdf"
)

func main() {
	inputFile := "example.pdf"
	outputFile := "extracted_text.txt"

	// PDF 파일에서 텍스트 추출
	err := extractText(inputFile, outputFile)
	if err != nil {
		fmt.Printf("Error extracting text from PDF: %v\n", err)
		return
	}

	// 추출된 텍스트 출력
	text, err := ioutil.ReadFile(outputFile)
	if err != nil {
		fmt.Printf("Error reading extracted text file: %v\n", err)
		return
	}

	fmt.Printf("Extracted text:\n\n%s\n", string(text))
}

func extractText(inputFile, outputFile string) error {
	// PDF 파일 열기
	f, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer f.Close()

	// 파일 크기 확인
	fileInfo, err := f.Stat()
	if err != nil {
		return err
	}

	// PDF 파일 읽기
	pdfReader, err := pdf.NewReader(f, fileInfo.Size())
	if err != nil {
		return err
	}

	// 추출된 텍스트 저장할 파일 생성
	o, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer o.Close()

	// 각 페이지에서 텍스트 추출
	numPages := pdfReader.NumPage()
	for pageIndex := 1; pageIndex <= numPages; pageIndex++ {
		page := pdfReader.Page(pageIndex)
		if page == (pdf.Page{}) {
			return fmt.Errorf("failed to retrieve page %d", pageIndex)
		}

		// 해당 페이지에서 텍스트 추출
		text := ""
		content := page.Content()
		for _, t := range content.Text {
			text += t.S
		}

		// 추출된 텍스트를 파일에 쓰기
		if _, err := o.WriteString(text + "\n"); err != nil {
			return err
		}
	}

	return nil
}
