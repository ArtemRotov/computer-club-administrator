package main

import (
	"github.com/ArtemRotov/computer-club-manager/internal/app"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("missing input argument (file)")
	}
	filename := os.Args[1]

	if err := app.Run(filename); err != nil {
		log.Fatal(err)
	}
}
