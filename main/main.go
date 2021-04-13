package main

//导入向文件内写日志方法包
import (
	"time"

	"github.com/day13Homeweok_filelog/filelog"
)

//编写一个日志库
//需求分析:
//1.支持往不同的地方输出日志
//2.日志分级别（1.debug 2.Trace 3.Info 4.Warning  5.Error  6.Fatal）
//3.日志要支持开关控制,上线之后只打印某个级别的信息
//4.完整的日志记录要包含时间，行号，文件名，日志级别，日志信息
//5.日志文件要切割
func main() {
	log := filelog.Newlog("Debug", "./", "day13Homeweok_filelog.log", 10*1024) //指定日志文件记录最低等级的日志，指定日志文件的存放目录，指定日志文件的文件名，指定日志文件的大小
	for {
		log.Debug("这是一个Debug级别的日志")
		log.Trace("这是一个Trace级别的日志")
		log.Info("这是一个Info级别的日志")
		log.Warning("这是一个Warning级别的日志")
		log.Error("这是一个Error级别的日志")
		log.Fatal("这是一个Fatal级别的日志")
		time.Sleep(time.Second)
	}

}
