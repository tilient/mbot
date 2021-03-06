package mbot

import (
	"encoding/binary"
	"log"
	"math"
	"sync"
	"time"

	"github.com/tarm/serial"
)

// --------------------------------------------------------

type Mbot struct {
	port *serial.Port
	mux  sync.Mutex
}

func MakeMbot(portname string) *Mbot {
	p, err := serial.OpenPort(&serial.Config{
		Name:        portname,
		Baud:        57600,
		ReadTimeout: 500 * time.Millisecond})
	if err != nil {
		log.Fatal(err)
	}
	return &Mbot{port: p}
}

func (bot *Mbot) Close() {
	bot.mux.Lock()
	defer bot.mux.Unlock()

	bot.port.Close()
}

func (bot *Mbot) cmd(cmd ...byte) []byte {
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

// --------------------------------------------------------
// Led Command
// --------------------------------------------------------

const (
	LedLeft  = 0x01
	LedRight = 0x02
	LedBoth  = 0x00
)

func (bot *Mbot) LedCmd(led byte, r byte, g byte, b byte) {
	bot.cmd(0xff, 0x55, 0x09, 0x00, 0x02, 0x08, 0x07, 0x02,
		led, r, g, b)
}

// --------------------------------------------------------
// Buzzer Command
// --------------------------------------------------------

func (bot *Mbot) BuzzerCmd(tone uint16, beat uint16) {
	bot.cmd(0xff, 0x55, 0x07, 0x00, 0x02, 0x22,
		byte(tone&0xff), byte((tone>>8)&0xff),
		byte(beat&0xff), byte((beat>>8)&0xff))
}

// --------------------------------------------------------
// Motor Command
// --------------------------------------------------------

const (
	LeftMotor  = 0x09
	RightMotor = 0x0a
)

func (bot *Mbot) MotorCmd(motor byte, speed int16) {
	if motor == RightMotor {
		speed = -speed
	}
	bot.cmd(0xff, 0x55, 0x06, 0x60, 0x02, 0x0a,
		motor, byte(speed&0xff), byte((speed>>8)&0xff))
}

// --------------------------------------------------------
// Line Sensor Command
// --------------------------------------------------------

func (bot *Mbot) LineSensorCmd() (bool, bool) {
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

// --------------------------------------------------------
// Ultrasonic Sensor Command
// --------------------------------------------------------

func (bot *Mbot) UltrasonicSensorCmd() float32 {
	res := bot.cmd(0xff, 0x55, 0x04, 0x00, 0x01, 0x01, 0x03)
	if len(res) < 10 {
		log.Fatal("wrong ultrasonic sensor result")
	}
	return 1.27 * float32frombytes(res[4:8])
}

// --------------------------------------------------------
// Light Sensor Command
// --------------------------------------------------------

func (bot *Mbot) LigtSensorCmd() float32 {
	res := bot.cmd(0xff, 0x55, 0x04, 0x05, 0x01, 0x03, 0x03)
	if len(res) < 10 {
		log.Fatal("wrong light sensor result")
	}
	return float32frombytes(res[4:8])
}

// --------------------------------------------------------
// Tools
// --------------------------------------------------------

func float32frombytes(bytes []byte) float32 {
	return math.Float32frombits(
		binary.LittleEndian.Uint32(bytes))
}
