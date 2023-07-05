package main

import (
	"fmt"
	"time"
)

const (
	TimeLayout = "15:04"
)

func main() {
	str1, _ := time.Parse(TimeLayout, "12:33")
	str2, _ := time.Parse(TimeLayout, "15:52")
	fmt.Println(str2.Sub(str1))

	//time, err := time.Parse(TimeLayout, str)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(time.Format(TimeLayout))
}
