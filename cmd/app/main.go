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
	if err != nil {
		log.Fatal(fmt.Sprintf("cannot open input file '%s'", filename))
	}
	defer file.Close()

	s, err := service.New(file)
	if err != nil {
		log.Fatal(err)
	}
	if err := s.Run(); err != nil {
		log.Fatal("cannot work")
	}
}
