// main.go

package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/unidoc/unipdf/v3/common/license"
	"os"
	"songtomtom/pdf-scrapper/pkg/pdfreader"
)

func init() {

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	licenseKey := os.Getenv("UNIDOC_LICENSE_API_KEY")
	if licenseKey == "" {
		panic("UNIDOC_LICENSE_API_KEY environment variable not set")
	}

	err = license.SetMeteredKey(licenseKey)
	if err != nil {
		panic(err)
	}
}

func main() {

	if len(os.Args) < 2 {
		fmt.Printf("Usage: go run main.go input.pdf\n")
		os.Exit(1)
	}

	inputPath := os.Args[1]

	err := pdfreader.OutputPdfText(inputPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
