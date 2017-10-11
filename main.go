package main

import (
	"fmt"
	"log"
	"time"

	"github.com/tarm/serial"
)

func sendCodes(s *serial.Port, str string) {
	_, err := s.Write([]byte(str))
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(50 * time.Millisecond)
}

func main() {
	fmt.Println("--- mBot ---")
	c := &serial.Config{Name: "/dev/ttyS4", Baud: 115200}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()

	for t := 0; t < 64; t++ {
		sendCodes(s, "\xff\x55\x09\x00\x02\x08\x07\x02\x00\xaa\x00\x00")
		//sendCodes(s, "\xff\x55\x07\x00\x02\x22\x7b\x00\xfa\x00")
		time.Sleep(150 * time.Millisecond)
		sendCodes(s, "\xff\x55\x09\x00\x02\x08\x07\x02\x00\x00\x00\xaa")
		//sendCodes(s, "\xff\x55\x07\x00\x02\x22\x9b\x01\xfa\x00")
		time.Sleep(150 * time.Millisecond)
	}
	sendCodes(s, "\xff\x55\x09\x00\x02\x08\x07\x02\x00\x00\x00\x00")

	fmt.Println("------------")
}
