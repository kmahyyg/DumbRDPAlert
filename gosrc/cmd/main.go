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
	netif "tailscale.com/net/interfaces"
)

var (
	ErrParamInvalid     = errors.New("does not have enough args")
	ErrCannotGetLocalIP = errors.New("cannot get local ip")
	pushConf            = &pushsdk.PushConfig{}
	gLogger             = &pushsdk.DumbLogger{}
)

func main() {
	// add logrotate
	logFdPrevinfo, err := os.Stat("rdpalarm-running.log")
	if err == nil {
		if logFdPrevinfo.Size() > 1073741824 {
			_ = os.Remove("rdpalarm-running.log")
		}
	}
	// log file
	logFd, err := os.OpenFile("rdpalarm-running.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
	curProg, err := os.Executable()
	if err != nil {
		gLogger.Critical("get exepath:", err)
	}
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
	curProgDir, _ := filepath.Split(curProg)
	curConfPath := filepath.Join(curProgDir, "rdpalert_pushconf.json")
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
	// pusher
	pusher, err := pushsdk.NewPusher(pushConf, gLogger)
	if err != nil {
		gLogger.Critical("new pusher:", err)
	}
	gLogger.Info("Pusher initalized.")
	// prepare content
	postBodys, err := PreparePushContents(curHostname, os.Args)
	if err != nil {
		gLogger.Critical("prepare contents:", err)
	}
	gLogger.Info("Push Contents prepared,")
	err = pusher.SetContents(postBodys)
	if err != nil {
		gLogger.Critical("post body set check:", err)
	}
	gLogger.Info("Push Contents checked.")
	err = pusher.SendPushRequests()
	if err != nil {
		gLogger.Critical("send push req to serv:", err)
	}
	gLogger.Info("Push Request sent to server. now exit. done.")
}

func printUsage() {
	_, _ = fmt.Fprintln(os.Stderr, "Usage: RDPAlarm.exe <Auth Domain> <Auth Username> <Auth IP>.")
}

func PreparePushContents(hostname string, args []string) ([]*pushsdk.PushContent, error) {
	cIPs := GetLocalIP()
	if cIPs == nil {
		gLogger.Critical("get local ip:", ErrCannotGetLocalIP)
	}
	gLogger.Info("Local IP got:", cIPs)
	pCont := make([]*pushsdk.PushContent, 0)
	authDomain := func() string {
		if len(args[1]) == 0 {
			return "localhost"
		} else {
			return args[1]
		}
	}()
	notiTitle := "RDP Login - Success"
	notiBody := fmt.Sprintf("From: %s - %s\\%s\nHost: %s, IPs: %s \n",
		args[3], authDomain, args[2], hostname, cIPs)
	for _, v := range pushConf.DeviceKeys {
		data := &pushsdk.PushContent{
			Title:     notiTitle,
			Body:      notiBody,
			DeviceKey: v,
		}
		data.Init()
		data.Level = pushConf.NotificationLevel
		data.Group = "Security_RDPAlert"
		data.Copy = hostname
		pCont = append(pCont, data)
	}
	return pCont, nil
}

// GetLocalIP returns the non loopback local IP of the host
func GetLocalIP() []string {
	dIf, _, err := netif.LocalAddresses()
	if err != nil {
		gLogger.Critical("get local addr:", err)
		return nil
	}
	if len(dIf) == 0 {
		return nil
	}
	res := make([]string, 0)
	for _, v := range dIf {
		res = append(res, v.String())
	}
	return res
}
