package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	url := "https://www.suneung.re.kr/boardCnts/fileDown.do?fileSeq=f7683448be641c3b3f47900702ea29f2"
	response, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile("file.pdf", body, 0644)
	if err != nil {
		panic(err)
	}

	fmt.Println("파일 다운로드 완료")
}
