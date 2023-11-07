package main

import (
	"log"

	"github.com/nawaltni/tracker/cmd"
)

func main() {
	if err := cmd.RootCommand().Execute(); err != nil {
		log.Fatal(err)
	}

}
