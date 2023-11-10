package main

import (
	"log"

	_ "github.com/amacneil/dbmate/v2/pkg/driver/postgres"
	"github.com/nawaltni/tracker/cmd"
)

func main() {
	if err := cmd.RootCommand().Execute(); err != nil {
		log.Fatal(err)
	}
}
