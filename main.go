package main

import (
	"flag"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"time"
	"os"
	"path"
	"strings"
	"strconv"

    "github.com/shirou/gopsutil/cpu"
    "github.com/shirou/gopsutil/disk"
    "github.com/shirou/gopsutil/mem"
    "github.com/shirou/gopsutil/net"
)

type config struct {
	Name string
	Interval int
	LogFormat string
	LogPath string
}

func main() {
	var filePath, logPath string
	flag.StringVar(&filePath, "path", "", "path of config file.")
	flag.Parse()

	c := config{}
	logFormat := ""
	if filePath != "" {
		data, err := ioutil.ReadFile(filePath)
		err = json.Unmarshal(data, &c)
		errHandler(err)
		configLogPath := c.LogPath
		if configLogPath != "" {
			logPath = configLogPath
		}
	} else {
		cwd := getCwd()
		logPath = path.Join(cwd, "psinfo_logs")
		fmt.Println("Log file created at: ", logPath)
	}

	logFormatDefault := "$logicalCores|$physicalCores|$percentPerCpu|$cpuPercent|$cpuModel|$memTotal|$memUsed|$memUsedPercent|$bytesRecv|$bytesSent|$diskTotle|$diskUsed|$diskUsedPercent"
	logFormat = c.LogFormat
	if logFormat == "" {
		logFormat = logFormatDefault
	}

	makeDirP(logPath)
	logContent := getPsInfo(logFormat)

	during := c.Interval
	if during == 0 {
		during = 5000
	}
	timer(during, logPath, logContent)
}

func getPsInfo(formatStr string) string {
	v, _ := mem.VirtualMemory()
    c, _ := cpu.Info()
	// d, _ := disk.Usage("/")
    nv, _ := net.IOCounters(false)
	var mbNum uint64 = 1024 * 1024
	// var cpuModel []string
	for _, sub_cpu := range c {
		modelname := sub_cpu.ModelName
		// cores := sub_cpu.Cores
		fmt.Printf("        CPU       : %v\n", modelname)
		// cpuModel = append(cpuModel, modelname)
	}
	count1, _ := cpu.Counts(true)
	count2, _ := cpu.Counts(false)
	per1, _ := cpu.Percent(time.Second, true)
	per2, _ := cpu.Percent(time.Second, false)
	fmt.Println("       CPU", count1, count2, per1, per2, nv) // 逻辑核：8，物理核：4 使用率：[30.693069306930692 0 19.801980198019802 0 8 0 5.9405940594059405 0] 总使用率：[10.099750623441397]
	diskInfo, _ := disk.Partitions(true)

	for _, v := range diskInfo {
		device := v.Mountpoint
		distDetial, _ := disk.Usage(device)
		if distDetial != nil {
			fmt.Printf("   %v     HD        : %v MB  Free: %v MB Usage:%f%%\n",v.Device, distDetial.Total/mbNum, distDetial.Free/mbNum, distDetial.UsedPercent)
		}
	}

	// |占位符|含义|
	// |--|--|
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
	// |$disk_totle|磁盘总空间|
	// |$disk_used|磁盘已使用空间|
	// |$disk_used_percent|磁盘使用占比|

	logicalCoresStr := strings.Replace(formatStr, "$logicalCores", strconv.Itoa(count1), -1)
	physicalCoresStr := strings.Replace(logicalCoresStr, "$physicalCores", strconv.Itoa(count2), -1)
	percentPerCpuStr := strings.Replace(physicalCoresStr, "$percentPerCpu", fmt.Sprintf("%f", per1), -1)
	cpuPercentStr := strings.Replace(percentPerCpuStr, "$cpuPercent", fmt.Sprintf("%f", per2), -1)
	cpuModelStr := strings.Replace(cpuPercentStr, "$cpuModel", "", -1)
	memTotalStr := strings.Replace(cpuModelStr, "$memTotal", fmt.Sprintf("%d", v.Total/mbNum), -1)
	memUsedStr := strings.Replace(memTotalStr, "$memUsed", fmt.Sprintf("%d", v.Used/mbNum), -1)
	memUsedPercentStr := strings.Replace(memUsedStr, "$memUsedPercent", fmt.Sprintf("%f", v.UsedPercent), -1)
	bytesRecvStr := strings.Replace(memUsedPercentStr, "$bytesRecv", "", -1)
	bytesSentStr := strings.Replace(bytesRecvStr, "$bytesSent", "", -1)
	diskTotleStr := strings.Replace(bytesSentStr, "$diskTotle", "", -1)
	diskUsedStr := strings.Replace(diskTotleStr, "$diskUsed", "", -1)
	rsStr := strings.Replace(diskUsedStr, "$diskUsedPercent", "", -1)

	return rsStr + "\n"
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
		dir, err := os.Getwd()
		errHandler(err)
		p = path.Join(dir, p)
	}
	rs := path.Dir(p)
	os.MkdirAll(rs, os.ModePerm)
}

func appendText(p string, content string) {
	file, err := os.OpenFile(p, os.O_CREATE | os.O_APPEND | os.O_RDWR, 0600)
	errHandler(err)
	defer file.Close()

	if _, err = file.WriteString(content); err != nil {
		panic(err)
	}
}

func timer(interval int, p, content string) {
	during := interval / 1000
	d := time.Duration(time.Second * time.Duration(during))
	t := time.NewTicker(d)
	defer t.Stop()
	for {
		<- t.C
		appendText(p, content)
	}
}
