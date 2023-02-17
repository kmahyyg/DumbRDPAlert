package pushsdk

import (
	"log"
	"os"
)

type DumbLogger struct {
	l *log.Logger
}

func (myl *DumbLogger) Init(logfileFd *os.File, lPrefix string) {
	myl.l = log.Default()
	myl.l.SetFlags(log.LstdFlags | log.Lmsgprefix)
	myl.l.SetPrefix(lPrefix)
	myl.l.SetOutput(logfileFd)
}

func (myl DumbLogger) Info(v ...any) {
	myl.l.Println("Info:", v)
}

func (myl DumbLogger) Error(v ...any) {
	myl.l.Println("Error:", v)
}

func (myl DumbLogger) Warn(v ...any) {
	myl.l.Println("Warning:", v)
}

func (myl DumbLogger) Critical(v ...any) {
	myl.l.Fatalln("CRITICAL:", v)
}

func (myl DumbLogger) Debug(v ...any) {
	myl.l.Println("Debug:", v)
}
