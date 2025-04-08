package utils

import (
	"errors"
	"io"
	"log"
	"os"
	"sync/atomic"
)

var (
	currentLogger           = new(dumbLogger)
	ErrLoggerNotInitialized = errors.New("logger not initialized")
)

const (
	debugEnv = "IS_IN_DEBUG"
)

func InitLogger(outputFilePath string, logPrefix string) error {
	currentLogger.fdpath = outputFilePath
	currentLogger.Init(logPrefix)
	currentLogger.b.Store(true)
	return nil
}

func GetLoggerInstance() (*dumbLogger, error) {
	if currentLogger.b.Load() == true {
		return currentLogger, nil
	} else {
		return nil, ErrLoggerNotInitialized
	}
}

func DestoryLoggerInstance() error {
	err := currentLogger.dispose()
	if err != nil {
		return err
	}
	currentLogger = new(dumbLogger)
	return nil
}

type dumbLogger struct {
	fdpath string
	l      *log.Logger
	b      atomic.Bool
	f      *os.File
}

func (myl *dumbLogger) Init(lPrefix string) {
	var err error
	myl.f, err = os.OpenFile(myl.fdpath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	myl.l = log.Default()
	myl.l.SetFlags(log.LstdFlags | log.Lmsgprefix)
	myl.l.SetPrefix(lPrefix)
	if _, exists := os.LookupEnv(debugEnv); exists {
		logmulti_opt := io.MultiWriter(os.Stderr, myl.f)
		myl.l.SetOutput(logmulti_opt)
		return
	}
	myl.l.SetOutput(myl.f)
}

func (myl *dumbLogger) dispose() error {
	if myl.b.Load() == false {
		return ErrLoggerNotInitialized
	}
	myl.b.Store(false)
	_ = myl.f.Sync()
	_ = myl.f.Close()
	return nil
}

func (myl *dumbLogger) Info(v ...any) {
	myl.l.Println("Info:", v)
}

func (myl *dumbLogger) Error(v ...any) {
	myl.l.Println("Error:", v)
}

func (myl *dumbLogger) Warn(v ...any) {
	myl.l.Println("Warning:", v)
}

func (myl *dumbLogger) Critical(v ...any) {
	myl.l.Fatalln("CRITICAL:", v)
}

func (myl *dumbLogger) Debug(v ...any) {
	myl.l.Println("Debug:", v)
}
