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

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

type config struct {
	Interval  int
	LogFormat string
	LogPath   string
}

var (
	logFormatDefault = "$dateTime|$logicalCores|$physicalCores|$percentPerCpu|$cpuPercent|$cpuModel|$memTotal|$memUsed|$memUsedPercent|$bytesRecv|$bytesSent|$diskTotal|$diskUsed|$diskUsedPercent"
	fmtLog           = logFormatDefault
	mbNum            = uint64(1024 * 1024)
	timeFormat       = "2006-01-02T15:04:05"
	during           = 60000
	recv             float32
	sent             float32
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
	normalInfo()

	timer(during, logPath)
}

// |$dateTime|日期时间戳|
// |$logical_cores|逻辑核数|
// |$physical_cores|物理核数|
// |$percent_percpu|单cpu使用率|
// |$cpu_percent|cpu综合使用率|
// |$cpu_model_name|cpu型号|
// |$mem_total|总内存|
// |$mem_used|已使用内存|
// |$mem_used_percent|内存使用率|
// |$bytes_recv|网卡下行速率|
// |$bytes_sent|网卡上行速率|
// |$disk_total|磁盘总空间|
// |$disk_used|磁盘已使用空间|
// |$disk_used_percent|磁盘使用占比|

func normalInfo() {
	v, _ := mem.VirtualMemory()
	c, _ := cpu.Info()
	logicalCount, _ := cpu.Counts(true)
	physicalCount, _ := cpu.Counts(false)

	var cpuModel []string
	for _, subCpu := range c {
		cpuModel = append(cpuModel, fmt.Sprintf(`"%s"`, subCpu.ModelName))
	}

	fmtLog = strings.Replace(fmtLog, "$memTotal", fmt.Sprintf("%dMB", v.Total/mbNum), -1)
	fmtLog = strings.Replace(fmtLog, "$logicalCores", strconv.Itoa(logicalCount), -1)
	fmtLog = strings.Replace(fmtLog, "$physicalCores", strconv.Itoa(physicalCount), -1)
	fmtLog = strings.Replace(fmtLog, "$cpuModel", fmt.Sprintf(`%s`, strings.Join(cpuModel, ",")), -1)
}

func getPsInfo(interval int) string {
	v, _ := mem.VirtualMemory()
	percentPerCpu, _ := cpu.Percent(time.Second, true)
	cpuPercent, _ := cpu.Percent(time.Second, false)
	diskInfo, _ := disk.Partitions(true)
	var diskTotal, diskUsed uint64
	for _, v := range diskInfo {
		device := v.Mountpoint
		distDetial, _ := disk.Usage(device)
		if distDetial != nil {
			diskTotal += distDetial.Total
			diskUsed += distDetial.Used
		}
	}

	nw, _ := net.IOCounters(false)
	parseNum := float32(uint64(interval) / 1000)
	var recvRate, sentRate float32
	if len(nw) > 0 && nw[0].Name == "all" {
		br := float32(nw[0].BytesRecv)
		bs := float32(nw[0].BytesSent)
		recvRate = (br - recv) / 1024 / parseNum
		sentRate = (bs - sent) / 1024 / parseNum

		// 初次获取上下行信息，矫正数据
		if recv == 0 || sent == 0 {
			recvRate = 0
			sentRate = 0
		}

		recv = br
		sent = bs
	}
	// recvRate := nw.BytesRecv / mbNum / interval / 1000
	// nw.BytesSent

	diskUsedPercent := float32(diskUsed) / float32(diskTotal)
	tmpLog := strings.Replace(fmtLog, "$diskTotal", fmt.Sprintf("%dGB", diskTotal/mbNum/1024), -1)
	tmpLog = strings.Replace(tmpLog, "$diskUsedPercent", fmt.Sprintf("%.2f", diskUsedPercent*100), -1)
	tmpLog = strings.Replace(tmpLog, "$diskUsed", fmt.Sprintf("%dGB", diskUsed/mbNum/1024), -1)
	tmpLog = strings.Replace(tmpLog, "$memUsedPercent", fmt.Sprintf("%.2f", v.UsedPercent), -1)
	tmpLog = strings.Replace(tmpLog, "$memUsed", fmt.Sprintf("%.2fMB", float32(v.Used)/float32(mbNum)), -1)
	tmpLog = strings.Replace(tmpLog, "$percentPerCpu", fmt.Sprintf("%.2f", percentPerCpu), -1)
	tmpLog = strings.Replace(tmpLog, "$cpuPercent", fmt.Sprintf("%.2f", cpuPercent[0]), -1)
	tmpLog = strings.Replace(tmpLog, "$bytesRecv", fmt.Sprintf("%.2fKB/s", recvRate), -1)
	tmpLog = strings.Replace(tmpLog, "$bytesSent", fmt.Sprintf("%.2fKB/s", sentRate), -1)
	tmpLog = strings.Replace(tmpLog, "$dateTime", time.Now().Format(timeFormat), -1)

	return tmpLog + "\n"
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
		content := getPsInfo(interval)
		appendText(p, content)
	}
}
