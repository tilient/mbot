package main

/*
to cross-compile for windows under linux:
  export GOOS=windows
  export GOARCH=amd64
  export CC=x86_64-w64-mingw32-gcc
  go build test01.go
*/

import (
	"fmt"
	"sync"
	"time"

	"github.com/tilient/mbot"
)

// ------------------------------------------------------------

func main() {
	fmt.Println("--- mbot ---")

	bot := mbot.MakeMbot("COM5")
	defer bot.Close()

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go blink(bot, wg)
	go drive(bot, wg)
	wg.Wait()

	fmt.Println("------------")
}

// ------------------------------------------------------------

func drive(bot *mbot.Mbot, wg *sync.WaitGroup) {
	defer wg.Done()
	for t := 0; t < 6000; t++ {
		bot.MotorCmd(mbot.RightMotor, 80)
		bot.MotorCmd(mbot.LeftMotor, 80)
		time.Sleep(5 * time.Millisecond)
		right, left := bot.LineSensorCmd()
		for (!right) || (!left) {
			bot.MotorCmd(mbot.RightMotor, -50)
			bot.MotorCmd(mbot.LeftMotor, -200)
			time.Sleep(120 * time.Millisecond)
			right, left = bot.LineSensorCmd()
		}
		for bot.UltrasonicSensorCmd() < 20.0 {
			bot.MotorCmd(mbot.LeftMotor, -100)
			time.Sleep(5 * time.Millisecond)
		}
	}
	bot.MotorCmd(mbot.RightMotor, 0)
	bot.MotorCmd(mbot.LeftMotor, 0)
}

func blink(bot *mbot.Mbot, wg *sync.WaitGroup) {
	defer wg.Done()
	for t := 0; t < 600; t++ {
		bot.LedCmd(mbot.LedLeft, 0xa0, 0x00, 0x00)
		bot.LedCmd(mbot.LedRight, 0x00, 0x00, 0xa0)
		time.Sleep(250 * time.Millisecond)
		bot.LedCmd(mbot.LedLeft, 0x00, 0x00, 0xa0)
		bot.LedCmd(mbot.LedRight, 0xa0, 0x00, 0x00)
		time.Sleep(250 * time.Millisecond)
	}
	bot.LedCmd(mbot.LedBoth, 0x00, 0x00, 0x00)
}
