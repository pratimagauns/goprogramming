package main

import (
	"fmt"
	"image/jpeg"
	"os"

	"github.com/unidoc/unipdf/extractor"
	"github.com/unidoc/unipdf/model"
)

// ExtractImagesToArchive will extract the images
// and save it in the outputPath
func ExtractImagesToArchive(inputPath, outputPath string) error {
	f, err := os.Open(inputPath)
	pdftext := ""
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
	err = os.Mkdir(outputPath, os.ModePerm)
	if err != nil {
		return err
	}

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
		pdftext = fmt.Sprintf("%s \n\n\nPage No: %d \n %s",
			pdftext,
			pageNum,
			text)
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
			fname := fmt.Sprintf("%s/p%d_%d.jpg", outputPath, i+1, idx)

			gimg, err := img.Image.ToGoImage()
			if err != nil {
				return err
			}

			imgf, err := os.Create(fname)
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
		extractText(outputPath, pdftext)
		fmt.Println("------------------------------")
		//------------------------------ IMAGE EXTRACT
	}
	fmt.Printf("Total: %d images\n", totalImages)

	return nil
}

func extractText(outputPath string, text string) {
	path := fmt.Sprintf("%s/%s", outputPath, "text.txt")
	fmt.Println(path)
	f, err := os.Create(path)
	if err != nil {
		fmt.Println(err)
	}

	l, err := f.WriteString(text)
	if err != nil {
		fmt.Println(err)
		f.Close()
	}

	fmt.Println(l, "bytes written successfully")
	err = f.Close()
	if err != nil {
		fmt.Println(err)
	}
}
