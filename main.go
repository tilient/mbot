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

type mbot struct {
	port *serial.Port
	cmd  []byte
}

func makeMbot(portname string) *mbot {
	bot := mbot{}
	c := &serial.Config{
		Name: portname,
		Baud: 57600,
  	ReadTimeout: 500 * time.Millisecond}
	p, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	bot.port = p
	bot.cmd = make([]byte, 0, 256)
	return &bot
}

func (bot *mbot)close() {
	bot.cmd = bot.cmd[0:0]
	bot.port.Close()
}

func (bot *mbot) addCmd(cmd ...byte) {
  bot.cmd = append(bot.cmd, cmd...)
}

func (bot *mbot) sendCmd() {
	_, err := bot.port.Write(bot.cmd)
	if err != nil {
		log.Fatal(err)
	}
	err = bot.port.Flush()
	if err != nil {
		log.Fatal(err)
	}
	bot.cmd = bot.cmd[0:0]
	time.Sleep(10 * time.Millisecond)
}

// ----------------------------------------------------------

const (
	ledLeft  = 0x01
  ledRight = 0x02
  ledBoth  = 0x00
)

func (bot *mbot)ledCmd(led byte,
                r byte, g byte, b byte) {
	bot.addCmd(0xff, 0x55, 0x09, 0x00, 0x02, 0x08, 0x07, 0x02,
             led, r, g, b)
}

// ----------------------------------------------------------

func (bot *mbot)buzzerCmd(tone uint16, beat uint16) {
	bot.addCmd(0xff, 0x55, 0x07, 0x00, 0x02, 0x22,
             byte(tone & 0xff), byte((tone >> 8) & 0xff),
	           byte(beat & 0xff), byte((beat >> 8) & 0xff))
}

// ----------------------------------------------------------

const (
	leftMotor = 0x09
	rightMotor = 0x0a
)

func (bot *mbot)motorCmd(motor byte, speed int16) {
	bot.addCmd(0xff, 0x55, 0x06, 0x60, 0x02, 0x0a, motor,
             byte(speed & 0xff), byte((speed >> 8) & 0xff))
}

// ----------------------------------------------------------

func main() {
	fmt.Println("--- mBot ---")
	bot := makeMbot("COM4")
	defer bot.close()

	//bot.motorCmd(leftMotor, -50)
	//bot.motorCmd(rightMotor, 50)
  //bot.sendCmd()
  //time.Sleep(500 * time.Millisecond)
	//bot.motorCmd(leftMotor, 0)
	//bot.motorCmd(rightMotor, 0)
  //bot.sendCmd()
	for t := 0; t < 16; t++ {
		bot.buzzerCmd(6000, 80)
		bot.ledCmd(ledLeft, 0xa0, 0x00, 0x00)
		bot.ledCmd(ledRight, 0x00, 0x00, 0xa0)
		time.Sleep(200 * time.Millisecond)
		bot.sendCmd()
		bot.buzzerCmd(3000, 80)
		bot.ledCmd(ledLeft, 0x00, 0x00, 0xa0)
		bot.ledCmd(ledRight, 0xa0, 0x00, 0x00)
		time.Sleep(200 * time.Millisecond)
		bot.sendCmd()
	}
	bot.ledCmd(ledBoth, 0x00, 0x00, 0x00)
	bot.sendCmd()

	fmt.Println("------------")
}
