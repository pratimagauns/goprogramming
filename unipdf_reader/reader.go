package main

import (
	"archive/zip"
	"fmt"
	"image/jpeg"
	"os"

	"github.com/unidoc/unipdf/extractor"
	"github.com/unidoc/unipdf/model"
	// pdf "github.com/unidoc/unipdf/model"
)

//go run reader.go files/input_pdf_3.pdf files/output.zip

var colorspaces = map[string]int{}
var filters = map[string]int{}

func main() {
	// Enable console debug-level logging when debugging:.
	//unicommon.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelDebug))

	if len(os.Args) < 3 {
		fmt.Printf("Syntax: go run pdf_extract_images.go input.pdf output.zip\n")
		os.Exit(1)
	}

	inputPath := os.Args[1]
	outputPath := os.Args[2]

	fmt.Printf("Input file: %s\n", inputPath)
	err := extractImagesToArchive(inputPath, outputPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// err = outputPdfText(inputPath)
	// if err != nil {
	// 	fmt.Printf("Error: %v\n", err)
	// 	os.Exit(1)
	// }
}

func extractImagesToArchive(inputPath, outputPath string) error {
	f, err := os.Open(inputPath)
	if err != nil {
		return err
	}

	defer f.Close()

	pdfReader, err := model.NewPdfReader(f)
	if err != nil {
		return err
	}

	isEncrypted, err := pdfReader.IsEncrypted()
	if err != nil {
		return err
	}

	// Try decrypting with an empty one.
	if isEncrypted {
		auth, err := pdfReader.Decrypt([]byte(""))
		if err != nil {
			// Encrypted and we cannot do anything about it.
			return err
		}
		if !auth {
			fmt.Println("Need to decrypt with password")
			return nil
		}
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return err
	}
	fmt.Printf("PDF Num Pages: %d\n", numPages)

	// Prepare output archive.
	zipf, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	defer zipf.Close()
	zipw := zip.NewWriter(zipf)

	totalImages := 0
	for i := 0; i < numPages; i++ {
		pageNum := i + 1
		fmt.Println("------------------------------")
		fmt.Printf("-----\nPage %d:\n", pageNum)

		page, err := pdfReader.GetPage(i + 1)
		if err != nil {
			return err
		}

		pextract, err := extractor.New(page)
		if err != nil {
			return err
		}

		//------------------------------ TEXT EXTRACT

		text, err := pextract.ExtractText()
		if err != nil {
			return err
		}

		//fmt.Printf("Page %d:\n", pageNum)
		fmt.Printf("\"%s\"\n", text)

		//------------------------------ TEXT EXTRACT END

		//------------------------------ IMAGE EXTRACT
		pimages, err := pextract.ExtractPageImages(nil)
		if err != nil {
			return err
		}

		fmt.Printf("%d Images\n", len(pimages.Images))
		for idx, img := range pimages.Images {
			fmt.Printf("Image %d - X: %.2f Y: %.2f, Width: %.2f, Height: %.2f\n",
				totalImages+idx+1, img.X, img.Y, img.Width, img.Height)
			fname := fmt.Sprintf("p%d_%d.jpg", i+1, idx)

			gimg, err := img.Image.ToGoImage()
			if err != nil {
				return err
			}

			imgf, err := zipw.Create(fname)
			if err != nil {
				return err
			}
			opt := jpeg.Options{Quality: 100}
			err = jpeg.Encode(imgf, gimg, &opt)
			if err != nil {
				return err
			}
		}
		totalImages += len(pimages.Images)
		fmt.Println("------------------------------")
		//------------------------------ IMAGE EXTRACT
	}
	fmt.Printf("Total: %d images\n", totalImages)

	// Make sure to check the error on Close.
	err = zipw.Close()
	if err != nil {
		return err
	}

	return nil
}

// outputPdfText prints out contents of PDF file to stdout.
// func outputPdfText(inputPath string) error {
// 	f, err := os.Open(inputPath)
// 	if err != nil {
// 		return err
// 	}

// 	defer f.Close()

// 	pdfReader, err := pdf.NewPdfReader(f)
// 	if err != nil {
// 		return err
// 	}

// 	numPages, err := pdfReader.GetNumPages()
// 	if err != nil {
// 		return err
// 	}

// 	fmt.Printf("--------------------\n")
// 	fmt.Printf("PDF to text extraction:\n")
// 	fmt.Printf("--------------------\n")
// 	for i := 0; i < numPages; i++ {
// 		pageNum := i + 1

// 		page, err := pdfReader.GetPage(pageNum)
// 		if err != nil {
// 			return err
// 		}

// 		ex, err := extractor.New(page)
// 		if err != nil {
// 			return err
// 		}

// 		text, err := ex.ExtractText()
// 		if err != nil {
// 			return err
// 		}

// 		fmt.Println("------------------------------")
// 		fmt.Printf("Page %d:\n", pageNum)
// 		fmt.Printf("\"%s\"\n", text)
// 		fmt.Println("------------------------------")
// 	}

// 	return nil
// }
