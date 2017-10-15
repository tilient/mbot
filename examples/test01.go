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
	"math"
	"sync"
	"time"

	"github.com/tilient/mbot"
)

// ------------------------------------------------

func rotateTest(bot *Mbot, wg *sync.WaitGroup) {
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

func sireneTest(bot *Mbot, wg *sync.WaitGroup) {
	time.Sleep(2000 * time.Millisecond)
	for t := 0; t < 4; t++ {
		bot.buzzerCmd(6000, 80)
		time.Sleep(200 * time.Millisecond)
		bot.buzzerCmd(3000, 80)
		time.Sleep(200 * time.Millisecond)
	}
	wg.Done()
}

func blinkTest(bot *Mbot, wg *sync.WaitGroup) {
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

func lineSensorTest(bot *Mbot, wg *sync.WaitGroup) {
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

func ultrasonicSensorTest(bot *Mbot, wg *sync.WaitGroup) {
	time.Sleep(1000 * time.Millisecond)
	for t := 0; t < 50; t++ {
		val := bot.ultrasonicSensorCmd()
		fmt.Println("ultrasonic:", val)
		time.Sleep(100 * time.Millisecond)
	}
	wg.Done()
}

func main() {
	fmt.Println("--- Mbot ---")
	bot := makeMbot("COM5")
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
