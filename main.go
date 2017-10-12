package main

// cfr.
//   http://learn.makeblock.com/en/mbot-serial-port-protocol/

import (
	"fmt"
	"log"
	"time"
	"sync"

	"github.com/tarm/serial"
)

// ----------------------------------------------------------

type mbot struct {
	port *serial.Port
	cmd  []byte
	mux  sync.Mutex
}

func makeMbot(portname string) *mbot {
	bot := mbot{}
	c := &serial.Config{
		Name: portname,
		Baud: 9600, //19200,
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
	bot.mux.Lock()
	defer bot.mux.Unlock()
	bot.port.Flush()
	bot.cmd = bot.cmd[0:0]
	bot.port.Close()
}

func (bot *mbot) addCmd(cmd ...byte) {
	bot.mux.Lock()
  bot.cmd = append(bot.cmd, cmd...)
  //bot.cmd = append(bot.cmd, 0x00, 0x00)
	bot.mux.Unlock()
}

func (bot *mbot) sendCmd() {
	if len(bot.cmd) < 1 {
		return
	}
	//fmt.Println("sendCmd:", bot.cmd)
	bot.mux.Lock()
	defer bot.mux.Unlock()
	n, err := bot.port.Write(bot.cmd)
	if err != nil {
		log.Fatal(err)
	}
	if n < len(bot.cmd) {
		log.Fatal(n, "is smaller then ", len(bot.cmd))
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
	if motor == rightMotor {
		speed = -speed
	}
	bot.addCmd(0xff, 0x55, 0x06, 0x60, 0x02, 0x0a, motor,
             byte(speed & 0xff), byte((speed >> 8) & 0xff))
}

// ----------------------------------------------------------

func rotateTest(bot *mbot, wg *sync.WaitGroup) {
	for i := 0; i < 2; i++ {
    bot.motorCmd(rightMotor, 200)
    bot.motorCmd(leftMotor, 200)
    bot.sendCmd()
    time.Sleep(500 * time.Millisecond)
    bot.motorCmd(rightMotor, -200)
    bot.motorCmd(leftMotor, -200)
    bot.sendCmd()
    time.Sleep(500 * time.Millisecond)
  }
  bot.motorCmd(rightMotor, 0)
  bot.motorCmd(leftMotor, 0)
  bot.sendCmd()
  wg.Done()
}

func sireneTest(bot *mbot, wg *sync.WaitGroup) {
  time.Sleep(100 * time.Millisecond)
	for t := 0; t < 8; t++ {
		bot.buzzerCmd(6000, 80)
		time.Sleep(200 * time.Millisecond)
		bot.sendCmd()
		bot.buzzerCmd(3000, 80)
		time.Sleep(200 * time.Millisecond)
		bot.sendCmd()
	}
  wg.Done()
}

func blinkTest(bot *mbot, wg *sync.WaitGroup) {
  time.Sleep(200 * time.Millisecond)
  for t := 0; t < 8; t++ {
  	bot.ledCmd(ledLeft, 0xa0, 0x00, 0x00)
  	bot.ledCmd(ledRight, 0x00, 0x00, 0xa0)
  	time.Sleep(200 * time.Millisecond)
  	bot.sendCmd()
  	bot.ledCmd(ledLeft, 0x00, 0x00, 0xa0)
  	bot.ledCmd(ledRight, 0xa0, 0x00, 0x00)
  	time.Sleep(200 * time.Millisecond)
  	bot.sendCmd()
  }
  bot.ledCmd(ledBoth, 0x00, 0x00, 0x00)
  bot.sendCmd()
  wg.Done()
}

func main() {
	fmt.Println("--- mBot ---")
	bot := makeMbot("COM4")
	defer bot.close()

	wg := sync.WaitGroup{}
	wg.Add(3)
  go rotateTest(bot, &wg)
  go sireneTest(bot, &wg)
  go blinkTest(bot, &wg)
	wg.Wait()
	fmt.Println("------------")
}
