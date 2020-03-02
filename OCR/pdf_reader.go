package main

import (
	"fmt"

	"github.com/otiai10/gosseract"
)

//let imagePath := "/Users/pratima/Documents/POC/goworkspace/src/github.com/go-programming/01/OCR/files/input.pdf"
func main() {
	fmt.Println("You about to start OCR")

	//readImage()
	client := gosseract.NewClient()
	defer client.Close()
	//client.SetImage("files/image_with_text.png")
	client.SetImage("files/image_with_text.png")
	text, err := client.Text()
	fmt.Println("error", err)
	fmt.Println(text)
}

// func readImage() {
// 	reader, err := os.Open(imagePath)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println(reader)
// 	defer reader.Close()
// }
