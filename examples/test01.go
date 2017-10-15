package main

/* cfr.
   http://learn.makeblock.com/en/...
     ...mbot-serial-port-protocol/

to compile for windows

export GOOS=windows
export GOARCH=amd64
export CC=x86_64-w64-mingw32-gcc
*/

import (
	"fmt"
	"sync"
	"time"

	"github.com/tilient/mbot"
)

// ------------------------------------------------

func rotateTest(bot *mbot.Mbot, wg *sync.WaitGroup) {
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
	wg.Done()
}

func sireneTest(bot *mbot.Mbot, wg *sync.WaitGroup) {
	time.Sleep(2000 * time.Millisecond)
	for t := 0; t < 4; t++ {
		bot.BuzzerCmd(6000, 80)
		time.Sleep(200 * time.Millisecond)
		bot.BuzzerCmd(3000, 80)
		time.Sleep(200 * time.Millisecond)
	}
	wg.Done()
}

func blinkTest(bot *mbot.Mbot, wg *sync.WaitGroup) {
	time.Sleep(200 * time.Millisecond)
	for t := 0; t < 12; t++ {
		bot.LedCmd(mbot.LedLeft, 0xa0, 0x00, 0x00)
		bot.LedCmd(mbot.LedRight, 0x00, 0x00, 0xa0)
		time.Sleep(200 * time.Millisecond)
		bot.LedCmd(mbot.LedLeft, 0x00, 0x00, 0xa0)
		bot.LedCmd(mbot.LedRight, 0xa0, 0x00, 0x00)
		time.Sleep(200 * time.Millisecond)
	}
	bot.LedCmd(mbot.LedBoth, 0x00, 0x00, 0x00)
	wg.Done()
}

func lineSensorTest(bot *mbot.Mbot, wg *sync.WaitGroup) {
	time.Sleep(1000 * time.Millisecond)
	for t := 0; t < 25; t++ {
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
	wg.Done()
}

func ultrasonicSensorTest(bot *mbot.Mbot, wg *sync.WaitGroup) {
	time.Sleep(1000 * time.Millisecond)
	for t := 0; t < 50; t++ {
		val := bot.UltrasonicSensorCmd()
		fmt.Println("ultrasonic:", val)
		time.Sleep(100 * time.Millisecond)
	}
	wg.Done()
}

func lightSensorTest(bot *mbot.Mbot, wg *sync.WaitGroup) {
	time.Sleep(1000 * time.Millisecond)
	for t := 0; t < 50; t++ {
		val := bot.LigtSensorCmd()
		fmt.Println("light:", val)
		time.Sleep(100 * time.Millisecond)
	}
	wg.Done()
}

func main() {
	fmt.Println("--- mbot ---")
	bot := mbot.MakeMbot("COM5")
	defer bot.Close()

	wg := sync.WaitGroup{}
	wg.Add(4)
	//go rotateTest(bot, &wg)
	//go sireneTest(bot, &wg)
	go blinkTest(bot, &wg)
	go lineSensorTest(bot, &wg)
	go ultrasonicSensorTest(bot, &wg)
	go lightSensorTest(bot, &wg)
	wg.Wait()
	fmt.Println("------------")
}
