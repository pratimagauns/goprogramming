package main

import (
	"fmt"
	"runtime"
	"sync"
)

var waitG sync.WaitGroup

func main() {
	fmt.Println("OS \t\t", runtime.GOOS)
	fmt.Println("ARCH \t\t", runtime.GOARCH)
	fmt.Println("CPUs \t\t", runtime.NumCPU())
	fmt.Println("Go Routines: \t", runtime.NumGoroutine())

	waitG.Add(1)
	go foo()

	fmt.Println("CPUs \t\t", runtime.NumCPU())
	fmt.Println("Go Routines \t\t", runtime.NumGoroutine())
	waitG.Wait()
}

func foo() {
	for i := 0; i < 10; i++ {
		fmt.Println("in foo:", i)
	}
	waitG.Done()
}
