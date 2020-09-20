package main

import (
	"fmt"
	"os"
	"strings"
)

const corpus = "" +
	"lazy cat jumps again and again and again"

func main (){
	words := strings.Fields(corpus)
	query := os.Args[1:]

	for _,q := range query{
		for i,w := range words {
			if strings.ToLower(q) == strings.ToLower(w) {
				fmt.Printf("Found word: %q at position #%-2d \n",w,(i + 1))
				break
			}
		}
	}
}