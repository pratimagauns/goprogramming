package main

import (
	"fmt"
	"os"
	"time"
)

//go run reader.go files/input_pdf_3.pdf files/output.zip

var colorspaces = map[string]int{}
var filters = map[string]int{}

func main() {
	// Enable console debug-level logging when debugging:.

	if len(os.Args) < 3 {
		fmt.Printf("Syntax: go run pdf_extract_images.go input.pdf output.zip\n")
		os.Exit(1)
	}

	downloadURL := os.Args[1]

	//Creating path for saving PDF
	timestamp := time.Now().Unix()
	inputPath := fmt.Sprintf("files/%d.%s", timestamp, "pdf")
	fmt.Println(inputPath)

	//Download File
	err := DownloadFile(
		inputPath,
		downloadURL)

	fmt.Println("File Downloaded")

	if err != nil {
		fmt.Printf("Error Downloading File: %v\n", err)
		os.Exit(1)
	}

	//Creating path for saving Images
	outputPath := fmt.Sprintf("%s/%s_%d", os.Args[2], "output", timestamp)

	//Extracting the images
	err = ExtractImagesToArchive(inputPath, outputPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
