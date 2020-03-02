package main

import (
	"os"
)

func main() {
	_, error := os.Open("no-file.txt")

	if error != nil {
		//fmt.Println("error happened", error)
		// log.Fatalln(error)
		panic(error)
	}
}
