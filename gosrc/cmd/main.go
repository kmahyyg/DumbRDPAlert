//go:build windows

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"rdpalert/embedded"
	"rdpalert/pushsdk"
	"rdpalert/utils"
)

const (
	LOGFILE_NAME  = "rdpalert_running.log"
	CONFJSON_NAME = "rdpalert_pushconf.json"
)

var (
	ErrParamInvalid = errors.New("does not have enough args")
	pushConf        = &pushsdk.PushConfig{}
)

func main() {
	// make sure log file is written to where program located, since CWD is SYSTEM32
	curExecPath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	// separate file and dir
	curWorkPath := filepath.Dir(curExecPath)
	finalLogFilePath := filepath.Join(curWorkPath, LOGFILE_NAME)
	// add logrotate
	logFdPrevinfo, err := os.Stat(finalLogFilePath)
	if err == nil {
		if logFdPrevinfo.Size() > 4194304 {
			// only preserve last 4MiB data
			_ = os.Remove(finalLogFilePath)
		}
	}
	// log file prepare and logger init
	err = utils.InitLogger(finalLogFilePath, "RDPAlert ")
	if err != nil {
		panic(err)
	}
	gLogger, err2 := utils.GetLoggerInstance()
	if err2 != nil {
		panic(err2)
	}
	defer func() { _ = utils.DestoryLoggerInstance() }()
	gLogger.Info("Logging file prepared.")
	gLogger.Debug("Argv: ", os.Args)
	gLogger.Info("Current Version: ", embedded.CurVersionStr)
	// static data ingestion, check to get hostname first
	_, err = os.Hostname()
	if err != nil {
		gLogger.Critical("get hostname:", err)
	}
	gLogger.Info("Current FilePath and Hostname got.")
	// params handling
	if len(os.Args) != 4 {
		printUsage()
		gLogger.Critical("param length check:", ErrParamInvalid)
	}
	gLogger.Info("Params checked.")
	// config load
	curConfPath := filepath.Join(curWorkPath, CONFJSON_NAME)
	confData, err := os.ReadFile(curConfPath)
	if err != nil {
		gLogger.Critical("read conf data:", err)
	}
	gLogger.Info("Config File Opened, Path: ", curConfPath)
	err = json.Unmarshal(confData, pushConf)
	if err != nil {
		gLogger.Critical("json unmarshal - pushconf:", err)
	}
	gLogger.Info("Config File Unmarshal Success.")
	// init pusher, while calling new method:
	// check config logic and if everything is fulfilled
	pusher, err := pushsdk.NewPusher(pushConf)
	if err != nil {
		gLogger.Critical("new pusher: ", err)
	}
	gLogger.Info("Pusher initialized.")
	// build generalized push content
	gpc, err := preparePushContent()
	if err != nil {
		gLogger.Critical("prepare push content:", err)
	}
	gLogger.Info("Push content prepared.")
	pusher.StageGeneralPushContent(gpc)
	gLogger.Info("Push content staged successfully.")
	// transform and send out
	err = pusher.SendPush()
	if err != nil {
		gLogger.Critical("send push content:", err)
	}
	gLogger.Info("Push content sent successfully.")
}

func printUsage() {
	_, _ = fmt.Fprintln(os.Stderr, "Usage: RDPAlarm.exe <Auth Domain> <Auth Username> <Auth IP>.")
}

func preparePushContent() (*pushsdk.GeneralPushContent, error) {
	gLogger, err := utils.GetLoggerInstance()
	if err != nil {
		return nil, err
	}
	cIPs, err := utils.GetLocalIP()
	if err != nil {
		panic(err)
	}
	if cIPs == nil {
		gLogger.Critical("get local ip:", utils.ErrCannotGetLocalIP)
	}
	gLogger.Info("Local IP got:", cIPs)
	args := os.Args
	hostname, _ := os.Hostname()
	authDomain := func() string {
		if len(args[1]) == 0 {
			return "localhost"
		} else {
			return args[1]
		}
	}()
	notiTitle := "RDP Login - Success"
	notiBody := fmt.Sprintf("From: %s - %s\\%s\nHost: %s, Host IPs: %s \n",
		args[3], authDomain, args[2], hostname, cIPs)
	notiShort := fmt.Sprintf("User %s from %s", args[2], args[3])
	gpc := &pushsdk.GeneralPushContent{
		Title:       notiTitle,
		ShortTitle:  notiShort,
		Description: notiBody,
	}
	return gpc, nil
}
