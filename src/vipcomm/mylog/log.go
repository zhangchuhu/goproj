package mylog

import (
	"fmt"
	"github.com/astaxie/beego"
	seelog "github.com/cihub/seelog"
	"path/filepath"
)

var (
	logLevel = "info"
)

/*
* 使用说明，main函数初始化
* 日志文件默认路径  /data/yy/gologs/进程名/进程名.log
* 日志默认级别     info
*
* func main () {
*	vipcomm.InitLog()
*   defer vipcomm.FlushLog()
*   ....
* }
 */

// console log
func loadAppConfig() {
	appConfig := `<seelog minlevel="` + logLevel + `">
	<outputs formatid="common">
		<console/>
	</outputs>
	<formats>
		<format id="common" format="%Date %Time %LEV [%RelFile %Func] %Msg%n" />
	</formats>
</seelog>
`
	logger, err := seelog.LoggerFromConfigAsBytes([]byte(appConfig))
	if err != nil {
		fmt.Println(err)
		return
	}
	seelog.ReplaceLogger(logger)
}

func init() {
	loadAppConfig()
}

/*
* 初始化日志配置
* @name  日志文件名
* @path  日志存放路径
* @level 日志等级过滤，trace  debug  info  warn  error  critical
 */
func InitLogByArgs(name, level, path string) {
	if name == "" {
		return
	}

	appConfig := `<seelog minlevel="`

	// level
	if level == "" {
		appConfig += "info"
	} else {
		appConfig += level
	}
	appConfig += `">
	<outputs formatid="common">
		<rollingfile type="size" filename="`

	// log file path
	confPath := name + ".log"
	if path == "" {
		appConfig += filepath.Join("/data/yy/gologs/", name, confPath)
	} else {
		appConfig += filepath.Join(path, name, confPath)
	}

	appConfig += `" maxsize="100000" maxrolls="5"/>
	</outputs>
	<formats>
		<format id="common" format="%Date %Time %LEV [%RelFile %Func] %Msg%n" />
	</formats>
</seelog>
`
	logger, err := seelog.LoggerFromConfigAsBytes([]byte(appConfig))
	if err != nil {
		fmt.Println(err)
		return
	}
	seelog.ReplaceLogger(logger)
}

// 初始化日志，日志默认路径在/data/yy/golog/
func InitLog() {
	var consoleLog = beego.AppConfig.String("consoleLog")
	logLevel = beego.AppConfig.String("logLevel")
	fmt.Println(logLevel)
	if logLevel != "info" && logLevel != "debug" && logLevel != "error" {
		logLevel = "info"
	}

	if consoleLog != "" {
		// 默认控制台日志
		loadAppConfig()
	} else {
		InitLogByArgs(beego.BConfig.AppName, logLevel, "")
	}
}

// 根据配置文件解析日志配置
func InitLogConfFile(confFile string) {
	logger, err := seelog.LoggerFromConfigAsFile(confFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	seelog.ReplaceLogger(logger)
}

func FlushLog() {
	seelog.Flush()
}

// gorpc errlog 封装
func ErrorFunc(fmt string, args ...interface{}) {
	seelog.Errorf(fmt, args)
}
