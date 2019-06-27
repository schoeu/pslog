package main

import (
	"flag"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"time"
	"os"
	"path"

    "github.com/shirou/gopsutil/cpu"
    "github.com/shirou/gopsutil/disk"
    // "github.com/shirou/gopsutil/host"
    "github.com/shirou/gopsutil/mem"
    "github.com/shirou/gopsutil/net"
)

type config struct {
	Name string
	Interval int
	Format string
	LogPath string
}

func main() {
	var filePath, logPath string
	flag.StringVar(&filePath, "path", "", "path of config file.")
	flag.Parse()

	c := config{}
	if filePath != "" {
		data, err := ioutil.ReadFile(filePath)
		err = json.Unmarshal(data, &c)
		errHandler(err)
		fmt.Println(c)
		configLogPath := c.LogPath
		if configLogPath == "" {
			logPath = configLogPath
		}
	} else {
		cwd := getCwd()
		logPath = path.Join(cwd, "psinfo_logs")
		fmt.Println("log file created at ", logPath)
	}
}

func getPsInfo() {
	v, _ := mem.VirtualMemory()
    c, _ := cpu.Info()
	// d, _ := disk.Usage("/")
    // n, _ := host.Info()
    nv, _ := net.IOCounters(false)
    // boottime, _ := host.BootTime()
    // btime := time.Unix(int64(boottime), 0).Format("2006-01-02 15:04:05")

	var mbNum uint64 = 1024 * 1024

	fmt.Printf("        Mem       : %v MB  Free: %v MB Used:%vMB Usage:%f%%\n", v.Total/mbNum, v.Available/mbNum, v.Used/mbNum, v.UsedPercent)
	
	for _, sub_cpu := range c {
		modelname := sub_cpu.ModelName
		// cores := sub_cpu.Cores
		fmt.Printf("        CPU       : %v\n", modelname)
	}
	count1, _ := cpu.Counts(true)
	count2, _ := cpu.Counts(false)
	per1, _ := cpu.Percent(time.Second, true)
	per2, _ := cpu.Percent(time.Second, false)
	fmt.Println("       CPU", count1, count2, per1, per2) // 虚拟核：8，物理核：4 使用率：[30.693069306930692 0 19.801980198019802 0 8 0 5.9405940594059405 0] 总使用率：[10.099750623441397]
	diskInfo, _ := disk.Partitions(true)

	for _, v := range diskInfo {
		device := v.Mountpoint
		distDetial, _ := disk.Usage(device)
		if distDetial != nil {
			fmt.Printf("   %v     HD        : %v MB  Free: %v MB Usage:%f%%\n",v.Device, distDetial.Total/1024/1024, distDetial.Free/1024/1024, distDetial.UsedPercent)
		}
	}

	fmt.Println(nv)
	// fmt.Println(disk.Partitions(true))

    // fmt.Printf("        Network: %v bytes / %v bytes\n", nv[0].BytesRecv, nv[0].BytesSent)
    // fmt.Printf("        SystemBoot:%v\n", btime)
    // fmt.Printf("        CPU Used    : used %f%% \n", cc[0])
    // fmt.Printf("        HD        : %v GB  Free: %v GB Usage:%f%%\n", d)
    // fmt.Printf("        OS        : %v(%v)   %v  \n", n.Platform, n.PlatformFamily, n.PlatformVersion)
	// fmt.Printf("        Hostname  : %v  \n", n)
	// fmt.Println(c)
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

func MakeDirP(p string) {
	if !path.IsAbs(p) {
		dir, err := os.Getwd()
		errHandler(err)
		p = path.Join(dir, p)
	}
	rs := path.Dir(p)
	os.MkdirAll(rs, os.ModePerm)
}
