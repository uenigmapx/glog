glog
====

在[golang/glog](https://github.com/golang/glog)的基础上做了一些修改

## 修改

1. 增加每天切割日志文件的功能,程序运行时指定 --dailyRolling=true参数即可
2. ~~将日志等级由原来的INFO WARN ERROR FATAL改为DEBUG INFO ERROR FATAL~~
2. 将日志等级修改到 DEBUG INFO WARN ERROR FATAL, 同时日志输出只有 DEBUG 到 ERROR, FATAL 不输出到日志文件
3. 增加日志输出等级设置, 当日志信息等级低于输出等级时则不输出日志信息
4. 将默认的刷新缓冲区时间由 20s 改为 5s
5. 让不同的输出级别只输出到各自的日志中
6. 添加日志分割颗粒度(-logparticle), 默认按日切割(d/day[default] -- 按日切割, m/month 按月切割)
7. 添加日志压缩选项(-logcompress), 默认无压缩(none[default]/zip/gzip/bzip2)
8. 添加压缩标准(-logcountpercompress), 当前以 文件数 关联, 开启后超过一定的文件数会自动压缩归档(TODO: 按照 log 文件名进行归档)
9. 当日志过大时清空日志，防止占满硬盘空间，可通过选项(--keepbig)，可选为(true,false[default])
10. 未使用flag.Parse()时，提示并保持文件输出
11. 增加配置文件配置日志记录, 优先级小于命令行参数

## 配置文件格式(log.json)

```json
{
      // display level: debug, info, warning, error[default], fatal
      "display_level": "error",
      // 文件记录等级: debug, info[default], warning, error, fatal
      "output_level": "info",
      // 日志路径, 默认临时文件夹
      "log_dir": "/tmp",
      // 是否保留大文件, 默认false
      "keep_big": "false",
      // 显示输出通道, (stdout[default], stderr)
      "display_channel": "stdout",
}
```

## 使用示例

```go
func main() {
    //初始化命令行参数
    flag.Parse()
    //退出时调用，确保日志写入文件中
    defer glog.Flush()

    // 解析后处理
    // 该方法已经包含 flag.Parse(), 因此可以选择注释上面的 flag.Parse()
    glog.Configure()

    /* 丢弃了该修改方法, 修改到 severity 的 flag.Value 方法上
    //一般在测试环境下设置输出等级为DEBUG，线上环境设置为INFO
    //glog.SetLevelString("DEBUG")
    */

    glog.Info("hello, glog")
    glog.Warning("warning glog")
    glog.Error("error glog")

    glog.Infof("info %d", 1)
    glog.Warningf("warning %d", 2)
    glog.Errorf("error %d", 3)
 }

```

```bash
./demo -log_dir=./log -rolling=true -outputseverity=DEBUG
```

## 详细参数

```bash
  -alsologtostderr
    	同时输出到文件和标准输出 <log to standard error as well as files>
  -log_backtrace_at value
    	当记录到 file:N , 则同时记录堆栈信息 <when logging hits line file:N, emit a stack trace>
  -log_dir string
    	If non-empty, write log files in this directory
  -logcompress string
    	压缩记录文件 <compress method(zip/gzip/none[default])> (default "none")
  -logcountpercompress value
        执行压缩需要的'最少'文件数<default is 0>
  -logparticle string
    	切割文件时的颗粒度 <particle size in rolling logfile (d/day--daily[default], m/month--monthly)> (default "d")
  -logtostderr
    	记录到标准错误输出而不是文件(覆盖 alsoToStderr) <log to standard error instead of files(cover alsoToStderr)>
  -outputseverity value
    	输出该等级之上的到记录文件 <logs at or above this content go to log file>
  -rolling
    	是否做按日(默认)或按月的文件切割 <weather to handle log files daily(default) or monthly>
  -rootdir string
    	用于暂存日志文件等的根目录 (default "/run/bin")
  -stderrthreshold value
    	输出该等级之上的到标准输出 <logs at or above this threshold go to stderr>
  -v value
    	V 记录器的记录等级 <log level for V logs>
  -vmodule value
    	文件过滤设置, 用 ',' 分隔 <comma-separated list of pattern=N settings for file-filtered logging>
```
