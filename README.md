# remote_log_aliyun_go

remote_log_aliyun sdk go 版本。是将阿里云日志服务进行二次封装。

## 安装

```bash
go get -u github.com/doubility/remote-log-go
```

## 快速开始

可拷贝 examples 中的例子

**重点：Logger 申明为全局变量，初始化一次！！！**

初始化项目参数

```go
    log := remote_log_aliyun_go.NewLogger("winner-test-project")
```

main.go (go-test 为 go.mod module)

```go
package main

import (
    "go-test/logger"
)

func main() {
	logger.Logger.Info("记录info日志") // http上传日志
	logger.Logger.Warn("记录warn日志") // http上传日志

	logger.Logger.Debug("debug日志") // console打印日志
}
```

## 详细说明

```code
// 日志类型 可在查询时筛选
-remote_log_aliyun_go.Debug
-remote_log_aliyun_go.Info
-remote_log_aliyun_go.Warn
-remote_log_aliyun_go.Error
-remote_log_aliyun_go.Access
```

```go
import (
    remote_log_go "github.com/doubility/remote-log-go"
)

// 申明日志存储方式，一种日志类型可选择多种存储方式

// 日志上传到服务器
// (info、warn、error、access日志使用http上传到服务器)
httpTransport := remote_log_go.NewHttpTransport(remote_log_go.Info, remote_log_go.Warn, remote_log_go.Error, remote_log_go.Access)

// 日志输出到控制台
// (debug日志使用console打印)
consoleTransport := remote_log_go.NewConsoleTransport(remote_log_go.Debug)

// 实例化
// appName string 应用的名称（查询日志时可使用）
// storageDays number 日志存储天数 (最小30天，最大360天)
// transport transport ...interface{} 日志处理方式 接受HttpTransport和ConsoleTransport
Logger := remote_log_go.NewLogger(appName, storageDays, transport)

// 初始化
err := Logger.init()

// 记录各种类型的日志
Logger.Debug(string);
Logger.Info(string);
Logger.Warn(string);
Logger.Error(error);
Logger.Access(string);
```

## 注意

1、需要环境变量`GO_APP_LOG_PATH`，上传失败的日志将保存在此目录下。

2、需要环境变量`GO_ALIYUAN_ENDPOINT`，阿里云Endpoint参数。

3、需要环境变量`GO_ALIYUAN_ACCESSKEYID`，阿里云访问密钥AccessKeyId。

4、需要环境变量`GO_ALIYUAN_ACCESSKEYSECRET`，阿里云访问密钥AccessKeySecret。

5、阿里云文档 https://help.aliyun.com/document_detail/286951.html

