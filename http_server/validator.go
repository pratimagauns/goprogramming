package main

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

var months = map[string]string{
	"JAN": "01",
	"FEB": "02",
	"MAR": "03",
	"APR": "04",
	"MAY": "05",
	"JUN": "06",
	"JUL": "07",
	"AUG": "08",
	"SEP": "09",
	"OCT": "10",
	"NOV": "11",
	"DEC": "12",
}

type validationError struct {
	info string
}

func (ce validationError) Error() string {
	return fmt.Sprintf("here is the error: %v", ce.info)
}

// GetURLValidator will provide the download url with the param provided
// Validates data
func GetURLValidator(params string) (string, validationError) {
	const queryPrefix = "https://www1.nseindia.com/content/historical/EQUITIES"

	params = strings.ToUpper(params)
	array := []rune(params)

	if len(array) == 0 {
		return "", validationError{info: "No query parameters provided"}
	}
	if len(array) < 8 || len(array) > 9 {
		return "", validationError{info: "Invalid parameters provided"}
	}

	re := regexp.MustCompile(`[A-Z]{3}`)
	if len(re.FindAllString(params, -1)) == 0 {
		fmt.Println("Invalid month parameters provided")
		return "", validationError{info: "1- Invalid month parameters provided"}
	}

	month := re.FindAllString(params, -1)[0]
	if len(month) < 0 {
		fmt.Println("Invalid month parameters provided")
		return "", validationError{info: "2- Invalid month parameters provided"}
	}

	indx := re.FindStringIndex(params)
	if indx[0] != 2 && indx[0] != 1 {
		fmt.Println("Month parameters wrong location", indx[0])
		return "", validationError{info: "Month parameters wrong location"}
	}

	fmt.Println("Month Found Correct", month)

	day := string(array[0:indx[0]])
	//Check for Index Out of Bound
	if len(day) < 0 {
		fmt.Println("Invalid day parameters")
		return "", validationError{info: "Invalid day parameters"}
	}
	year := string(array[indx[1]:len(array)])
	//Check for Index Out of Bound
	if len(year) != 4 {
		fmt.Println("Invalid year parameters")
		return "", validationError{info: "Invalid year parameters"}
	}

	//fmt.Println("Found...", day, year)

	if len(day) == 1 {
		day = fmt.Sprintf("0%s", day)
		params = fmt.Sprintf("%s%s%s", day, month, year)
	}
	//RFC3339     = "2006-01-02T15:04:05Z07:00"
	formattedDate := fmt.Sprintf("%s-%s-%sT00:00:00.000Z", year, months[month], day)

	fmt.Println(formattedDate)
	t, err := time.Parse(time.RFC3339, formattedDate) //Check format possible
	if err != nil {
		fmt.Println("Invalid Date")
		return "", validationError{info: "Invalid Date"}
	}
	fmt.Println(t)

	filename := fmt.Sprintf("cm%sbhav.csv.zip", params)
	downloadURL := fmt.Sprintf("%s/%s/%s/%s", queryPrefix, year, month, filename)

	fmt.Println(downloadURL)
	return downloadURL, validationError{info: ""}
}
