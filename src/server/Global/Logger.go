package Global
import (
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

)
/*
#include<stdio.h>
char str[80];
const char* build_time(void)
{
   	sprintf(str, "%s%s", __DATE__,__TIME__);
	return str;
}
import "C"
*/

var Logger *logs.BeeLogger;
var (
	buildTime = ""
	version = ""
)

type LoggerManager struct {

	beego.Controller
}

func Init_Logs() {

	t := time.Now()
	timestamp := t.UTC().Unix()
	timestring := strconv.FormatInt(timestamp,10)
	version = timestring+"-"+ buildTime
	Logger = logs.NewLogger()
	Logger.SetLogger(logs.AdapterConsole)
	Logger.SetLogger(logs.AdapterFile,`{"filename":"logs/`+timestring+buildTime+`Tetris.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10,"color":true}`)
	Logger.Debug("this is a debug message")
	Logger.EnableFuncCallDepth(true)
	Logger.Trace("version:",timestring+"-"+buildTime)

}