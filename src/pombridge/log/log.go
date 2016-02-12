package log
import (
	"log"
	"os"
	"fmt"
)

var level = 0

const (
	DebugLevel = iota
	InfoLevel = iota
	WarningLevel = iota
	ErrorLevel = iota
	FatalLevel = iota
)

var pLog = log.New(os.Stdout, "pombridge ", log.Ltime)

func SetLevel(l int) {
	level = l
}

func D(msg ...interface{}) {
	if (level > DebugLevel) { return }
	pLog.Print("[D] " + fmt.Sprint(msg...))
}

func I(msg ...interface{}) {
	if (level > InfoLevel) { return }
	pLog.Print("[I] " + fmt.Sprint(msg...))
}

func W(msg ...interface{}) {
	if (level > WarningLevel) { return }
	pLog.Print("[W] " + fmt.Sprint(msg...))
}

func E(msg ...interface{}) {
	if (level > ErrorLevel) { return }
	pLog.Print("[E] " + fmt.Sprint(msg...))
}

func F(msg ...interface{}) {
	pLog.Print("[F] " + fmt.Sprint(msg...))
	os.Exit(-1)
}
