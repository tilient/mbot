package main

// cfr.
//   http://learn.makeblock.com/en/mbot-serial-port-protocol/

import (
	"fmt"
	"log"
	"time"

	"github.com/tarm/serial"
)

// ----------------------------------------------------------

func sendCmd(s *serial.Port, cmd []byte){
	_, err := s.Write(cmd)
	if err != nil {
		log.Fatal(err)
	}
	err = s.Flush()
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(50 * time.Millisecond)
}

// ----------------------------------------------------------

const (
	ledLeft  = 0x01
  ledRight = 0x02
  ledBoth  = 0x00
)

func sendLedCmd(s *serial.Port, position byte,
                r byte, g byte, b byte) {
	sendCmd(s, []byte{0xff, 0x55, 0x09, 0x00, 0x02, 0x08,
	                  0x07, 0x02, position, r, g, b})
}

// ----------------------------------------------------------

func sendBuzzerCmd(s *serial.Port,
                   tone uint16, beat uint16) {
	toneLow := byte(tone & 0xff)
	toneHigh := byte((tone >> 8) & 0xff)
	beatLow := byte(beat & 0xff)
	beatHigh := byte((beat >> 8) & 0xff)
	sendCmd(s, []byte{0xff, 0x55, 0x07, 0x00, 0x02, 0x22,
	                  toneLow, toneHigh, beatLow, beatHigh})
}

// ----------------------------------------------------------

func main() {
	fmt.Println("--- mBot ---")
	c := &serial.Config{Name: "COM4", Baud: 57600,
  	ReadTimeout: 500 * time.Millisecond}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()

	for t := 0; t < 8; t++ {
	  sendCmd(s, []byte{
		  0xff, 0x55, 0x09, 0x00, 0x02, 0x08,
	    0x07, 0x02, 0x01, 0xa0, 0x00, 0x00,
		  0xff, 0x55, 0x09, 0x00, 0x02, 0x08,
	    0x07, 0x02, 0x02, 0x00, 0x00, 0xa0 })
		//sendBuzzerCmd(s, 6000, 4)
		//sendLedCmd(s, ledLeft, 0xa0, 0x00, 0x00)
		//sendLedCmd(s, ledRight, 0x00, 0x00, 0xa0)
		time.Sleep(25 * time.Millisecond)
	  sendCmd(s, []byte{
		  0xff, 0x55, 0x09, 0x00, 0x02, 0x08,
	    0x07, 0x02, 0x01, 0x00, 0x00, 0xa0,
		  0xff, 0x55, 0x09, 0x00, 0x02, 0x08,
	    0x07, 0x02, 0x02, 0xa0, 0x00, 0x00 })
		//sendLedCmd(s, ledLeft, 0x00, 0x00, 0xa0)
		//sendLedCmd(s, ledRight, 0xa0, 0x00, 0x00)
		time.Sleep(25 * time.Millisecond)
	}
	sendLedCmd(s, ledBoth, 0x00, 0x00, 0x00)

	fmt.Println("------------")
}
