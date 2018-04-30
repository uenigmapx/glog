glog
====

在[golang/glog](https://github.com/golang/glog)的基础上做了一些修改。

## 修改的地方:
1. 增加每天切割日志文件的功能,程序运行时指定 --dailyRolling=true参数即可
2. 将日志等级由原来的INFO WARN ERROR FATAL改为DEBUG INFO ERROR FATAL
3. 增加日志输出等级设置, 当日志信息等级低于输出等级时则不输出日志信息
4. 将默认的刷新缓冲区时间由 20s 改为 5s
5. 添加日志分割颗粒度(-logparticle), 默认按日切割(d/day[default] -- 按日切割, m/month 按月切割)
6. 添加日志压缩选项(-logcompress), 默认无压缩(none[default]/zip/gzip/bzip2)


##使用示例
```go
func main() {
    //初始化命令行参数
    flag.Parse()
    //退出时调用，确保日志写入文件中
    defer glog.Flush()

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

启动时指定
log_dir 参数将日志文件保存到特定的目录
dailyrolling 启动日切割日志选项
outputseverity 设置日志文件记录基准

```bash
./demo --log_dir=./log --dailyrolling=true -outputseverity=DEBUG
```
