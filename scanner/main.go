package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func main(){
	args := os.Args[1:]
	if len(args) != 1 {
		fmt.Println("Please type a search word")
		return
	}
	query := args[0]

	rx := regexp.MustCompile(`[^a-z]+`)
	in := bufio.NewScanner(os.Stdin)
	in.Split(bufio.ScanWords)

	words := make(map[string]bool)
	for in.Scan() {
		word := strings.ToLower(in.Text())
		word = rx.ReplaceAllString(word, "")
		if len(word) > 2 {
			words[word] = true
		}
		// fmt.Println("Scanned text:",in.Text())
	}
	
	for word := range words {
		fmt.Println(word)
	}

	if words[query]{
		fmt.Printf("The input contains %q.\n", query)
		return
	}

	fmt.Printf("The input does not contain %q.\n", query)
	// if err := in.Err(); err != nil {
	// 	fmt.Println("ERROR:", err)
	// }
}