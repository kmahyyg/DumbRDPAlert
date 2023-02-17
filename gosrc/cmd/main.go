//go:build windows

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"rdpalert/embedded"
	"rdpalert/pushsdk"
)

var (
	ErrParamInvalid     = errors.New("does not have enough args")
	ErrCannotGetLocalIP = errors.New("cannot get local ip")
	pushConf            = &pushsdk.PushConfig{}
	gLogger             = &pushsdk.DumbLogger{}
)

func main() {
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
			return "."
		} else {
			return args[1]
		}
	}()
	notiTitle := "RDP Login - Success"
	notiBody := fmt.Sprintf("Login from %s As %s\\%s .\nCurrent host: %s , IPs: %s .\n",
		args[3], authDomain, args[2], hostname, cIPs)
	for _, v := range pushConf.DeviceKeys {
		data := &pushsdk.PushContent{
			Title:             notiTitle,
			Body:              notiBody,
			DeviceKey:         v,
			Level:             pushConf.NotificationLevel,
			Group:             "Security_RDPAlert",
			AutomaticallyCopy: "1",
			Copy:              hostname,
		}
		pCont = append(pCont, data)
	}
	return pCont, nil
}

// GetLocalIP returns the non loopback local IP of the host
func GetLocalIP() []string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		gLogger.Critical("get local ip failed:", err)
		return nil
	}
	gLogger.Debug("got ip addrs, Length: ", len(addrs))
	res := make([]string, 0)
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				payload := ipnet.IP.String()
				gLogger.Info("found local ip:", payload)
				res = append(res, payload)
			}
		}
	}
	if len(res) == 0 {
		gLogger.Debug("get local ip with length zero.")
		return nil
	}
	return res
}
