package main

import (
	"fmt"
	"github.com/ArtemRotov/computer-club-administrator/internal/service"
	"log"
	"os"
)

func main() {
	filename := os.Args[1]

	file, err := os.Open(filename)
	defer file.Close()

	if err != nil {
		log.Fatal(fmt.Sprintf("cannot open file '%s'", filename))
	}

	_, _ = service.New(file)
}
