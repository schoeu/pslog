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

var (
    logFormatDefault = "$logicalCores|$physicalCores|$percentPerCpu|$cpuPercent|$cpuModel|$memTotal|$memUsed|$memUsedPercent|$bytesRecv|$bytesSent|$diskTotle|$diskUsed|$diskUsedPercent"
    fmtLog = logFormatDefault
    kbNum = uint64(1024 * 1024)
)

func main() {
    var filePath, logPath string
    flag.StringVar(&filePath, "path", "", "path of config file.")
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
    } else {
        cwd := getCwd()
        logPath = path.Join(cwd, "psinfo_logs")
        fmt.Println("Log file created at: ", logPath)
    }

    logFormat := c.LogFormat
    if logFormat == "" {
        logFormat = logFormatDefault
    }

    makeDirP(logPath)
    normalInfo()

    during := c.Interval
    if during == 0 {
        during = 5000
    }
    timer(during, logPath)
}

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

func normalInfo() {
    v, _ := mem.VirtualMemory()
    c, _ := cpu.Info()
    logicalCount, _ := cpu.Counts(true)
    physicalCount, _ := cpu.Counts(false)
    // d, _ := disk.Usage("/")
    nv, _ := net.IOCounters(false)
    fmt.Println(nv)
    var cpuModel []string
    for _, sub_cpu := range c {
        cpuModel = append(cpuModel, fmt.Sprintf(`"%s"`, sub_cpu.ModelName))
    }

    fmtLog = strings.Replace(fmtLog, "$memTotal", fmt.Sprintf("%dMB", v.Total/kbNum), -1)
    fmtLog = strings.Replace(fmtLog, "$logicalCores", strconv.Itoa(logicalCount), -1)
    fmtLog = strings.Replace(fmtLog, "$physicalCores", strconv.Itoa(physicalCount), -1)
    fmtLog = strings.Replace(fmtLog, "$cpuModel", fmt.Sprintf(`%s`, strings.Join(cpuModel, ",")), -1)
}

func getPsInfo() string {
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
    
    diskUsedPercent := float32(diskUsed) / float32(diskTotal)
    tmpLog := strings.Replace(fmtLog, "$diskTotle", fmt.Sprintf("%dGB", diskTotal / kbNum / 1024 ), -1)
    tmpLog = strings.Replace(tmpLog, "$diskUsedPercent", fmt.Sprintf("%.2f", diskUsedPercent), -1)
    tmpLog = strings.Replace(tmpLog, "$diskUsed", fmt.Sprintf("%dGB", diskUsed / kbNum / 1024), -1)
    tmpLog = strings.Replace(tmpLog, "$memUsedPercent", fmt.Sprintf("%.2f", v.UsedPercent), -1)
    tmpLog = strings.Replace(tmpLog, "$memUsed", fmt.Sprintf("%.2fMB", float32(v.Used)/float32(kbNum)), -1)
    tmpLog = strings.Replace(tmpLog, "$percentPerCpu", fmt.Sprintf("%.2f", percentPerCpu), -1)
    tmpLog = strings.Replace(tmpLog, "$cpuPercent", fmt.Sprintf("%.2f", cpuPercent[0]), -1)
    tmpLog = strings.Replace(tmpLog, "$bytesRecv", "", -1)
    tmpLog = strings.Replace(tmpLog, "$bytesSent", "", -1)

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
        cwd := getCwd()
        p = path.Join(cwd, p)
    }
    rs := path.Dir(p)
    os.MkdirAll(rs, os.ModePerm)
}

func appendText(p string, content string) {
    file, err := os.OpenFile(p, os.O_CREATE | os.O_APPEND | os.O_RDWR, 0600)
    errHandler(err)
    defer file.Close()

    _, err = file.WriteString(content)
    errHandler(err)
}

func timer(interval int, p string) {
    during := interval / 1000
    d := time.Duration(time.Second * time.Duration(during))
    t := time.NewTicker(d)
    defer t.Stop()

    for {
        <- t.C
        content := getPsInfo()
        appendText(p, content)
    }
}
