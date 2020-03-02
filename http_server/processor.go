package main

import (
	"encoding/csv"
	"io"
	"log"
	"net/url"
	"os"
)

// ProcessDateRequest will download the data and save to local storage
// Give a date string
func ProcessDateRequest(queryMap url.Values) ([]string, string) {

	//fmt.Println("Processing...", queryMap.Get("date"))

	downloadURL, validationerr := GetURLValidator(queryMap.Get("date"))
	if validationerr.info != "" {
		//log.Fatal(validationerr)
		return nil, validationerr.info
	}

	//Download file
	filename, err := DownloadFile("./", downloadURL)

	if err != nil {
		//log.Fatal(err)
		return nil, err.Error()
	}

	records := ParseCSV(filename)
	return records, ""
}

// ParseCSV will parse records from the csv file
func ParseCSV(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	records := make(chan []string)
	go func() {
		parser := csv.NewReader(file)

		defer close(records)
		for {
			record, err := parser.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				file.Close()
				log.Fatal(err)
			}

			records <- record
		}
	}()

	recordsList := PrintRecords(records)

	//Delete file after records fetched
	err = os.Remove(filename)
	if err != nil {
		file.Close()
		log.Fatal(err)
	}

	return recordsList
}

// PrintRecords  prints the slice of records from the channel provided
func PrintRecords(records chan []string) []string {
	var recordsList []string
	for record := range records {
		//fmt.Println(record)
		recordsList = append(recordsList, record...)
	}
	return recordsList
}
