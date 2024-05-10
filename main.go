package main

import (
	"log"
	"os"
)

func main() {
	repoURL := os.Args[1]
	hitFileName := os.Args[2]
	file, err := os.Open(hitFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	// convert output path go repo url

}
