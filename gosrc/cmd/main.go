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
	netif "tailscale.com/net/netmon"
)

const (
	LOGFILE_NAME  = "rdpalarm-running.log"
	CONFJSON_NAME = "rdpalert_pushconf.json"
)

var (
	ErrParamInvalid     = errors.New("does not have enough args")
	ErrCannotGetLocalIP = errors.New("cannot get local ip")
	pushConf            = &pushsdk.PushConfig{}
	gLogger             = &pushsdk.DumbLogger{}
)

func main() {
	// make sure log file is written to where program located, since CWD is SYSTEM32
	curExecPath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	// seperate file and dir
	curWorkPath := filepath.Dir(curExecPath)
	finalLogFilePath := filepath.Join(curWorkPath, LOGFILE_NAME)
	// add logrotate
	logFdPrevinfo, err := os.Stat(finalLogFilePath)
	if err == nil {
		if logFdPrevinfo.Size() > 4194304 {
			_ = os.Remove(finalLogFilePath)
		}
	}
	// log file
	logFd, err := os.OpenFile(finalLogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	gLogger.Init(logFd, "RDPAlarm ")
	defer func() {
		_ = logFd.Sync()
		_ = logFd.Close()
	}()
	gLogger.Info("Logging file prepared.")
	gLogger.Debug("Argv: ", os.Args)
	gLogger.Info("Current Version: ", embedded.CurVersionStr)
	// static data
	curHostname, err := os.Hostname()
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
	// config
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

}

func printUsage() {
	_, _ = fmt.Fprintln(os.Stderr, "Usage: RDPAlarm.exe <Auth Domain> <Auth Username> <Auth IP>.")
}

// GetLocalIP returns the non loopback local IP of the host
func GetLocalIP() ([]string, error) {
	dIf, _, err := netif.LocalAddresses()
	if err != nil {
		gLogger.Critical("get local addr:", err)
		return nil, err
	}
	if len(dIf) == 0 {
		return nil, ErrCannotGetLocalIP
	}
	res := make([]string, 0)
	for _, v := range dIf {
		res = append(res, v.String())
	}
	return res, nil
}
