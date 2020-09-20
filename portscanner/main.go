package main

import (
	"fmt"
	"net"
	"sort"
)

func worker(ports, results chan int) {
	for p := range ports {
		//fmt.Println(p, " scanning")
		address := fmt.Sprintf("127.0.0.1:%d", p)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			results <- 0
			//fmt.Println(p, " not open")
			continue
		}
		conn.Close()
		//fmt.Println(p, " open")
		results <- p
	}
}

func main() {
	ports := make(chan int, 100)
	results := make(chan int)
	var openports []int

	// Create a worker pool of 100
	for i := 0; i < cap(ports); i++ {
		go worker(ports, results)
	}

	// Add jobs to ports channel
	go func() {
		for i := 1; i <= 1024; i++ {
			ports <- i
		}
	}()

	// Get open ports from results and add to openports array
	for i := 0; i < 1024; i++ {
		port := <-results
		if port != 0 {
			openports = append(openports, port)
		}
	}

	// Close channels
	close(ports)
	close(results)
	fmt.Println("All scanned")
	sort.Ints(openports)
	for _, port := range openports {
		fmt.Printf("%d open\n", port)
	}
}
