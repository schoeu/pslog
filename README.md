# pslog

## 简介
无依赖，跨平台，零部署成本的服务端硬件资源监控。

获取机器基础数据，包括CPU，磁盘，网络，内存相关信息，按配置写入指定日志文件，供下游分析展现等。

## 快速使用
### 1. 下载

下载对应机器版本可执行文件

- [mac版本](https://github.com/schoeu/psloger/raw/master/pslog_mac)
- [linux 32位版本](https://github.com/schoeu/pslog/raw/master/pslog_linux32)
- [linux 62位版本](https://github.com/schoeu/pslog/raw/master/pslog_linux64)
- [windows 32位版本](https://github.com/schoeu/pslog/raw/master/pslog_32.exe)
- [windows 64位版本](https://github.com/schoeu/pslog/raw/master/pslog_64.exe)

### 2. 开箱即用
以mac为例
```
nohup ./pslog_mac &
```

## 配置
配置文件可以定义输出的日志格式，日志路径，输出时间间隔。配置非必须。
``` json
{
    "interval": 60000,
    "logFormat": "$logicalCores|$physicalCores|$percentPerCpu|$cpuPercent|$cpuModel|$memTotal|$memUsed|$memUsedPercent|$bytesRecv|$bytesSent|$diskTotle|$diskUsed|$diskUsedPercent",
    "logPath": "./psinfo_logs"
}
```

其中：

- `interval`为输出日志间隔，非必须，默认值为`60000`（60秒）
- `logPath`为日志输出路径，非必须，默认值为`./psinfo_logs`，则会自动生成对应目录及文件，并写入日志
- `logFormat`为输出的日志格式，可自定义，如只需要内存使用率，cpu使用率，磁盘使用率并以`^`间隔的，则该格式字符串为`$cpuPercent^$memUsedPercent^$diskUsedPercent`，日志中每行内容格式为`16.69^74.65^`


详细的字段介绍如下表

|占位符|含义|
|--|--|
|$logicalCores|逻辑核数|
|$physicalCores|物理核数|
|$percentPerCpu|单cpu使用率|
|$cpuPercent|cpu综合使用率|
|$cpuModel|cpu型号|
|$memTotal|总内存|
|$memUsed|已使用内存|
|$memUsedPercent|内存使用率|
|$bytesRecv|网卡下行速率|
|$bytesSent|网卡上行速率|
|$diskTotle|磁盘总空间|
|$diskUsed|磁盘已使用空间|
|$diskUsedPercent|磁盘使用占比|

## MIT License

Copyright (c) 2019 Schoeu

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

