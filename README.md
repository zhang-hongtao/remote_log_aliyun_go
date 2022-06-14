# remote_log_aliyun_go

remote_log_aliyun sdk go 版本。是将阿里云日志服务进行二次封装。

## 安装

```bash
go get -u github.com/zhang-hongtao/remote_log_aliyun_go
```

## 快速开始

可拷贝 examples 中的例子

**重点**

初始化项目参数

```go
package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"remote_log_aliyun_go"
	"syscall"
	"time"
)

func main() {
	logger := remote_log_aliyun_go.NewLogger("项目名称")
	err := logger.Init()
	if err != nil {
		fmt.Print("初始化错误:", err.Error())
	}
    logger.Logger.Info("记录info日志") // http上传日志
	logger.Logger.Warn("记录warn日志") // http上传日志

	logger.Logger.Debug("debug日志") // console打印日志
	
    /// 程序退出时 主动关闭服务
    logger.Close()
	
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

## 注意

1、需要环境变量`GO_APP_LOG_PATH`，上传失败的日志将保存在此目录下。

2、需要环境变量`GO_ALIYUAN_ENDPOINT`，阿里云Endpoint参数。

3、需要环境变量`GO_ALIYUAN_ACCESSKEYID`，阿里云访问密钥AccessKeyId。

4、需要环境变量`GO_ALIYUAN_ACCESSKEYSECRET`，阿里云访问密钥AccessKeySecret。

5、程序在退出时需主动调用 `logger.Close()` 退出程序

6、使用该项目时需主动在日志服务控制台创建项目和日志名称

7、阿里云文档 https://help.aliyun.com/document_detail/286951.html

