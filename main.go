package main

// cfr.
//   http://learn.makeblock.com/en/...
//     ...mbot-serial-port-protocol/
//
/* to compile for windows
export GOOS=windows
export GOARCH=amd64
export CC=x86_64-w64-mingw32-gcc
*/

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/tarm/serial"
)

// ------------------------------------------------

type mbot struct {
	port *serial.Port
	mux  sync.Mutex
}

func makeMbot(portname string) *mbot {
	bot := mbot{}
	c := &serial.Config{
		Name:        portname,
		Baud:        57600,
		ReadTimeout: 500 * time.Millisecond}
	p, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	bot.port = p
	return &bot
}

func (bot *mbot) close() {
	bot.port.Close()
}

func (bot *mbot) cmd(cmd ...byte) []byte {
	bot.mux.Lock()
	defer bot.mux.Unlock()

	n, err := bot.port.Write(cmd)
	if err != nil {
		log.Fatal(err)
	}
	if n < len(cmd) {
		log.Fatal(n, "is smaller then ", len(cmd))
	}
	err = bot.port.Flush()
	if err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, 32)
	n, err = bot.port.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	return buf[:n]
}

// ------------------------------------------------

const (
	ledLeft  = 0x01
	ledRight = 0x02
	ledBoth  = 0x00
)

func (bot *mbot) ledCmd(led byte,
	r byte, g byte, b byte) {
	bot.cmd(0xff, 0x55, 0x09, 0x00, 0x02, 0x08,
		0x07, 0x02, led, r, g, b)
}

// ------------------------------------------------

func (bot *mbot) buzzerCmd(tone uint16, beat uint16) {
	bot.cmd(0xff, 0x55, 0x07, 0x00, 0x02, 0x22,
		byte(tone&0xff), byte((tone>>8)&0xff),
		byte(beat&0xff), byte((beat>>8)&0xff))
}

// ------------------------------------------------

const (
	leftMotor  = 0x09
	rightMotor = 0x0a
)

func (bot *mbot) motorCmd(motor byte, speed int16) {
	if motor == rightMotor {
		speed = -speed
	}
	bot.cmd(0xff, 0x55, 0x06, 0x60, 0x02, 0x0a,
		motor,
		byte(speed&0xff), byte((speed>>8)&0xff))
}

// ------------------------------------------------

func (bot *mbot) lineSensorCmd() (bool, bool) {
	res := bot.cmd(0xff, 0x55, 0x04, 0x60, 0x01, 0x11, 0x02)
	if len(res) < 10 {
		log.Fatal("wrong line sensor result")
	}
	if (res[6] == 0x00) && (res[7] == 0x00) {
		return false, false
	}
	if (res[6] == 0x40) && (res[7] == 0x40) {
		return true, true
	}
	if (res[6] == 0x00) && (res[7] == 0x40) {
		return true, false
	}
	if (res[6] == 0x80) && (res[7] == 0x3f) {
		return false, true
	}
	log.Fatal("wrong line sensor result")
	return false, false
}

// ------------------------------------------------

func (bot *mbot) ultrasonicSensorCmd() uint32 {
	res := bot.cmd(0xff, 0x55, 0x04, 0x00, 0x01, 0x01, 0x03)
	if len(res) < 10 {
		log.Fatal("wrong ultrasonic sensor result")
	}
	return uint32(res[7]-64)*256 + uint32(res[6])
}

// ------------------------------------------------

func rotateTest(bot *mbot, wg *sync.WaitGroup) {
	bot.motorCmd(rightMotor, 100)
	bot.motorCmd(leftMotor, -100)
	time.Sleep(1500 * time.Millisecond)
	bot.motorCmd(rightMotor, -100)
	bot.motorCmd(leftMotor, 100)
	time.Sleep(3000 * time.Millisecond)
	bot.motorCmd(rightMotor, 100)
	bot.motorCmd(leftMotor, -100)
	time.Sleep(1500 * time.Millisecond)
	bot.motorCmd(rightMotor, 0)
	bot.motorCmd(leftMotor, 0)
	wg.Done()
}

func sireneTest(bot *mbot, wg *sync.WaitGroup) {
	time.Sleep(2000 * time.Millisecond)
	for t := 0; t < 4; t++ {
		bot.buzzerCmd(6000, 80)
		time.Sleep(200 * time.Millisecond)
		bot.buzzerCmd(3000, 80)
		time.Sleep(200 * time.Millisecond)
	}
	wg.Done()
}

func blinkTest(bot *mbot, wg *sync.WaitGroup) {
	time.Sleep(200 * time.Millisecond)
	for t := 0; t < 12; t++ {
		bot.ledCmd(ledLeft, 0xa0, 0x00, 0x00)
		bot.ledCmd(ledRight, 0x00, 0x00, 0xa0)
		time.Sleep(200 * time.Millisecond)
		bot.ledCmd(ledLeft, 0x00, 0x00, 0xa0)
		bot.ledCmd(ledRight, 0xa0, 0x00, 0x00)
		time.Sleep(200 * time.Millisecond)
	}
	bot.ledCmd(ledBoth, 0x00, 0x00, 0x00)
	wg.Done()
}

func lineSensorTest(bot *mbot, wg *sync.WaitGroup) {
	time.Sleep(1000 * time.Millisecond)
	for t := 0; t < 25; t++ {
		right, left := bot.lineSensorCmd()
		if left {
			bot.ledCmd(ledLeft, 0x00, 0xa0, 0x00)
		} else {
			bot.ledCmd(ledLeft, 0x00, 0x00, 0x00)
		}
		if right {
			bot.ledCmd(ledRight, 0x00, 0xa0, 0x00)
		} else {
			bot.ledCmd(ledRight, 0x00, 0x00, 0x00)
		}
		time.Sleep(200 * time.Millisecond)
	}
	bot.ledCmd(ledBoth, 0x00, 0x00, 0x00)
	wg.Done()
}

func ultrasonicSensorTest(bot *mbot, wg *sync.WaitGroup) {
	time.Sleep(1000 * time.Millisecond)
	for t := 0; t < 50; t++ {
		val := bot.ultrasonicSensorCmd()
		fmt.Println("ultrasonic:", val)
		time.Sleep(100 * time.Millisecond)
	}
	wg.Done()
}

func main() {
	fmt.Println("--- mBot ---")
	bot := makeMbot("COM4")
	defer bot.close()

	wg := sync.WaitGroup{}
	wg.Add(5)
	go rotateTest(bot, &wg)
	go sireneTest(bot, &wg)
	go blinkTest(bot, &wg)
	go lineSensorTest(bot, &wg)
	go ultrasonicSensorTest(bot, &wg)
	wg.Wait()
	fmt.Println("------------")
}
