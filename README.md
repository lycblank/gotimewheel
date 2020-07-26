# gotimewheel
golang实现的分层时间轮，参照linux的时间轮实现

# 安装
``` shell
go get -u github.com/lycblank/gotimewheel
```

# 使用
```go
tw := NewTimeWheel(time.Second)
tw.AddTimer(2*time.Second, func(arg interface{}){
    beforeTime := arg.(time.Time)
    nowTime := time.Now()
    fmt.Println(beforeTime.Unix(), nowTime.Unix())
    }, time.Now())
time.Sleep(3*time.Second)
```
