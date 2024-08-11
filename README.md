# etop
**etop** is system monitor tool inspired by [atop](https://github.com/Atoptool/atop) and 
[below](https://github.com/facebookincubator/below) for linux, written in golang.

## Feature

* **Cgroup v2** collect cgroup v2 if available.
* **Persistent record** record all samples into disk file. so it is easy to investigate historical issue.
* **Dump structured Information** dump mode not only output plain text, but also json which can import into database. even send data through OTLP to any OTel backends 
(e.g `Grafana`)
* **Filter and sort of process view** any field can be filtered and sorted.

## Installation

### Prebuilt binaries
Download binary from [release page](https://github.com/xixiliguo/etop/releases).   
RPM file was proviced too, it include systemd unit file.  
### Source
install from source code
```bash
go install github.com/xixiliguo/etop@latest
```

## Quickstart

record sampels by every 5 seconds
```
etop record -i 5
```
enter interactive mode to view data 
```
etop report
```
dump cpu info by text format (default from 00:00 until now)
```
etop dump cpu
```