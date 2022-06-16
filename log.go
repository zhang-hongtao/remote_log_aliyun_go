package remote_log_aliyun_go

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/aliyun-log-go-sdk/producer"
	"google.golang.org/protobuf/proto"
)

/*
 * @Author: your name
 * @Date: 2022-06-07 09:18:11
 * @LastEditTime: 2022-06-08 10:57:27
 * @LastEditors: zhanghongtaodeMacBook-Pro.local
 * @Description: 使用阿里云日志SDK log 日志存储
 * @FilePath: /winnerLog/log/Log.go
 */
var (
	producerInstance *producer.Producer
	logger           *Logger
)

type Logger struct {
	appName         string
	endpoint        string
	projectName     string
	logStoreName    string
	accessKeyId     string
	accessKeySecret string
}

/**
* @description: 实例化对象
* @param {string} projectName log的项目名称
* @param {string} appName 应用名称
* @param {string} logName 日志文件名
* @return 返回Logger对象
 */
func NewLogger(projectName, appName, logName string) *Logger {
	if projectName == "" {
		panic(errors.New("projectName cannot be empty"))
	}
	if appName == "" {
		panic(errors.New("appname cannot be empty"))
	}
	if logName == "" {
		panic(errors.New("logName cannot be empty"))
	}
	goPath := os.Getenv("GO_APP_LOG_PATH")
	if goPath != "" {
		ErrorLogPath = fmt.Sprintf("%v/%v/remote_logs", goPath, appName)
		os.MkdirAll(ErrorLogPath, os.ModePerm)
	} else {
		panic(errors.New("invalid env GO_APP_LOG_PATH"))
	}
	// 配置AccessKey、服务入口、Project名称、Logstore名称等相关信息。
	// 日志服务的服务入口。更多信息，请参见服务入口。
	// 此处以杭州为例，其它地域请根据实际情况填写。
	endpoint := os.Getenv("GO_ALIYUAN_ENDPOINT")
	if endpoint == "" {
		panic(errors.New("invalid env GO_ALIYUAN_ENDPOINT"))
	}
	// 阿里云访问密钥AccessKey。更多信息，请参见访问密钥。阿里云账号AccessKey拥有所有API的访问权限，风险很高。强烈建议您创建并使用RAM用户进行API访问或日常运维。
	accessKeyId := os.Getenv("GO_ALIYUAN_ACCESSKEYID")
	if accessKeyId == "" {
		panic(errors.New("invalid env GO_ALIYUAN_ACCESSKEYID"))
	}
	accessKeySecret := os.Getenv("GO_ALIYUAN_ACCESSKEYSECRET")
	if accessKeySecret == "" {
		panic(errors.New("invalid env GO_ALIYUAN_ACCESSKEYSECRET"))
	}
	// RAM用户角色的临时安全令牌。此处取值为空，表示不使用临时安全令牌。更多信息，请参见授权用户角色。
	// securityToken = ""
	// 创建LogStore。
	// logStoreName = "remote_logs_" + appName
	logger = &Logger{
		projectName:     projectName,
		appName:         appName,
		logStoreName:    logName,
		endpoint:        endpoint,
		accessKeyId:     accessKeyId,
		accessKeySecret: accessKeySecret,
	}
	return logger
}

/**
 * @description: 初始化log 实例
 */
func (l *Logger) Init() {
	producerConfig := producer.GetDefaultProducerConfig()
	producerConfig.Endpoint = logger.endpoint
	producerConfig.AccessKeyID = logger.accessKeyId
	producerConfig.AccessKeySecret = logger.accessKeySecret
	producerInstance = producer.InitProducer(producerConfig)
	producerInstance.Start() // 启动producer实例
}

/**
 * @description:保存log
 * @param {string} content
 * @param {string} level
 */
func (l *Logger) SavaLog(content string, level string) {
	log := &sls.Log{
		Time: proto.Uint32(uint32(time.Now().Unix())),
		Contents: []*sls.LogContent{{
			Key:   proto.String("message"),
			Value: proto.String(content),
		}, {
			Key:   proto.String("level"),
			Value: proto.String(level),
		}},
	}
	ip, err := externalIP()
	if err != nil {
		fmt.Println(err)
	}
	err1 := producerInstance.SendLogWithCallBack(logger.projectName, logger.logStoreName, logger.appName, ip.String(), log, &Callback{})
	if err1 != nil {
		formatConsoleErr(content, level, err1.Error())
	}
}

func (l *Logger) Debug(a ...interface{}) {
	fmt.Println(a...)
}

func (l *Logger) Info(message string) {
	l.SavaLog(message, "Info")
}

func (l *Logger) Warn(message string) {
	l.SavaLog(message, "Warn")
}

func (l *Logger) Error(message error) {
	l.SavaLog(message.Error(), "Error")
}

func (l *Logger) Access(message string) {
	l.SavaLog(message, "Access")
}

func (l *Logger) Close() {
	producerInstance.SafeClose()
}

/**
 * 发送日志的回调
 */
type Callback struct {
}

func (callback *Callback) Success(result *producer.Result) {
	// attemptList := result.GetReservedAttempts() // 遍历获得所有的发送记录
	// for _, attempt := range attemptList {
	// 	fmt.Println(attempt)
	// }
}

func (callback *Callback) Fail(result *producer.Result) {
	if !result.IsSuccessful() {
		FailUploadLog(formatConsole(result))
	}
}

/**
 * @description: 日志格式化 日志发送失败
 * @param {*logger.LogInfo} log
 * @return {*}
 */
func formatConsole(result *producer.Result) string {
	var s strings.Builder
	s.WriteString("ErrCode:")
	s.WriteString(result.GetErrorCode()) // 获得最后一次发送失败错误码
	s.WriteString(" ErrMsg:")
	s.WriteString(result.GetErrorMessage()) // 获得最后一次发送失败信息
	s.WriteString(" RequestId:")
	s.WriteString(result.GetRequestId()) // 获得最后一次发送失败请求Id
	s.WriteString(" TimeStampMs:")
	s.WriteString(fmt.Sprint(result.GetTimeStampMs())) // 获得最后一次发送失败请求时间
	s.WriteString(" ReservedAttempts:")
	s.WriteString(fmt.Sprint(result.GetReservedAttempts())) // 获得producerBatch 每次尝试被发送的信息
	return s.String()
}

/**
 * @description: 日志服务 报错的日志
 * @param {*logger.LogInfo} log
 * @return {*}
 */
func formatConsoleErr(content, level, errMsg string) string {
	var s strings.Builder
	s.WriteString(fmt.Sprint(time.Now().Unix())) // 报错时间
	s.WriteString(" content:")
	s.WriteString(content) // 日志内容
	s.WriteString(" level:")
	s.WriteString(level) //日志级别
	s.WriteString(" err:")
	s.WriteString(errMsg) // 日志错误信息
	return s.String()
}
