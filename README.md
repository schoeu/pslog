# pslog

## 简介
无依赖，跨平台，零部署成本的服务端硬件资源监控。

获取机器基础数据，包括CPU，磁盘，网络，内存相关信息，按配置写入指定日志文件，供下游分析展现等。

## 快速使用
### 1. 下载

下载对应机器版本可执行文件

- [mac版本](https://github.com/schoeu/psloger/raw/master/pslog_mac)
- [linux 32位版本](https://github.com/schoeu/pslog/raw/master/pslog_linux32)
- [linux 64位版本](https://github.com/schoeu/pslog/raw/master/pslog_linux64)
- [windows 32位版本](https://github.com/schoeu/pslog/raw/master/pslog_32.exe)
- [windows 64位版本](https://github.com/schoeu/pslog/raw/master/pslog_64.exe)

### 2. 开箱即用
以mac为例
```
chmod +x ./pslog_mac
nohup ./pslog_mac &
```

### 3. 配置

方法一：指定配置文件
```
./pslog_mac --path ./path_to_your_config_file
```

配置文件可以自定义输出的日志格式，日志路径，输出时间间隔。配置文件非必须。
``` json
{
    "interval": 60000,
    "logFormat": "$dateTime|$cpuPercent",
    "logPath": "~/.psinfo.log"
}
```

或

方法二：命令行参数

```
./pslog_mac --logFormat '$dateTime' --logPath './ps_logs' --interval 60000
```

配置文件与命令行参数含义相同，其中：

- `interval`为输出日志间隔，非必须，默认值为`60000`（60秒）
- `logPath`为日志输出路径，非必须，默认值为`~/.psinfo.log`，则会自动生成对应目录及文件，并写入日志
- `logFormat`为输出的日志格式，非必须，默认值为`"$dateTime|$logicalCores|$physicalCores|$percentPerCpu|$cpuPercent|$cpuModel|$memTotal|$memUsed|$memUsedPercent|$bytesRecv|$bytesSent|$diskTotal|$diskUsed|$diskUsedPercent"`。

用户可以使用其中任意字段自行拼接。比如需要内存使用率，cpu使用率，且想以`^`间隔，则该格式字符串为`$cpuPercent^$memUsedPercent`，生效后，日志中单行内容为`16.69^74.65`，代表cpu，内存使用率分别为`16.69%`，`74.65%`。


暂时支持以下字段

|占位符|含义|示例|备注|
|--|--|--|--|
|$dateTime|日期时间戳|2019-06-28T17:37:11|当前时间戳|
|$logicalCores|逻辑核数|8||
|$physicalCores|物理核数|4||
|$percentPerCpu|单cpu使用率|[33.66 3.00 30.00 3.96 30.00 3.96 27.72 3.96]|展现每一个逻辑核的使用率|
|$cpuPercent|cpu综合使用率|6.64|使用率为6.64%|
|$cpuModel|cpu型号|"Intel(R) Core(TM) i7-4750HQ CPU @ 2.00GHz"|多类核会以`,`隔开|
|$memTotal|总内存|8192MB|8GB，此处以MB来展现|
|$memUsed|已使用内存|5516.53MB|已使用了5516.53MB|
|$memUsedPercent|内存使用率|67.34|已使用占比67.34%|
|$bytesRecv|网卡下行速率|4.00KB/s|下行速率|
|$bytesSent|网卡上行速率|1.50KB/s|上行速率|
|$diskTotal|磁盘总空间|467GB|磁盘总计467G，不包括隐藏分区|
|$diskUsed|磁盘已使用空间|159GB|已使用159GB|
|$diskUsedPercent|磁盘使用占比|34.20|磁盘使用了34.20%|
|$load|负载|1.56,1.72,1.88|分别代表1分钟，5分钟，15分钟的系统负载|

**重要**：如果配合`pslog_agent`日志上报服务使用，则`logFormat`与`logPath`保持默认即可。

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
