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
	wg.Add(6)
	go rotateTest(bot, wg)
	go sireneTest(bot, wg)
	go blinkTest(bot, wg)
	go lineSensorTest(bot, wg)
	go ultrasonicSensorTest(bot, wg)
	go lightSensorTest(bot, wg)
	wg.Wait()

	fmt.Println("------------")
}

// ------------------------------------------------------------

func rotateTest(bot *mbot.Mbot, wg *sync.WaitGroup) {
	defer wg.Done()
	bot.MotorCmd(mbot.RightMotor, 100)
	bot.MotorCmd(mbot.LeftMotor, -100)
	time.Sleep(1500 * time.Millisecond)
	bot.MotorCmd(mbot.RightMotor, -100)
	bot.MotorCmd(mbot.LeftMotor, 100)
	time.Sleep(3000 * time.Millisecond)
	bot.MotorCmd(mbot.RightMotor, 100)
	bot.MotorCmd(mbot.LeftMotor, -100)
	time.Sleep(1500 * time.Millisecond)
	bot.MotorCmd(mbot.RightMotor, 0)
	bot.MotorCmd(mbot.LeftMotor, 0)
}

func sireneTest(bot *mbot.Mbot, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(2000 * time.Millisecond)
	for t := 0; t < 7; t++ {
		bot.BuzzerCmd(6000, 80)
		time.Sleep(200 * time.Millisecond)
		bot.BuzzerCmd(3000, 80)
		time.Sleep(200 * time.Millisecond)
	}
}

func blinkTest(bot *mbot.Mbot, wg *sync.WaitGroup) {
	defer wg.Done()
	for t := 0; t < 12; t++ {
		bot.LedCmd(mbot.LedLeft, 0xa0, 0x00, 0x00)
		bot.LedCmd(mbot.LedRight, 0x00, 0x00, 0xa0)
		time.Sleep(200 * time.Millisecond)
		bot.LedCmd(mbot.LedLeft, 0x00, 0x00, 0xa0)
		bot.LedCmd(mbot.LedRight, 0xa0, 0x00, 0x00)
		time.Sleep(200 * time.Millisecond)
	}
	bot.LedCmd(mbot.LedBoth, 0x00, 0x00, 0x00)
}

func lineSensorTest(bot *mbot.Mbot, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(1000 * time.Millisecond)
	for t := 0; t < 20; t++ {
		right, left := bot.LineSensorCmd()
		if left {
			bot.LedCmd(mbot.LedLeft, 0x00, 0xa0, 0x00)
		} else {
			bot.LedCmd(mbot.LedLeft, 0x00, 0x00, 0x00)
		}
		if right {
			bot.LedCmd(mbot.LedRight, 0x00, 0xa0, 0x00)
		} else {
			bot.LedCmd(mbot.LedRight, 0x00, 0x00, 0x00)
		}
		time.Sleep(200 * time.Millisecond)
	}
	bot.LedCmd(mbot.LedBoth, 0x00, 0x00, 0x00)
}

func ultrasonicSensorTest(bot *mbot.Mbot, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(2000 * time.Millisecond)
	for t := 0; t < 30; t++ {
		val := bot.UltrasonicSensorCmd()
		fmt.Println("ultrasonic:", val)
		time.Sleep(100 * time.Millisecond)
	}
}

func lightSensorTest(bot *mbot.Mbot, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(3000 * time.Millisecond)
	for t := 0; t < 20; t++ {
		val := bot.LigtSensorCmd()
		fmt.Println("light:", val)
		time.Sleep(100 * time.Millisecond)
	}
}
