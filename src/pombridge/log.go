package pombridge
import (
	"log"
	"os"
	"fmt"
)

type PomLogger struct {
	level int
}
var Log PomLogger

const (
	DebugLevel = iota
	InfoLevel = iota
	WarningLevel = iota
	ErrorLevel = iota
	FatalLevel = iota
)

var pLog = log.New(os.Stdout, "pombridge ", log.Ltime)

func (Log *PomLogger) SetLevel(level int) {
	Log.level = level
}

func (Log *PomLogger) D(msg ...interface{}) {
	if (Log.level > DebugLevel) { return }
	pLog.Print("[D] " + fmt.Sprint(msg...))
}

func (Log *PomLogger) I(msg ...interface{}) {
	if (Log.level > InfoLevel) { return }
	pLog.Print("[I] " + fmt.Sprint(msg...))
}

func (Log *PomLogger) W(msg ...interface{}) {
	if (Log.level > WarningLevel) { return }
	pLog.Print("[W] " + fmt.Sprint(msg...))
}

func (Log *PomLogger) E(msg ...interface{}) {
	if (Log.level > ErrorLevel) { return }
	pLog.Print("[E] " + fmt.Sprint(msg...))
}

func (Log *PomLogger) F(msg ...interface{}) {
	pLog.Print("[F] " + fmt.Sprint(msg...))
	os.Exit(-1)
}
