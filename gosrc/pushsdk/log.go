package pushsdk

import (
	"io"
	"log"
	"os"
)

type DumbLogger struct {
	l *log.Logger
}

func (myl *DumbLogger) Init(logfileFd *os.File, lPrefix string) {
	mwr := io.MultiWriter(os.Stdout, logfileFd)
	myl.l = log.Default()
	myl.l.SetFlags(log.LstdFlags | log.Lshortfile | log.LUTC | log.Lmsgprefix)
	myl.l.SetPrefix(lPrefix)
	myl.l.SetOutput(mwr)
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
