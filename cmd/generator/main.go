package main

import (
  "log"
	"inkbottle/internal"
)

func main() {
	if err := internal.Build(); err != nil {
		log.Fatal(err)
	}
	log.Println("build done -> docs/")
}
