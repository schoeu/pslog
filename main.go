package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/schoeu/gopsinfo"
)

type config struct {
	Interval  int
	LogFormat string
	LogPath   string
}

var (
	logFormatDefault = "$dateTime|$logicalCores|$physicalCores|$percentPerCpu|$cpuPercent|$cpuModel|$memTotal|$memUsed|$memUsedPercent" +
		"|$recvSpeed|$sentSpeed|$diskTotal|$diskUsed|$diskUsedPercent|$load|$os|$platform|$platformFamily|$platformVersion"
	fmtLog           = logFormatDefault
	kbNum            = uint64(1024)
	mbNum            = kbNum * kbNum
	timeFormat       = "2006-01-02T15:04:05"
	during           = 60000
)

func main() {
	var filePath, logPath, formatStr, dist string
	var interval int
	homeDir, _ := os.UserHomeDir()

	flag.StringVar(&filePath, "path", "", "configuration file path.")
	flag.StringVar(&formatStr, "logFormat", logFormatDefault, "log format string.")
	flag.StringVar(&dist, "logPath", path.Join(homeDir, ".psinfo.log"), "log file path.")
	flag.IntVar(&interval, "interval", during, "interval timer")
	flag.Parse()

	c := config{}
	if filePath != "" {
		data, err := ioutil.ReadFile(filePath)
		err = json.Unmarshal(data, &c)
		errHandler(err)
		configLogPath := c.LogPath
		if configLogPath != "" {
			logPath = configLogPath
		}

		logFormat := c.LogFormat
		if logFormat != "" {
			fmtLog = logFormat
		}

		if c.Interval != 0 {
			during = c.Interval
		}
	} else {
		logPath = dist
		fmtLog = formatStr
		during = interval
	}

	fmt.Println("Log file created at: ", logPath)

	makeDirP(logPath)
	timer(during, logPath)
}

func getLogContent(interval int) string {
	rs := gopsinfo.GetPsInfo(interval)
	diskUsedPercent := float64(rs.DiskUsed) / float64(rs.DiskTotal)
	tmpLog := strings.Replace(fmtLog, "$diskTotal", fmt.Sprintf("%dGB", rs.DiskTotal/mbNum/kbNum), -1)
	tmpLog = strings.Replace(tmpLog, "$diskUsedPercent", parseFloatNum(diskUsedPercent*100), -1)
	tmpLog = strings.Replace(tmpLog, "$diskUsed", fmt.Sprintf("%dGB", rs.DiskUsed/mbNum/kbNum), -1)
	tmpLog = strings.Replace(tmpLog, "$memUsedPercent", parseFloatNum(rs.MemUsedPercent), -1)
	tmpLog = strings.Replace(tmpLog, "$memUsed", fmt.Sprintf("%.2fMB", float32(rs.MemUsed)/float32(mbNum)), -1)
	tmpLog = strings.Replace(tmpLog, "$percentPerCpu", fmt.Sprintf("%.2f", rs.PercentPerCpu), -1)
	tmpLog = strings.Replace(tmpLog, "$cpuPercent", parseFloatNum(rs.CpuPercent), -1)
	tmpLog = strings.Replace(tmpLog, "$recvSpeed", fmt.Sprintf("%.2fKB/s", rs.RecvSpeed / float64(kbNum)), -1)
	tmpLog = strings.Replace(tmpLog, "$sentSpeed", fmt.Sprintf("%.2fKB/s", rs.SentSpeed / float64(kbNum)), -1)
	//fmt.Println(rs.RecvSpeed)
	tmpLog = strings.Replace(tmpLog, "$dateTime", time.Now().Format(timeFormat), -1)
	tmpLog = strings.Replace(tmpLog, "$load", strings.Join(rs.Load, ","), -1)
	tmpLog = strings.Replace(tmpLog, "$memTotal", fmt.Sprintf("%dMB", rs.MemTotal/mbNum), -1)
	tmpLog = strings.Replace(tmpLog, "$logicalCores", strconv.Itoa(rs.LogicalCores), -1)
	tmpLog = strings.Replace(tmpLog, "$physicalCores", strconv.Itoa(rs.PhysicalCores), -1)
	tmpLog = strings.Replace(tmpLog, "$cpuModel", fmt.Sprintf(`%s`, strings.Join(rs.CpuModel, ",")), -1)
	tmpLog = strings.Replace(tmpLog, "$os", rs.Os, -1)
	tmpLog = strings.Replace(tmpLog, "$platform", rs.Platform, -1)
	tmpLog = strings.Replace(tmpLog, "$platformFamily", rs.PlatformFamily, -1)
	tmpLog = strings.Replace(tmpLog, "$platformVersion", rs.PlatformVersion, -1)

	return tmpLog + "\n"
}

func parseFloatNum(n float64) string {
	return fmt.Sprintf("%.2f", n)
}

func errHandler(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func getCwd() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	return dir
}

func makeDirP(p string) {
	if !path.IsAbs(p) {
		p = path.Join(getCwd(), p)
	}
	rs := path.Dir(p)
	os.MkdirAll(rs, os.ModePerm)
}

func appendText(p string, content string) {
	file, err := os.OpenFile(p, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0600)
	errHandler(err)
	defer file.Close()

	_, err = file.WriteString(content)
	errHandler(err)
}

func timer(interval int, p string) {
	d := time.Duration(time.Millisecond * time.Duration(interval))
	t := time.NewTicker(d)
	defer t.Stop()

	for {
		<-t.C
		content := getLogContent(interval)
		appendText(p, content)
	}
}
